package api

import (
	"github.com/sirupsen/logrus"
	"github.com/virtual-vgo/vvgo/pkg/login"
	"github.com/virtual-vgo/vvgo/pkg/tracing"
	"net/http"
	"time"
)

// PasswordLoginHandler authenticates requests using form values user and pass and a static map of valid combinations.
// If the user pass combo exists in the map, then a login cookie with the mapped roles is create and sent in the response.
type PasswordLoginHandler struct {
	Sessions *login.Store

	// Logins is a map of login user and pass to a slice of roles for that login.
	Logins map[[2]string][]login.Role
}

func (x PasswordLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "password_login")
	defer span.Send()

	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var identity login.Identity
	if err := x.Sessions.ReadSessionFromRequest(ctx, r, &identity); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	user := r.FormValue("user")
	pass := r.FormValue("pass")

	gotRoles, ok := x.Logins[[2]string{user, pass}]
	if !ok {
		logger.WithFields(logrus.Fields{
			"user": user,
		}).Error("password authentication failed")
		unauthorized(w)
		return
	}

	identity = login.Identity{
		Kind:  login.KindPassword,
		Roles: gotRoles,
	}

	cookie, err := x.Sessions.NewCookie(ctx, &identity, LoginCookieDuration)
	if err != nil {
		logger.WithError(err).Error("store.NewCookie() failed")
		internalServerError(w)
		return
	}

	http.SetCookie(w, cookie)
	logger.WithFields(logrus.Fields{
		"identity": identity.Kind,
		"roles":    identity.Roles,
	}).Info("authorization succeeded")
	http.Redirect(w, r, "/", http.StatusFound)
}

// LogoutHandler deletes the login session from the incoming request, if it exists.
type LogoutHandler struct {
	Sessions *login.Store
}

func (x LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "logout_handler")
	defer span.Send()

	if err := x.Sessions.DeleteSessionFromRequest(ctx, r); err != nil {
		logger.WithError(err).Error("x.Sessions.DeleteSessionFromRequest failed")
		internalServerError(w)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}


type OAuthState struct {
	Namespace   string
	RedirectURL string
}

func (x OAuthState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "discord_login")
	defer span.Send()

	statusBytes := make([]byte, 32)
	rand.Read(statusBytes)
	state := strconv.FormatUint(binary.BigEndian.Uint64(statusBytes[:16]), 16)
	value := strconv.FormatUint(binary.BigEndian.Uint64(statusBytes[16:]), 16)
	if err := redis.Do(ctx, redis.Cmd(nil, "SETEX", x.Namespace+":oauth_state:"+state, "300", "")); err != nil {
		logger.WithError(err).Error("redis.Do() failed")
		internalServerError(w)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "vvgo-discord-pre",
		Value:    value,
		Path:     "/login/discord",
		Domain:   "",
		Expires:  time.Now().Add(300 * time.Second),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, x.RedirectURL+`&state=`+state, http.StatusFound)
}

