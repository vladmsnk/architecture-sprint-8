package main

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"os"
	"reports-api/internal/app/reports"
	uc "reports-api/internal/usecase/reports"
)

var (
	keycloakAddrEnv  = "KEYCLOAK_ADDR"
	keycloakRealmEnv = "KEYCLOAK_REALM"
	apiAddrEnv       = "APP_ADDR"
)

func main() {
	ctx := context.Background()

	host := os.Getenv(apiAddrEnv)

	keySet, err := extractJSONWebKeySetFromKeycloak(ctx)
	if err != nil {
		log.Fatal(err)
	}

	getter := uc.NewGetter()

	e := reports.New(keySet, getter)

	if err := e.Run(host); err != nil {
		log.Fatal(err)
	}
}

func extractJSONWebKeySetFromKeycloak(ctx context.Context) (jwk.Set, error) {
	keycloakAddr := os.Getenv(keycloakAddrEnv)
	keycloakRealm := os.Getenv(keycloakRealmEnv)

	jwkURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", keycloakAddr, keycloakRealm)

	keySet, err := jwk.Fetch(ctx, jwkURL)
	if err != nil {
		return nil, fmt.Errorf("jwk.Fetch: %w", err)
	}

	return keySet, err
}
