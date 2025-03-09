package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Marattttt/new_new_portfolio/auth/pkg/models"
	"golang.org/x/oauth2"
)

const (
	googleUserInfoUrl = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type googleAuth struct {
	conf   oauth2.Config
	logger *slog.Logger
}

func (g googleAuth) generateLoginURL(state string) string {
	url := g.conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url
}

func (g googleAuth) getUserInfo(ctx context.Context, code string) (*models.User, error) {
	token, err := g.conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging: %w", err)
	}


	g.logger.Debug("Finished code for token exchange", slog.Any("token", token))

	client := g.conf.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("getting user info: %w", err)
	}
	defer resp.Body.Close()

	userInfo := make(map[string]any)

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	g.logger.Debug("Finished requesting user info", slog.Any("userInfo", userInfo))

	user, err := g.parseUserInfo(token, userInfo)

	return user, nil
}

func (g googleAuth) parseUserInfo(
	token *oauth2.Token,
	userInfo map[string]any,
) (
	*models.User, error,
) {
	user := &models.User{
		Email: userInfo["email"].(string),
		Name:  userInfo["name"].(string),
	}

	if user.Name == "" || user.Email == "" {
		return nil, fmt.Errorf("Invalid email '%s' or username '%s'", user.Email, user.Name)
	}

	user.Google = &models.GoogleOauth{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	g.logger.Debug("Decoded data",
		slog.Any("user", *user),
		slog.Any("googleAuth", *user.Google),
	)

	if token.Expiry == (time.Time{}) {
		user.Google.Expiry = time.Now().Add(time.Hour * 24 * 365)
	} else {
		user.Google.Expiry = token.Expiry
	}

	return user, nil
}
