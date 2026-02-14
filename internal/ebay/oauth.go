package ebay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	sandboxTokenURL    = "https://api.sandbox.ebay.com/identity/v1/oauth2/token"
	productionTokenURL = "https://api.ebay.com/identity/v1/oauth2/token"
)

// TokenResponse represents the OAuth token response from eBay
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope,omitempty"` // Space-separated list of granted scopes
}

// GetApplicationToken gets an application token using client credentials
// This is useful for public API calls that don't require user authorization
func (c *Client) GetApplicationToken() (*TokenResponse, error) {
	tokenURL := sandboxTokenURL
	if c.config.Environment == "PRODUCTION" {
		tokenURL = productionTokenURL
	}

	// Create the credentials for Basic Auth
	credentials := c.config.AppID + ":" + c.config.CertID
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	// Prepare the request body
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "https://api.ebay.com/oauth/api_scope")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshAccessToken refreshes an expired access token using the refresh token
func (c *Client) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	tokenURL := sandboxTokenURL
	if c.config.Environment == "PRODUCTION" {
		tokenURL = productionTokenURL
	}

	credentials := c.config.AppID + ":" + c.config.CertID
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("scope", "https://api.ebay.com/oauth/api_scope https://api.ebay.com/oauth/api_scope/sell.inventory https://api.ebay.com/oauth/api_scope/sell.fulfillment https://api.ebay.com/oauth/api_scope/sell.account https://api.ebay.com/oauth/api_scope/sell.finances")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	// Update the client's access token
	c.config.AccessToken = tokenResp.AccessToken

	return &tokenResp, nil
}

// GetUserAuthorizationURL generates the URL for user authorization
// Users need to visit this URL to grant your application access to their eBay account
// Note: eBay requires the RuName as the redirect_uri parameter, not the actual callback URL
func (c *Client) GetUserAuthorizationURL(state string) string {
	authURL := "https://auth.sandbox.ebay.com/oauth2/authorize"
	if c.config.Environment == "PRODUCTION" {
		authURL = "https://auth.ebay.com/oauth2/authorize"
	}

	params := url.Values{}
	params.Set("client_id", c.config.AppID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", c.config.RedirectURI) // This should be the RuName
	params.Set("scope", "https://api.ebay.com/oauth/api_scope https://api.ebay.com/oauth/api_scope/sell.inventory https://api.ebay.com/oauth/api_scope/sell.fulfillment https://api.ebay.com/oauth/api_scope/sell.account https://api.ebay.com/oauth/api_scope/sell.finances")
	if state != "" {
		params.Set("state", state)
	}

	return authURL + "?" + params.Encode()
}

// ExchangeCodeForToken exchanges an authorization code for access and refresh tokens
func (c *Client) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	tokenURL := sandboxTokenURL
	if c.config.Environment == "PRODUCTION" {
		tokenURL = productionTokenURL
	}

	credentials := c.config.AppID + ":" + c.config.CertID
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.config.RedirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token exchange request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token exchange response: %w", err)
	}

	// Log full response for debugging OAuth issues
	log.Printf("[DEBUG] Token Exchange Response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token exchange response: %w", err)
	}

	// Log the granted scopes for debugging
	log.Printf("🔑 OAuth tokens obtained - Granted scopes: %s", tokenResp.Scope)

	// Update the client's tokens
	c.config.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		c.config.RefreshToken = tokenResp.RefreshToken
	}

	return &tokenResp, nil
}

// TokenManager helps manage token expiration and refresh
type TokenManager struct {
	client        *Client
	accessToken   string
	refreshToken  string
	expiresAt     time.Time
	refreshBefore time.Duration // Refresh token this duration before expiry
}

// NewTokenManager creates a new token manager
func NewTokenManager(client *Client, accessToken, refreshToken string, expiresIn int) *TokenManager {
	return &TokenManager{
		client:        client,
		accessToken:   accessToken,
		refreshToken:  refreshToken,
		expiresAt:     time.Now().Add(time.Duration(expiresIn) * time.Second),
		refreshBefore: 5 * time.Minute, // Refresh 5 minutes before expiry
	}
}

// GetValidToken returns a valid access token, refreshing if necessary
func (tm *TokenManager) GetValidToken() (string, error) {
	// Check if token needs refresh
	if time.Now().Add(tm.refreshBefore).After(tm.expiresAt) {
		if tm.refreshToken == "" {
			return "", fmt.Errorf("token expired and no refresh token available")
		}

		// Refresh the token
		tokenResp, err := tm.client.RefreshAccessToken(tm.refreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}

		tm.accessToken = tokenResp.AccessToken
		tm.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

		// Update refresh token if a new one was provided
		if tokenResp.RefreshToken != "" {
			tm.refreshToken = tokenResp.RefreshToken
		}
	}

	return tm.accessToken, nil
}
