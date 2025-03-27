package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Marattttt/new_new_portfolio/auth/internal/config"
	"github.com/Marattttt/new_new_portfolio/auth/pkg/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthResourceServer int

const GoogleResourceServer OAuthResourceServer = iota

type StateProvider interface {
	GenerateState(OAuthResourceServer) string
	InvalidateState(state string) error
}

type Service struct {
	logger *slog.Logger
	state  StateProvider
	google googleAuth
}

func NewService(conf *config.OAuthConfig, sp StateProvider, logger *slog.Logger) *Service {
	return &Service{
		logger: logger,
		state:  sp,
		google: googleAuth{
			logger: logger.With(slog.String("resourceServer", "google")),
			conf: oauth2.Config{
				ClientID:     conf.Google.ClientID,
				ClientSecret: conf.Google.ClientSecret,
				RedirectURL:  conf.Google.Redirect,
				Scopes:       conf.Google.Scopes,
				Endpoint:     google.Endpoint,
			},
		},
	}
}

func (s Service) GenerateLoginURL(resource OAuthResourceServer) (string) {
	switch resource {
	case GoogleResourceServer:
		s.logger.Debug("Using google resource provider")
		state := s.state.GenerateState(GoogleResourceServer)
		url := s.google.generateLoginURL(state)
		return  url

	default:
		panic(fmt.Sprintf("no such OAuthResourceServer: %d", resource))
	}
}

func (s Service) HandleCallback(ctx context.Context, resource OAuthResourceServer, code string) (*models.User, error) {
	switch resource {
	case GoogleResourceServer:
		u, err := s.google.getUserInfo(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("google: %w", err)
		}
		return u, nil

	default:
		panic(fmt.Sprintf("no such OAuthResourceServer: %d", resource))
	}
}
