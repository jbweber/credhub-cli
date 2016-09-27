package commands

import (
	"fmt"

	"net/url"

	"github.com/pivotal-cf/credhub-cli/actions"
	"github.com/pivotal-cf/credhub-cli/client"
	"github.com/pivotal-cf/credhub-cli/config"
	. "github.com/pivotal-cf/credhub-cli/util"
)

type ApiCommand struct {
	Server            ApiPositionalArgs `positional-args:"yes"`
	ServerFlagUrl     string            `short:"s" long:"server" description:"URI of API server to target"`
	SkipTlsValidation bool              `long:"skip-tls-validation" description:"Skip certificate validation of the API endpoint. Not recommended!"`
}

type ApiPositionalArgs struct {
	ServerUrl string `positional-arg-name:"SERVER" description:"URI of API server to target"`
}

func (cmd ApiCommand) Execute([]string) error {
	cfg := config.ReadConfig()
	serverUrl := targetUrl(cmd)

	if serverUrl == "" {
		fmt.Println(cfg.ApiURL)
	} else {
		existingCfg := cfg
		err := GetApiInfo(&cfg, serverUrl, cmd.SkipTlsValidation)
		if err != nil {
			return err
		}

		fmt.Println("Setting the target url:", cfg.ApiURL)

		if existingCfg.AuthURL != cfg.AuthURL {
			RevokeTokenIfNecessary(existingCfg)
			cfg = MarkTokensAsRevokedInConfig(cfg)
		}
		config.WriteConfig(cfg)
	}

	return nil
}

func GetApiInfo(cfg *config.Config, serverUrl string, skipTlsValidation bool) error {
	serverUrl = AddDefaultSchemeIfNecessary(serverUrl)
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	cfg.ApiURL = parsedUrl.String()

	cfg.InsecureSkipVerify = skipTlsValidation
	cmInfo, err := actions.NewInfo(client.NewHttpClient(*cfg), *cfg).GetServerInfo()
	if err != nil {
		return err
	}
	cfg.AuthURL = cmInfo.AuthServer.Url

	if parsedUrl.Scheme != "https" {
		Warning("Warning: Insecure HTTP API detected. Data sent to this API could be intercepted" +
			" in transit by third parties. Secure HTTPS API endpoints are recommended.")
	} else {
		if skipTlsValidation {
			Warning("Warning: The targeted TLS certificate has not been verified for this connection.")
		}
	}

	return nil
}

func targetUrl(cmd ApiCommand) string {
	if cmd.Server.ServerUrl != "" {
		return cmd.Server.ServerUrl
	} else {
		return cmd.ServerFlagUrl
	}
}
