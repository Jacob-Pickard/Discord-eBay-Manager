package webhook

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"ebaymanager-bot/internal/ebay"

	"github.com/bwmarrin/discordgo"
)

// OAuthCallback represents a pending OAuth authorization
type OAuthCallback struct {
	State       string
	Discord     *discordgo.Session
	Interaction *discordgo.Interaction
	CreatedAt   time.Time
}

var (
	oauthCallbacks = make(map[string]*OAuthCallback)
	callbacksMutex sync.RWMutex
	ebayClient     *ebay.Client // eBay client for token exchange
)

// SetupOAuthHandlers adds OAuth callback endpoints to the webhook server
func (s *Server) SetupOAuthHandlers() {
	http.HandleFunc("/webhook/oauth/callback", s.handleOAuthCallback)
	http.HandleFunc("/webhook/oauth/declined", s.handleOAuthDeclined)
	log.Println("📍 OAuth callback endpoints registered")
}

// SetEbayClient sets the eBay client for token exchange
func SetEbayClient(client *ebay.Client) {
	ebayClient = client
	log.Println("✅ eBay client configured for OAuth")
}

// RegisterOAuthCallback registers a pending OAuth authorization
func RegisterOAuthCallback(state string, discord *discordgo.Session, interaction *discordgo.Interaction) {
	callbacksMutex.Lock()
	defer callbacksMutex.Unlock()

	oauthCallbacks[state] = &OAuthCallback{
		State:       state,
		Discord:     discord,
		Interaction: interaction,
		CreatedAt:   time.Now(),
	}

	log.Printf("📝 Registered OAuth callback for state: %s", state)

	// Clean up old callbacks
	go cleanupOldCallbacks()
}

// RegisterOAuthCallback implements the WebhookServer interface for Server
func (s *Server) RegisterOAuthCallback(state string, discord *discordgo.Session, interaction *discordgo.Interaction) {
	RegisterOAuthCallback(state, discord, interaction)
}

// cleanupOldCallbacks removes OAuth callbacks older than 10 minutes
func cleanupOldCallbacks() {
	callbacksMutex.Lock()
	defer callbacksMutex.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	for state, callback := range oauthCallbacks {
		if callback.CreatedAt.Before(cutoff) {
			delete(oauthCallbacks, state)
			log.Printf("🗑️ Cleaned up expired OAuth callback: %s", state)
		}
	}
}

// handleOAuthCallback processes OAuth authorization callbacks from eBay
func (s *Server) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	log.Printf("📨 OAuth callback received - State: %s, Has code: %v, Error: %s", state, code != "", errorParam)

	if errorParam != "" {
		errorDesc := r.URL.Query().Get("error_description")
		sendOAuthError(state, fmt.Sprintf("%s: %s", errorParam, errorDesc))

		html := `<!DOCTYPE html>
<html><head><title>Authorization Failed</title><style>body{font-family:Arial;max-width:600px;margin:50px auto;padding:20px}.error{color:red}</style></head>
<body><h1 class="error">❌ Authorization Failed</h1><p>` + errorDesc + `</p><p>You can close this window and try again in Discord with <code>/ebay-authorize</code></p></body></html>`

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}

	if code == "" {
		sendOAuthError(state, "No authorization code received")
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Show success page immediately
	html := `<!DOCTYPE html>
<html><head><title>Authorization Successful</title><style>body{font-family:Arial;max-width:600px;margin:50px auto;padding:20px}.success{color:green}</style></head>
<body><h1 class="success">✅ Authorization Successful!</h1><p>Your eBay account has been connected successfully.</p><p><strong>You can now close this window and return to Discord.</strong></p></body></html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))

	// Process token exchange in background
	go processOAuthToken(state, code)
}

// handleOAuthDeclined handles when user declines authorization
func (s *Server) handleOAuthDeclined(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")

	log.Printf("❌ OAuth declined - State: %s", state)
	sendOAuthError(state, "Authorization declined by user")

	html := `<!DOCTYPE html>
<html><head><title>Authorization Declined</title><style>body{font-family:Arial;max-width:600px;margin:50px auto;padding:20px}.error{color:red}</style></head>
<body><h1 class="error">❌ Authorization Declined</h1><p>You declined the authorization request.</p><p>You can close this window. If you want to try again, use <code>/ebay-authorize</code> in Discord.</p></body></html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// processOAuthToken exchanges the code for tokens and notifies Discord
func processOAuthToken(state, code string) {
	log.Printf("🔄 Processing OAuth token exchange for state: %s", state)

	callbacksMutex.RLock()
	callback, exists := oauthCallbacks[state]
	callbacksMutex.RUnlock()

	if !exists || callback == nil {
		log.Printf("⚠️ No callback found for state: %s", state)
		return
	}

	// Check if eBay client is configured
	if ebayClient == nil {
		log.Println("❌ eBay client is nil")
		callback.Discord.FollowupMessageCreate(callback.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ **Server configuration error**\n\nEbay client not properly configured. Contact administrator.",
		})
		return
	}

	// Exchange code for token using ebayClient
	_, err := ebayClient.ExchangeCodeForToken(code)
	if err != nil {
		log.Printf("❌ Failed to exchange code for token: %v", err)
		callback.Discord.FollowupMessageCreate(callback.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("❌ **Failed to get access token:** %v\n\nTry `/ebay-authorize` again or use `/ebay-code` for manual entry.", err),
		})

		// Remove callback
		callbacksMutex.Lock()
		delete(oauthCallbacks, state)
		callbacksMutex.Unlock()
		return
	}

	// Success! Notify Discord
	log.Printf("✅ OAuth tokens obtained successfully for state: %s", state)
	callback.Discord.FollowupMessageCreate(callback.Interaction, true, &discordgo.WebhookParams{
		Content: "✅ **Authorization Successful!**\n\nYour eBay account has been connected.\nAccess token and refresh token have been saved.\n\n🎉 You can now use all eBay commands!\n\n💡 The bot will automatically refresh your token every 90 minutes.",
	})

	// Clean up callback
	callbacksMutex.Lock()
	delete(oauthCallbacks, state)
	callbacksMutex.Unlock()

	log.Printf("📨 Notified Discord and cleaned up OAuth callback for state: %s", state)
}

// sendOAuthError sends an error message to Discord for a pending OAuth flow
func sendOAuthError(state, errorMsg string) {
	callbacksMutex.RLock()
	callback, exists := oauthCallbacks[state]
	callbacksMutex.RUnlock()

	if !exists || callback == nil {
		log.Printf("⚠️ No callback found for state: %s", state)
		return
	}

	callback.Discord.FollowupMessageCreate(callback.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("❌ **Authorization failed:** %s\n\n💡 Try again with `/ebay-authorize`", errorMsg),
	})

	// Remove callback
	callbacksMutex.Lock()
	delete(oauthCallbacks, state)
	callbacksMutex.Unlock()
}
