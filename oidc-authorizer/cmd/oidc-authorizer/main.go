package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/osbuild/community-gateway/oidc-authorizer/internal/config"
	"github.com/osbuild/community-gateway/oidc-authorizer/pkg/identity"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func main() {
	conf := config.Config{}
	err := config.LoadConfigFromEnv(&conf)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, conf.Provider)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Info(r.Header)
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			http.Error(w, "Authorization header is empty", http.StatusBadRequest)
			return
		}

		tokenData := strings.Split(bearerToken, " ")
		if len(tokenData) != 2 {
			http.Error(w, "Authorization header is not a properly formatted Bearer token", http.StatusBadRequest)
			return
		}

		if tokenData[0] != "Bearer" {
			http.Error(w, "Authorization header is not a properly formatted Bearer token", http.StatusBadRequest)
			return
		}

		oauth2Token := &oauth2.Token{
			AccessToken: tokenData[1],
			TokenType:   "Bearer",
		}
		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
		if err != nil {
			logrus.Errorf("Failed to fetch user info: %v", err)
			http.Error(w, fmt.Sprintf("Failed to get userinfo: %v", err), http.StatusInternalServerError)
			return
		}

		if userInfo.Subject == "" {
			logrus.Warning("User with empty subject")
			http.Error(w, "UserInfo header has empty subject", http.StatusInternalServerError)
			return
		}

		idHeader := identity.Header{
			User: userInfo.Subject,
		}
		idHB64, err := idHeader.Base64()
		if err != nil {
			logrus.Errorf("Failed to create identity header: %v", err)
			http.Error(w, fmt.Sprintf("Failed to create identity header: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("x-fedora-identity", idHB64)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logrus.Info("Listening")
	logrus.Fatal(http.ListenAndServe("127.0.0.1:5556", nil))
}
