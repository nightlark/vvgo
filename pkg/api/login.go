package api

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"github.com/virtual-vgo/vvgo/pkg/discord"
	"github.com/virtual-vgo/vvgo/pkg/redis"
	"github.com/virtual-vgo/vvgo/pkg/tracing"
	"net/http"
	"strconv"
	"time"
)

type DiscordLoginHandler struct {
	GuildID        discord.GuildID
	RoleVVGOMember string
	Discord        *discord.Client
}

var ErrNotAMember = errors.New("not a member")

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

func (x DiscordLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "discord_oauth_handler")
	defer span.Send()

	handleError := func(err error) bool {
		if err != nil {
			tracing.AddError(ctx, err)
			logger.WithError(err).Error("discord authentication failed")
			unauthorized(w)
			return false
		}
		return true
	}

	// read the state param
	state := r.FormValue("state")
	if state == "" {
		handleError(errors.New("no state param"))
		return
	}

	// check if it exists in redis
	redis.Do(ctx, redis.Cmd(gotValue, "GET", x.Namespace+":oauth_state:"+state))

	// check against the cookie value
	cookie, err := r.Cookie("vvgo-discord-pre")
	if ok := handleError(err); !ok {
		return
	}



	// get an oauth token from discord
	code := r.FormValue("code")
	oauthToken, err := x.Discord.QueryOAuth(ctx, code)
	if ok := handleError(err); !ok {
		return
	}

	// get the user id
	discordUser, err := x.Discord.QueryIdentity(ctx, oauthToken)
	if ok := handleError(err); !ok {
		return
	}

	// check if this user is in our guild
	guildMember, err := x.Discord.QueryGuildMember(ctx, x.GuildID, discordUser.ID)
	if ok := handleError(err); !ok {
		return
	}

	// check that they have the member role
	var ok bool
	for _, role := range guildMember.Roles {
		if role == x.RoleVVGOMember {
			ok = true
			break
		}
	}
	if !ok {
		handleError(ErrNotAMember)
		return
	}
	w.Write([]byte("authorized"))
}
