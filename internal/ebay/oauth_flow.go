package ebay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// OAuthServer handles the OAuth callback
type OAuthServer struct {
	client       *Client
	server       *http.Server
	tokenChan    chan *TokenResponse
	errorChan    chan error
	authURL      string
	state        string
}

// NewOAuthServer creates a new OAuth callback server
func NewOAuthServer(client *Client) *OAuthServer {
	return &OAuthServer{
		client:    client,
		tokenChan: make(chan *TokenResponse, 1),
		errorChan: make(chan error, 1),
		state:     fmt.Sprintf("state_%d", time.Now().Unix()),
	}
}

// StartAuthFlow generates the authorization URL and starts the callback server
func (s *OAuthServer) StartAuthFlow() (string, error) {
	// Generate authorization URL
	s.authURL = s.client.GetUserAuthorizationURL(s.state)

	// Start HTTP server to handle callback
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", s.handleCallback)
	mux.HandleFunc("/", s.handleRoot)

	s.server = &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("OAuth server error: %v", err)
		}
	}()

	log.Println("OAuth server started on http://localhost:3000")
	return s.authURL, nil
}

// handleRoot shows a simple landing page
func (s *OAuthServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>eBay OAuth</title>
		<style>
			body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
			.success { color: green; }
			.error { color: red; }
			a { color: #0066c0; }
		</style>
	</head>
	<body>
		<h1>eBay Manager Bot - OAuth Setup</h1>
		<p>Click the link below to authorize the bot with your eBay account:</p>
		<p><a href="` + s.authURL + `" target="_blank">Authorize with eBay</a></p>
		<p>After authorizing, you'll be redirected back here.</p>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// handleCallback processes the OAuth callback
func (s *OAuthServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	if errorParam != "" {
		errorDesc := r.URL.Query().Get("error_description")
		s.errorChan <- fmt.Errorf("authorization failed: %s - %s", errorParam, errorDesc)
		
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Authorization Failed</title></head>
		<body>
			<h1 class="error">Authorization Failed</h1>
			<p>Error: ` + errorParam + `</p>
			<p>` + errorDesc + `</p>
			<p>You can close this window.</p>
		</body>
		</html>
		`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}

	if state != s.state {
		s.errorChan <- fmt.Errorf("invalid state parameter")
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	if code == "" {
		s.errorChan <- fmt.Errorf("no authorization code received")
		http.Error(w, "No code", http.StatusBadRequest)
		return
	}

	// Exchange code for tokens
	tokens, err := s.client.ExchangeCodeForToken(code)
	if err != nil {
		s.errorChan <- fmt.Errorf("failed to exchange code: %w", err)
		
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Token Exchange Failed</title></head>
		<body>
			<h1 class="error">Token Exchange Failed</h1>
			<p>` + err.Error() + `</p>
			<p>You can close this window and try again.</p>
		</body>
		</html>
		`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}

	s.tokenChan <- tokens

	// Save tokens to .env file
	if err := s.saveTokensToEnv(tokens); err != nil {
		log.Printf("Warning: Failed to save tokens to .env: %v", err)
	}

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Authorization Success</title>
		<style>
			body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
			.success { color: green; }
		</style>
	</head>
	<body>
		<h1 class="success">✅ Authorization Successful!</h1>
		<p>Your eBay account has been connected successfully.</p>
		<p>Access token and refresh token have been saved.</p>
		<p><strong>You can now close this window and return to Discord.</strong></p>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// WaitForTokens waits for the OAuth flow to complete
func (s *OAuthServer) WaitForTokens(timeout time.Duration) (*TokenResponse, error) {
	select {
	case tokens := <-s.tokenChan:
		s.Stop()
		return tokens, nil
	case err := <-s.errorChan:
		s.Stop()
		return nil, err
	case <-time.After(timeout):
		s.Stop()
		return nil, fmt.Errorf("authorization timeout")
	}
}

// Stop shuts down the OAuth server
func (s *OAuthServer) Stop() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}

// SaveTokensToEnv updates the .env file with new tokens (exported for manual code submission)
func (s *OAuthServer) SaveTokensToEnv(tokens *TokenResponse) error {
	return s.saveTokensToEnv(tokens)
}

// saveTokensToEnv updates the .env file with new tokens
func (s *OAuthServer) saveTokensToEnv(tokens *TokenResponse) error {
	envPath := ".env"
	
	// Read current .env file
	data, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	
	// Update token lines
	for i, line := range lines {
		if strings.HasPrefix(line, "EBAY_ACCESS_TOKEN=") {
			lines[i] = "EBAY_ACCESS_TOKEN=" + tokens.AccessToken
		} else if strings.HasPrefix(line, "EBAY_REFRESH_TOKEN=") {
			if tokens.RefreshToken != "" {
				lines[i] = "EBAY_REFRESH_TOKEN=" + tokens.RefreshToken
			}
		}
	}

	// Write back to file
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(envPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	log.Println("Tokens saved to .env file")
	return nil
}

// AutoRefreshToken automatically refreshes the access token when it expires
func (c *Client) AutoRefreshToken(refreshToken string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Refreshing eBay access token...")
		
		tokens, err := c.RefreshAccessToken(refreshToken)
		if err != nil {
			log.Printf("Failed to refresh token: %v", err)
			continue
		}

		log.Println("Access token refreshed successfully")

		// Update tokens in .env file
		envPath := ".env"
		data, err := os.ReadFile(envPath)
		if err != nil {
			log.Printf("Failed to read .env for token update: %v", err)
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		
		for i, line := range lines {
			if strings.HasPrefix(line, "EBAY_ACCESS_TOKEN=") {
				lines[i] = "EBAY_ACCESS_TOKEN=" + tokens.AccessToken
			}
			if tokens.RefreshToken != "" && strings.HasPrefix(line, "EBAY_REFRESH_TOKEN=") {
				lines[i] = "EBAY_REFRESH_TOKEN=" + tokens.RefreshToken
			}
		}

		newContent := strings.Join(lines, "\n")
		if err := os.WriteFile(envPath, []byte(newContent), 0644); err != nil {
			log.Printf("Failed to save refreshed tokens: %v", err)
		}
	}
}

// GetTokenInfo returns information about the current token
func (c *Client) GetTokenInfo() (map[string]interface{}, error) {
	endpoint := "/commerce/identity/v1/user"
	
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	if err := json.Unmarshal(respBody, &info); err != nil {
		return nil, err
	}

	return info, nil
}
