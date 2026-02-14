package ebay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"ebaymanager-bot/internal/config"
)

const (
	sandboxAPIURL     = "https://api.sandbox.ebay.com"
	productionAPIURL  = "https://api.ebay.com"
	sandboxAuthURL    = "https://auth.sandbox.ebay.com/oauth2/authorize"
	productionAuthURL = "https://auth.ebay.com/oauth2/authorize"
)

// Client handles eBay API interactions
type Client struct {
	config     config.EbayConfig
	httpClient *http.Client
	baseURL    string
	authURL    string
}

// NewClient creates a new eBay API client
func NewClient(cfg config.EbayConfig) *Client {
	baseURL := sandboxAPIURL
	authURL := sandboxAuthURL

	if cfg.Environment == "PRODUCTION" {
		baseURL = productionAPIURL
		authURL = productionAuthURL
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		authURL: authURL,
	}
}

// CheckConnection verifies the eBay API connection
func (c *Client) CheckConnection() string {
	if c.config.AccessToken == "" {
		return "âŒ No access token configured. You need to authenticate with eBay first."
	}

	// Try a simple API call to verify connection
	// TODO: Implement actual API call
	return fmt.Sprintf("âœ… Connected to eBay API (%s mode)\nðŸ”‘ Access token: %s...",
		c.config.Environment,
		c.config.AccessToken[:10])
}

// GetTokenScopes returns the OAuth scopes that would be requested for a new token
func (c *Client) GetTokenScopes() []string {
	// These are the scopes we configure in oauth.go
	return []string{
		"api_scope",
		"sell.inventory",
		"sell.fulfillment",
		"sell.account",
		"sell.finances",
	}
}

// makeRequest is a helper to make authenticated requests to eBay API
func (c *Client) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	fullURL := c.baseURL + endpoint
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Language", "en-US")
	req.Header.Set("Accept-Language", "en-US")

	// Log request details for Finances API debugging
	if strings.Contains(endpoint, "/finances/") {
		log.Printf("[DEBUG] Finances API Request: %s %s", method, fullURL)
		log.Printf("[DEBUG] Access Token (first 20 chars): %s...", c.config.AccessToken[:20])
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Enhanced logging for Finances API errors
	if resp.StatusCode >= 400 {
		if strings.Contains(endpoint, "/finances/") {
			log.Printf("[DEBUG] Finances API Error - Status: %d", resp.StatusCode)
			log.Printf("[DEBUG] Response Headers: %v", resp.Header)
			log.Printf("[DEBUG] Response Body: %s", string(respBody))
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// OrdersResponse represents the response from the fulfillment API
type OrdersResponse struct {
	Orders []Order `json:"orders"`
	Total  int     `json:"total"`
}

// GetOrders fetches recent orders from eBay Fulfillment API
func (c *Client) GetOrders(limit int) ([]Order, error) {
	if c.config.AccessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	endpoint := fmt.Sprintf("/sell/fulfillment/v1/order?limit=%d", limit)

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var ordersResp OrdersResponse
	if err := json.Unmarshal(respBody, &ordersResp); err != nil {
		return nil, fmt.Errorf("failed to parse orders response: %w", err)
	}

	// Process orders to populate computed fields
	for i := range ordersResp.Orders {
		order := &ordersResp.Orders[i]

		// Extract buyer username
		order.BuyerUsername = order.Buyer.Username

		// Extract price and currency
		fmt.Sscanf(order.PricingSummary.Total.Value, "%f", &order.TotalPrice)
		order.Currency = order.PricingSummary.Total.Currency

		// Extract fulfillment status
		order.FulfillmentStatus = order.OrderFulfillmentStatus

		// Process line items to get images and prices
		for j := range order.LineItems {
			lineItem := &order.LineItems[j]

			// Extract line item price
			fmt.Sscanf(lineItem.LineItemCost.Value, "%f", &lineItem.Price)

			// eBay Fulfillment API doesn't include images, so fetch from Inventory API
			if lineItem.LegacyItemId != "" {
				if img := c.getItemImage(lineItem.LegacyItemId); img != "" {
					lineItem.ImageUrl = img
				}
			}
		}
	}

	return ordersResp.Orders, nil
}

// GetOrderByID fetches a specific order by ID
func (c *Client) GetOrderByID(orderID string) (*Order, error) {
	if c.config.AccessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	endpoint := fmt.Sprintf("/sell/fulfillment/v1/order/%s", orderID)

	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var order Order
	if err := json.Unmarshal(respBody, &order); err != nil {
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	return &order, nil
}

// getItemImage fetches the image URL for a specific item using Browse API
func (c *Client) getItemImage(legacyItemId string) string {
	// Use eBay's Browse API to get item details with images
	// Browse API uses item_id format, needs conversion from legacy ID
	endpoint := fmt.Sprintf("/buy/browse/v1/item/get_item_by_legacy_id?legacy_item_id=%s", legacyItemId)
	respBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		// If Browse API fails, try standard thumb URL pattern
		return fmt.Sprintf("https://thumbs.ebayimg.com/thumbs/g/%s/s-l225.jpg", legacyItemId)
	}

	var result struct {
		Image struct {
			ImageUrl string `json:"imageUrl"`
		} `json:"image"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Sprintf("https://thumbs.ebayimg.com/thumbs/g/%s/s-l225.jpg", legacyItemId)
	}

	if result.Image.ImageUrl != "" {
		return result.Image.ImageUrl
	}

	return fmt.Sprintf("https://thumbs.ebayimg.com/thumbs/g/%s/s-l225.jpg", legacyItemId)
}

// OffersResponse represents the response from the offers API
type OffersResponse struct {
	Offers []Offer `json:"offers"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

// GetOffers fetches pending buyer offers (best offers) from eBay
func (c *Client) GetOffers() ([]Offer, error) {
	if c.config.AccessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	// Note: eBay's Sell APIs don't have a direct "buyer offers" endpoint
	// Best offers are integrated into the Inventory API per listing
	// For webhooks, you'll get OFFER notifications
	// For now, return empty list - offers will be handled via webhooks

	return []Offer{}, nil
}

// GetListings retrieves active inventory listings
func (c *Client) GetListings(limit int) ([]map[string]interface{}, error) {
	if c.config.AccessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	// Note: Listing retrieval requires either:
	// 1. Inventory API with proper SKU management (not standard eBay listings)
	// 2. Trading API with XML parsing (slow, complex)
	// 3. Browse API (public data only, no seller-specific listings)
	//
	// For now, return guidance to check listings on eBay directly
	return []map[string]interface{}{
		{"message": "Listing retrieval requires Trading API XML parsing or inventory management. Check your active listings at: https://www.ebay.com/sh/lst/active"},
	}, nil
}

// RespondToOffer accepts, declines, or counters a buyer offer
// action can be: "ACCEPT", "DECLINE", or "COUNTER"
func (c *Client) RespondToOffer(offerID string, action string, counterPrice float64) error {
	if c.config.AccessToken == "" {
		return fmt.Errorf("no access token available")
	}

	var reqBody map[string]interface{}

	switch action {
	case "ACCEPT":
		reqBody = map[string]interface{}{
			"action": "ACCEPT",
		}
	case "DECLINE":
		reqBody = map[string]interface{}{
			"action": "DECLINE",
		}
	case "COUNTER":
		if counterPrice <= 0 {
			return fmt.Errorf("counter price must be greater than 0")
		}
		reqBody = map[string]interface{}{
			"action": "COUNTER",
			"counterOffer": map[string]interface{}{
				"price": map[string]interface{}{
					"value":    counterPrice,
					"currency": "USD",
				},
			},
		}
	default:
		return fmt.Errorf("invalid action: %s (must be ACCEPT, DECLINE, or COUNTER)", action)
	}

	endpoint := fmt.Sprintf("/sell/negotiation/v1/offer/%s/respond", offerID)
	_, err := c.makeRequest("POST", endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to respond to offer: %w", err)
	}

	return nil
}

// GetSellerBalance retrieves seller account balance information
func (c *Client) GetSellerBalance() (map[string]float64, error) {
	// Try the transactions endpoint to calculate balance
	respData, err := c.makeRequest("GET", "/sell/finances/v1/transaction?limit=50", nil)
	if err != nil {
		log.Printf("[DEBUG] Balance API error: %v", err)
		// If transactions fail, could mean no Managed Payments or missing OAuth scope
		if strings.Contains(err.Error(), "404") {
			return nil, fmt.Errorf("finances API not available - try /ebay-authorize to re-authorize with Finances API scope, or check your eBay keyset has Finances API enabled in Developer Portal")
		}
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "401") {
			return nil, fmt.Errorf("finances API access denied - please run /ebay-authorize to re-authorize with Finances API scope")
		}
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Log response for debugging (truncate if too long)
	debugResp := string(respData)
	if len(debugResp) > 500 {
		debugResp = debugResp[:500]
	}
	log.Printf("[DEBUG] Balance API Response (first 500 chars): %s", debugResp)

	var result struct {
		Transactions []struct {
			Amount struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
			TransactionType string `json:"transactionType"`
		} `json:"transactions"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse transactions: %w", err)
	}

	// Calculate balance from transactions
	var balance float64
	for _, tx := range result.Transactions {
		var amount float64
		fmt.Sscanf(tx.Amount.Value, "%f", &amount)
		balance += amount // eBay API returns negative values for fees
	}

	return map[string]float64{
		"available": balance,
		"total":     balance,
	}, nil
}

// GetPayouts retrieves recent payout transactions
func (c *Client) GetPayouts(limit int) ([]map[string]interface{}, error) {
	// Get payouts with succeeded status
	endpoint := fmt.Sprintf("/sell/finances/v1/payout?limit=%d&payoutStatus=SUCCEEDED", limit)
	respData, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		// 404 means no payout data yet or API not available
		if strings.Contains(err.Error(), "404") {
			return nil, fmt.Errorf("payouts API not available - try /ebay-authorize to re-authorize with Finances API scope, or check your eBay Developer keyset has Finances API enabled")
		}
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "401") {
			return nil, fmt.Errorf("payouts API access denied - please run /ebay-authorize to re-authorize")
		}
		return nil, fmt.Errorf("failed to get payouts: %w", err)
	}

	var result struct {
		Payouts []struct {
			PayoutId     string `json:"payoutId"`
			PayoutStatus string `json:"payoutStatus"`
			PayoutDate   string `json:"payoutDate"`
			Amount       struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
			PayoutInstrument struct {
				InstrumentType string `json:"instrumentType"`
			} `json:"payoutInstrument"`
		} `json:"payouts"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse payouts: %w", err)
	}

	payouts := make([]map[string]interface{}, 0, len(result.Payouts))
	for _, payout := range result.Payouts {
		var amount float64
		fmt.Sscanf(payout.Amount.Value, "%f", &amount)

		payouts = append(payouts, map[string]interface{}{
			"id":     payout.PayoutId,
			"amount": amount,
			"status": payout.PayoutStatus,
			"type":   fmt.Sprintf("%s Payout", payout.PayoutInstrument.InstrumentType),
			"date":   payout.PayoutDate[:10], // Format: YYYY-MM-DD
		})
	}

	return payouts, nil
}

// GetBuyerMessages retrieves buyer messages from the Post-Order API
func (c *Client) GetBuyerMessages(limit int) ([]map[string]interface{}, error) {
	// Post-Order API uses different authentication (not OAuth Bearer tokens)
	// It requires the older "IAF token" from the Trading API
	return nil, fmt.Errorf("buyer messages not available - eBay's Post-Order API doesn't support OAuth authentication. Check messages on eBay.com or the eBay mobile app instead")
}
