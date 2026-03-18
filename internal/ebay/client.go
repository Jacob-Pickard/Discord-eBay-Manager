package ebay

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

// SaveTokensToEnv persists the current access and refresh tokens to the .env file
func (c *Client) SaveTokensToEnv() error {
	envPath := ".env"
	data, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "EBAY_ACCESS_TOKEN=") {
			lines[i] = "EBAY_ACCESS_TOKEN=" + c.config.AccessToken
		} else if strings.HasPrefix(line, "EBAY_REFRESH_TOKEN=") && c.config.RefreshToken != "" {
			lines[i] = "EBAY_REFRESH_TOKEN=" + c.config.RefreshToken
		}
	}

	if err := os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	log.Println("✅ Tokens persisted to .env file")
	return nil
}

// CheckConnection verifies the eBay API connection
func (c *Client) CheckConnection() string {
	if c.config.AccessToken == "" {
		return "âŒ **Not authorized**\n\nNo access token found. Run `/ebay-authorize` to connect your eBay account."
	}

	// Test with a lightweight API call
	_, err := c.makeRequest("GET", "/sell/account/v1/privilege", nil)
	if err != nil {
		return fmt.Sprintf("âš ï¸ **Token configured but API test failed**\n\nEnvironment: %s\nError: %v\n\nâš¡ Try re-authorizing with `/ebay-authorize`",
			c.config.Environment, err)
	}

	return fmt.Sprintf("âœ… **Connected to eBay API**\n\nEnvironment: `%s`\nToken: `%s...` âœ…\n\n💡 Use `/ebay-scopes` to see your token's permissions.",
		c.config.Environment,
		c.config.AccessToken[:10])
}

// GetTokenScopes returns the OAuth scopes configured for this application
func (c *Client) GetTokenScopes() map[string]interface{} {
	// Return scope info with token status
	result := make(map[string]interface{})
	result["hasToken"] = c.config.AccessToken != ""
	result["environment"] = c.config.Environment
	result["requestedScopes"] = []string{
		"https://api.ebay.com/oauth/api_scope",
		"https://api.ebay.com/oauth/api_scope/sell.inventory",
		"https://api.ebay.com/oauth/api_scope/sell.fulfillment",
		"https://api.ebay.com/oauth/api_scope/sell.account",
		"https://api.ebay.com/oauth/api_scope/sell.finances",
		"https://api.ebay.com/oauth/api_scope/commerce.identity.readonly",
		"https://api.ebay.com/oauth/api_scope/commerce.notification.subscription",
	}
	return result
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
	// Finances API uses apiz.ebay.com instead of api.ebay.com
	if strings.Contains(endpoint, "/sell/finances/") {
		fullURL = strings.Replace(fullURL, "api.ebay.com", "apiz.ebay.com", 1)
	}
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

	// Log all outbound requests
	log.Printf("[API] %s %s", method, fullURL)

	// Browse and Commerce APIs require a marketplace ID
	if strings.Contains(endpoint, "/buy/") || strings.Contains(endpoint, "/commerce/") {
		req.Header.Set("X-EBAY-C-MARKETPLACE-ID", "EBAY_US")
	}

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

	// Enhanced logging for all API errors
	if resp.StatusCode >= 400 {
		log.Printf("[API ERROR] %s %s => HTTP %d: %s", method, fullURL, resp.StatusCode, string(respBody))
		if strings.Contains(endpoint, "/finances/") {
			log.Printf("[DEBUG] Finances Response Headers: %v", resp.Header)
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[API] %s %s => HTTP %d (%d bytes)", method, fullURL, resp.StatusCode, len(respBody))

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

	// Try to fetch offers from Negotiation API
	// Note: eBay's Negotiation API may require offer IDs from notifications
	// If this fails, advise users to use webhook notifications
	endpoint := "/sell/negotiation/v1/offer"

	respData, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		// Return error with helpful message
		return nil, fmt.Errorf("eBay API doesn't support listing all offers directly. Use webhook notifications to receive offer alerts in real-time")
	}

	var response OffersResponse
	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse offers response: %w", err)
	}

	// Filter for PENDING offers only
	pendingOffers := make([]Offer, 0)
	for _, offer := range response.Offers {
		if offer.Status == "PENDING" {
			pendingOffers = append(pendingOffers, offer)
		}
	}

	return pendingOffers, nil
}

// GetSellerUsername retrieves the authenticated seller's eBay username via the Identity API.
// Requires commerce.identity.readonly scope (granted after re-authorizing with /ebay-authorize).
func (c *Client) GetSellerUsername() (string, error) {
	respData, err := c.makeRequest("GET", "/commerce/identity/v1/user/", nil)
	if err != nil {
		return "", fmt.Errorf("identity API error: %w", err)
	}
	var result struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(respData, &result); err != nil {
		return "", fmt.Errorf("failed to parse identity response: %w", err)
	}
	if result.Username == "" {
		return "", fmt.Errorf("no username returned by Identity API")
	}
	return result.Username, nil
}

// GetListings retrieves active listings via the Trading API GetMyeBaySelling.
// Uses the seller's OAuth access token — works for all traditionally-listed items.
func (c *Client) GetListings(limit int) ([]Listing, error) {
	if c.config.AccessToken == "" {
		return nil, fmt.Errorf("no access token - run /ebay-authorize first")
	}
	if limit <= 0 || limit > 200 {
		limit = 10
	}

	// Trading API — GetMyeBaySelling returns all active listings for the authenticated seller
	tradingURL := "https://api.ebay.com/ws/api.dll"
	if c.config.Environment != "PRODUCTION" {
		tradingURL = "https://api.sandbox.ebay.com/ws/api.dll"
	}

	reqBody := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<GetMyeBaySellingRequest xmlns="urn:ebay:apis:eBLBaseComponents">
  <ErrorLanguage>en_US</ErrorLanguage>
  <WarningLevel>High</WarningLevel>
  <ActiveList>
    <Include>true</Include>
    <Pagination>
      <EntriesPerPage>%d</EntriesPerPage>
      <PageNumber>1</PageNumber>
    </Pagination>
  </ActiveList>
  <DetailLevel>ReturnAll</DetailLevel>
</GetMyeBaySellingRequest>`, limit)

	req, err := http.NewRequest("POST", tradingURL, strings.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to build Trading API request: %w", err)
	}
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("X-EBAY-API-SITEID", "0")
	req.Header.Set("X-EBAY-API-COMPATIBILITY-LEVEL", "967")
	req.Header.Set("X-EBAY-API-CALL-NAME", "GetMyeBaySelling")
	req.Header.Set("X-EBAY-API-IAF-TOKEN", c.config.AccessToken)

	log.Printf("[API] POST %s (Trading API: GetMyeBaySelling)", tradingURL)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Trading API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Trading API response: %w", err)
	}
	log.Printf("[API] POST Trading API => HTTP %d (%d bytes)", resp.StatusCode, len(body))
	log.Printf("[DEBUG] Trading API response: %s", string(body))

	var result struct {
		XMLName xml.Name `xml:"GetMyeBaySellingResponse"`
		Ack     string   `xml:"Ack"`
		Errors  []struct {
			LongMessage string `xml:"LongMessage"`
		} `xml:"Errors"`
		ActiveList struct {
			ItemArray struct {
				Items []struct {
					ItemID        string `xml:"ItemID"`
					Title         string `xml:"Title"`
					Quantity      int    `xml:"Quantity"`
					SellingStatus struct {
						CurrentPrice struct {
							CurrencyID string  `xml:"currencyID,attr"`
							Value      float64 `xml:",chardata"`
						} `xml:"CurrentPrice"`
						QuantityRemaining int `xml:"QuantityRemaining"`
					} `xml:"SellingStatus"`
					ShippingDetails struct {
						ShippingServiceOptions []struct {
							ShippingServiceCost struct {
								Value float64 `xml:",chardata"`
							} `xml:"ShippingServiceCost"`
						} `xml:"ShippingServiceOptions"`
						ShippingType string `xml:"ShippingType"`
					} `xml:"ShippingDetails"`
					ConditionDisplayName string `xml:"ConditionDisplayName"`
					PictureDetails       struct {
						GalleryURL string   `xml:"GalleryURL"`
						PictureURL []string `xml:"PictureURL"`
					} `xml:"PictureDetails"`
					ListingDetails struct {
						ViewItemURL string `xml:"ViewItemURL"`
					} `xml:"ListingDetails"`
				} `xml:"Item"`
			} `xml:"ItemArray"`
		} `xml:"ActiveList"`
	}

	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Trading API response: %w", err)
	}
	if result.Ack != "Success" && result.Ack != "Warning" {
		msg := "unknown error"
		if len(result.Errors) > 0 {
			msg = result.Errors[0].LongMessage
		}
		return nil, fmt.Errorf("Trading API error: %s", msg)
	}

	items := result.ActiveList.ItemArray.Items
	listings := make([]Listing, 0, len(items))
	for _, item := range items {
		price := item.SellingStatus.CurrentPrice.Value
		currency := item.SellingStatus.CurrentPrice.CurrencyID
		if currency == "" {
			currency = "USD"
		}
		qty := item.SellingStatus.QuantityRemaining
		if qty == 0 {
			qty = item.Quantity
		}

		shipping := "See listing"
		if item.ShippingDetails.ShippingType == "Free" {
			shipping = "Free"
		} else if len(item.ShippingDetails.ShippingServiceOptions) > 0 {
			cost := item.ShippingDetails.ShippingServiceOptions[0].ShippingServiceCost.Value
			if cost == 0 {
				shipping = "Free"
			} else {
				shipping = fmt.Sprintf("$%.2f", cost)
			}
		}

		// GetMyeBaySelling only returns GalleryURL (s-l140.jpg, 140px).
		// eBay's CDN supports larger sizes via URL suffix substitution.
		imageURL := item.PictureDetails.GalleryURL
		if len(item.PictureDetails.PictureURL) > 0 {
			imageURL = item.PictureDetails.PictureURL[0]
		}
		// Upgrade thumbnail to 500px version by replacing size suffix
		imageURL = strings.Replace(imageURL, "s-l140.jpg", "s-l500.jpg", 1)
		imageURL = strings.Replace(imageURL, "s-l96.jpg", "s-l500.jpg", 1)

		listings = append(listings, Listing{
			Title:      item.Title,
			Price:      price,
			Currency:   currency,
			Shipping:   shipping,
			Quantity:   qty,
			Condition:  item.ConditionDisplayName,
			ImageURL:   imageURL,
			ListingURL: item.ListingDetails.ViewItemURL,
		})
	}

	return listings, nil
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
	// Use seller_funds_summary endpoint to get pending payout amount
	respData, err := c.makeRequest("GET", "/sell/finances/v1/seller_funds_summary", nil)
	if err != nil {
		log.Printf("[DEBUG] Balance API error: %v", err)
		if strings.Contains(err.Error(), "404") {
			return nil, fmt.Errorf("finances API not available - ensure your eBay account is enrolled in Managed Payments")
		}
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "401") {
			return nil, fmt.Errorf("finances API access denied - run /ebay-authorize to re-authorize with Finances API scope")
		}
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	log.Printf("[DEBUG] Balance API Response: %s", string(respData))

	var result struct {
		AvailableFunds struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"availableFunds"`
		TotalBalance struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"totalBalance"`
	}

	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse balance response: %w", err)
	}

	// Parse the available funds (pending payout amount)
	var available float64
	fmt.Sscanf(result.AvailableFunds.Value, "%f", &available)

	var total float64
	fmt.Sscanf(result.TotalBalance.Value, "%f", &total)

	return map[string]float64{
		"available": available,
		"total":     total,
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
