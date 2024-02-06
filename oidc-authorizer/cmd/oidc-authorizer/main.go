package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/osbuild/community-gateway/oidc-authorizer/internal/config"
	"github.com/osbuild/community-gateway/oidc-authorizer/pkg/identity"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// https://datatracker.ietf.org/doc/rfc7662/
type IntrospectionResponse struct {
	Active   bool   `json:"active"`
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	Scope    string `json:"scope"`
}

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

		var idHeader identity.Identity
		// fallback to userinfo if client id/client secret not specified
		if conf.IntrospectURL == "" {
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
			idHeader = identity.Identity{
				User: userInfo.Subject,
			}
		} else {
			//
			body := url.Values{}
			body.Set("client_id", conf.ClientID)
			body.Set("client_secret", conf.ClientSecret)

			req, err := http.NewRequest("POST", conf.IntrospectURL, strings.NewReader(body.Encode()))
			if err != nil {
				logrus.Errorf("Failed to create a new token introspection request: %v", err)
				http.Error(w, fmt.Sprintf("Failed to introspect token: %v", err), http.StatusInternalServerError)
				return
			}

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			query := req.URL.Query()
			query.Add("token", oauth2Token.AccessToken)
			req.URL.RawQuery = query.Encode()

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logrus.Errorf("Failed to create a new token introspection request: %v", err)
				http.Error(w, fmt.Sprintf("Failed to introspect token: %v", err), http.StatusInternalServerError)
				return
			}

			if resp.StatusCode == http.StatusUnauthorized {
				http.Error(w, fmt.Sprintf("Introspection returned 401: %v", err), http.StatusUnauthorized)
				return
			}

			if resp.StatusCode != http.StatusOK {
				logrus.Errorf("Introspection returned unexected status : %d", resp.StatusCode)
				http.Error(w, fmt.Sprintf("Failed to introspect token: %v", err), http.StatusInternalServerError)
				return
			}

			var introResp IntrospectionResponse
			err = json.NewDecoder(resp.Body).Decode(&introResp)
			if err != nil {
				logrus.Errorf("Failed to parse token introspection response: %v", err)
				http.Error(w, fmt.Sprintf("Failed to introspect token: %v", err), http.StatusInternalServerError)
				return
			}

			idHeader = identity.Identity{
				User: introResp.Username,
			}
		}

		idHB64, err := idHeader.Base64()
		if err != nil {
			logrus.Errorf("Failed to create identity header: %v", err)
			http.Error(w, fmt.Sprintf("Failed to create identity header: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add(identity.FedoraIDHeader, idHB64)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logrus.Info("Listening")
	logrus.Fatal(http.ListenAndServe(":5556", nil))
}
