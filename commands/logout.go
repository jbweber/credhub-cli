package commands

import (
	"fmt"

	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	"github.com/pivotal-cf/credhub-cli/repositories"
)

type LogoutCommand struct {
}

func (cmd LogoutCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	RevokeTokenIfNecessary(cfg)
	cfg = MarkTokensAsRevokedInConfig(cfg)
	config.WriteConfig(cfg)
	fmt.Println("Logout Successful")
	return nil
}

func RevokeTokenIfNecessary(cfg config.Config) {
	if cfg.RefreshToken != "" && cfg.RefreshToken != "revoked" {
		authRepository := repositories.NewAuthRepository(client.NewHttpClient(cfg), false)
		request, err := client.NewTokenRevocationRequest(cfg)
		if err == nil {
			authRepository.SendRequest(request, "logout")
		}
	}
}

func MarkTokensAsRevokedInConfig(cfg config.Config) config.Config {
	cfg.AccessToken = "revoked"
	cfg.RefreshToken = "revoked"
	return cfg
}
