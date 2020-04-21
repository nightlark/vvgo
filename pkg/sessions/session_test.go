package sessions

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestSession_String(t *testing.T) {
	secret := Secret{0x560febda7eae12b8, 0xc0cecc7851ca8906, 0x2623d26de389ebcb, 0x5a3097fc6ef622a1}
	session := Session{
		ID:      0x7b7cc95133c4265d,
		Expires: time.Unix(0, 0x1607717a7c5d32e1),
	}
	got := session.String(secret)
	wantCookieValue := "efc691b4e7b839e1-V-b5d412436dc4119b-V-b2e0dbf93cf52bb7-G-bcf63f4fdd4adf89-O-7b7cc95133c4265d1607717a7c5d32e1"
	assert.Equal(t, wantCookieValue, got, "value")
}

func TestSession_ReadCookie(t *testing.T) {
	secret := Secret{0x560febda7eae12b8, 0xc0cecc7851ca8906, 0x2623d26de389ebcb, 0x5a3097fc6ef622a1}
	src := http.Cookie{
		Value: "efc691b4e7b839e1-V-b5d412436dc4119b-V-b2e0dbf93cf52bb7-G-bcf63f4fdd4adf89-O-7b7cc95133c4265d1607717a7c5d32e1",
	}
	wantSession := Session{
		ID:      0x7b7cc95133c4265d,
		Expires: time.Unix(0, 0x1607717a7c5d32e1),
	}

	var gotSession Session
	require.NoError(t, gotSession.ReadCookie(secret, &src), "Read()")
	assert.Equal(t, wantSession, gotSession, "session")
}

func TestSecret(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		token := NewSecret()
		assert.NoError(t, token.Validate(), "validate")
		t.Logf("new token: %s", token.String())
		assert.NotEqual(t, token.String(), NewSecret().String())
	})
	t.Run("decode", func(t *testing.T) {
		arg := "196ddf804c7666d4-8d32ff4a91a530bc-c5c7cde4a26096ad-67758135226bfb2e"
		expected := Secret{0x196ddf804c7666d4, 0x8d32ff4a91a530bc, 0xc5c7cde4a26096ad, 0x67758135226bfb2e}
		got, _ := DecodeSecret(arg)
		assert.Equal(t, expected, got)
	})
	t.Run("string", func(t *testing.T) {
		expected := "196ddf804c7666d4-8d32ff4a91a530bc-c5c7cde4a26096ad-67758135226bfb2e"
		arg := Secret{0x196ddf804c7666d4, 0x8d32ff4a91a530bc, 0xc5c7cde4a26096ad, 0x67758135226bfb2e}
		got := arg.String()
		assert.Equal(t, expected, got)
	})
	t.Run("validate/success", func(t *testing.T) {
		arg := Secret{0x196ddf804c7666d4, 0x8d32ff4a91a530bc, 0xc5c7cde4a26096ad, 0x67758135226bfb2e}
		assert.NoError(t, arg.Validate())
	})
	t.Run("validate/fail", func(t *testing.T) {
		arg := Secret{0, 0x8d32ff4a91a530bc, 0xc5c7cde4a26096ad, 0x67758135226bfb2e}
		assert.Equal(t, ErrInvalidSecret, arg.Validate())
	})
}
