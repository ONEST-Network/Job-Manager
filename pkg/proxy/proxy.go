package proxy

import (
	"os"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/config"
	"github.com/sirupsen/logrus"
)

func SetProxyENVs() {
	if config.Config.HttpProxy != "" {
		if err := os.Setenv("HTTP_PROXY", config.Config.HttpProxy); err != nil {
			logrus.Fatalf("Failed to set HTTP_PROXY env to '%s', err: %v", config.Config.HttpProxy, err)
		}
	}

	if config.Config.HttpsProxy != "" {
		if err := os.Setenv("HTTPS_PROXY", config.Config.HttpsProxy); err != nil {
			logrus.Fatalf("Failed to set HTTPS_PROXY env to '%s', err: %v", config.Config.HttpsProxy, err)
		}
	}

	if config.Config.NoProxy != "" {
		if err := os.Setenv("NO_PROXY", config.Config.NoProxy); err != nil {
			logrus.Fatalf("Failed to set NO_PROXY env to '%s', err: %v", config.Config.NoProxy, err)
		}
	}
}
