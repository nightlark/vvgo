package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/virtual-vgo/vvgo/pkg/sessions"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestLoginHandler_ServeHTTP(t *testing.T) {
	t.Run("get/redirect", func(t *testing.T) {
		ts := httptest.NewServer(LoginHandler{
			Sessions: sessions.NewStore(sessions.Opts{}),
		})
		defer ts.Close()

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		require.NoError(t, err, "cookiejar.New")
		client := noFollow(&http.Client{Jar: jar})
		resp, err := client.Get(ts.URL)
		require.NoError(t, err, "client.Get")
		assert.Equal(t, http.StatusFound, resp.StatusCode)
		assert.Equal(t, "/", resp.Header.Get("Location"), "location")
	})

	t.Run("get/view", func(t *testing.T) {
		wantCode := http.StatusOK
		wantBytes, err := ioutil.ReadFile("testdata/login.html")
		if err != nil {
			t.Fatalf("ioutil.ReadFile() failed: %v", err)
		}

		ts := httptest.NewServer(LoginHandler{
			Sessions: sessions.NewStore(sessions.Opts{}),
		})
		defer ts.Close()

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		require.NoError(t, err, "cookiejar.New")
		client := noFollow(&http.Client{Jar: jar})
		resp, err := client.Get(ts.URL)
		require.NoError(t, err, "client.Get")
		assert.Equal(t, wantCode, resp.StatusCode)
		var respBody bytes.Buffer
		_, err = respBody.ReadFrom(resp.Body)
		require.NoError(t, err, "resp.Body.Read() failed")
		origBody := strings.TrimSpace(respBody.String())

		m := minify.New()
		m.AddFunc("text/html", html.Minify)
		var gotBuf bytes.Buffer
		if err := m.Minify("text/html", &gotBuf, &respBody); err != nil {
			panic(err)
		}
		gotBody := gotBuf.String()

		var wantBuf bytes.Buffer
		if err := m.Minify("text/html", &wantBuf, bytes.NewReader(wantBytes)); err != nil {
			panic(err)
		}
		wantBody := wantBuf.String()
		if !assert.Equal(t, wantBody, gotBody, "body") {
			t.Logf("Got Body:\n%s\n", origBody)
		}
	})

	t.Run("post/failure", func(t *testing.T) {
		ts := httptest.NewServer(LoginHandler{
			Sessions: sessions.NewStore(sessions.Opts{}),
		})
		defer ts.Close()

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		require.NoError(t, err, "cookiejar.New")
		client := noFollow(&http.Client{Jar: jar})

		resp, err := client.Post(ts.URL, "text/plain", nil)
		require.NoError(t, err, "client.Get")
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		var gotBody bytes.Buffer
		gotBody.ReadFrom(resp.Body)
		assert.Equal(t, "authorization failed", strings.TrimSpace(gotBody.String()), "body")
	})

	t.Run("post/success", func(t *testing.T) {
		ts := httptest.NewServer(LoginHandler{
			Sessions: sessions.NewStore(sessions.Opts{}),
		})
		defer ts.Close()
		tsRealURL, err := url.Parse(ts.URL)
		require.NoError(t, err, "url.Parse")

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		require.NoError(t, err, "cookiejar.New")
		client := noFollow(&http.Client{Jar: jar})
		t.Log("current cookies:")
		for _, cookie := range jar.Cookies(tsRealURL) {
			t.Logf("%s: %s", cookie.Name, cookie.Value)
		}

		urlValues := make(url.Values)
		urlValues.Add("user", "jackson@jacksonargo.com")
		urlValues.Add("pass", "jackson")
		resp, err := client.PostForm(ts.URL, urlValues)
		require.NoError(t, err, "client.Get")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var gotBody bytes.Buffer
		gotBody.ReadFrom(resp.Body)
		assert.Equal(t, "welcome jackson, have a cookie!", strings.TrimSpace(gotBody.String()), "body")

	})

	t.Run("post/success+repeat", func(t *testing.T) {
		ts := httptest.NewServer(LoginHandler{
			Sessions: sessions.NewStore(sessions.Opts{}),
		})
		defer ts.Close()

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		require.NoError(t, err, "cookiejar.New")
		client := noFollow(&http.Client{Jar: jar})

		urlValues := make(url.Values)
		urlValues.Add("user", "jackson@jacksonargo.com")
		urlValues.Add("pass", "jackson")
		resp, err := client.PostForm(ts.URL, urlValues)
		require.NoError(t, err, "client.PostForm")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var gotBody bytes.Buffer
		gotBody.ReadFrom(resp.Body)
		assert.Equal(t, "welcome jackson, have a cookie!", strings.TrimSpace(gotBody.String()), "body")

		resp, err = client.Get(ts.URL)
		require.NoError(t, err, "client.Get")
		assert.Equal(t, http.StatusFound, resp.StatusCode)
		assert.Equal(t, "/", resp.Header.Get("Location"), "location")
	})
}

func noFollow(client *http.Client) *http.Client {
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return client
}

func TestSession_RenderCookie(t *testing.T) {
	secret := Secret{0x560febda7eae12b8, 0xc0cecc7851ca8906, 0x2623d26de389ebcb, 0x5a3097fc6ef622a1}
	var dest http.Cookie
	session := Session{
		ID:      0x7b7cc95133c4265d,
		Expires: time.Unix(0, 0x1607717a7c5d32e1),
	}
	session.RenderCookie(secret, &dest)
	wantCookieExpires := session.Expires
	wantCookieValue := "efc691b4e7b839e1-The-b5d412436dc4119b-Earth-b2e0dbf93cf52bb7-Is-bcf63f4fdd4adf89-Flat-7b7cc95133c4265d1607717a7c5d32e1"
	assert.Equal(t, wantCookieExpires, dest.Expires, "expires")
	assert.Equal(t, wantCookieValue, dest.Value, "value")
}

func TestSession_ReadCookie(t *testing.T) {
	secret := Secret{0x560febda7eae12b8, 0xc0cecc7851ca8906, 0x2623d26de389ebcb, 0x5a3097fc6ef622a1}
	src := http.Cookie{
		Value: "efc691b4e7b839e1-The-b5d412436dc4119b-Earth-b2e0dbf93cf52bb7-Is-bcf63f4fdd4adf89-Flat-7b7cc95133c4265d1607717a7c5d32e1",
	}
	wantSession := Session{
		ID:      0x7b7cc95133c4265d,
		Expires: time.Unix(0, 0x1607717a7c5d32e1),
	}

	var gotSession Session
	require.NoError(t, gotSession.ReadCookie(secret, &src), "Read()")
	assert.Equal(t, wantSession, gotSession, "session")
}
