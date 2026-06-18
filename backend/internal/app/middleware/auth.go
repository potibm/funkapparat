package middleware

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/potibm/funkapparat/internal/app/domain"
	sloggin "github.com/samber/slog-gin"
)

func AuthMiddleware(ctx context.Context, issuerURL, clientID string, skipTLSVerify bool) (gin.HandlerFunc, error) {
	// 1. HTTP client with optional TLS verification
	const oidcHTTPTimeout = 10 * time.Second

	baseTransport, _ := http.DefaultTransport.(*http.Transport)
	transport := baseTransport.Clone()
	client := &http.Client{
		Transport: transport,
		Timeout:   oidcHTTPTimeout,
	}

	if skipTLSVerify {
		// #nosec G402 -- for local dev environments
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // NOSONAR
	}

	// 2. Add the custom HTTP client to the OIDC context
	setupCtx, cancel := context.WithTimeout(ctx, oidcHTTPTimeout)
	defer cancel()

	oidcCtx := oidc.ClientContext(setupCtx, client)

	// 3. Initialize the OIDC Provider
	provider, err := oidc.NewProvider(oidcCtx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("error initializing the OIDC Provider: %w", err)
	}

	// 4. Configure the Verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	// 5. Return the Gin middleware
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]

		requestCtx := oidc.ClientContext(c.Request.Context(), client)

		idToken, err := verifier.Verify(requestCtx, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userID := idToken.Subject

		c.Set("userID", userID)
		sloggin.AddCustomAttributes(c, slog.String("user_id", userID))

		ctxWithUser := context.WithValue(c.Request.Context(), domain.UserIDKey, userID)
		c.Request = c.Request.WithContext(ctxWithUser)

		c.Next()
	}, nil
}
