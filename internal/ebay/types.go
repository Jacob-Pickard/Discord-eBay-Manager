package ebay

import "time"

// Order represents an eBay order
type Order struct {
	OrderID                      string                   `json:"orderId"`
	CreationDate                 time.Time                `json:"creationDate"`
	Buyer                        Buyer                    `json:"buyer"`
	BuyerUsername                string                   `json:"buyerUsername"` // Computed field
	PricingSummary               PricingSummary           `json:"pricingSummary"`
	TotalPrice                   float64                  `json:"totalPrice"` // Computed field
	Currency                     string                   `json:"currency"`   // Computed field
	FulfillmentStartInstructions []FulfillmentInstruction `json:"fulfillmentStartInstructions"`
	OrderFulfillmentStatus       string                   `json:"orderFulfillmentStatus"`
	FulfillmentStatus            string                   `json:"fulfillmentStatus"` // Computed field
	LineItems                    []LineItem               `json:"lineItems"`
}

// Buyer represents the buyer information
type Buyer struct {
	Username string `json:"username"`
}

// PricingSummary contains order pricing details
type PricingSummary struct {
	Total struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"total"`
}

// FulfillmentInstruction contains shipping details
type FulfillmentInstruction struct {
	ShippingStep struct {
		ShipTo Address `json:"shipTo"`
	} `json:"shippingStep"`
}

// LineItem represents a single item in an order
type LineItem struct {
	LineItemID   string `json:"lineItemId"`
	Title        string `json:"title"`
	Quantity     int    `json:"quantity"`
	LineItemCost struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"lineItemCost"`
	Price float64 `json:"price"` // Computed field
	SKU   string  `json:"sku"`
	Image struct {
		ImageUrl string `json:"imageUrl"`
	} `json:"image"`
	ImageUrl     string `json:"imageUrl"` // Computed field
	LegacyItemId string `json:"legacyItemId"`
}

// Address represents a shipping address
type Address struct {
	Name       string `json:"name"`
	Street1    string `json:"street1"`
	Street2    string `json:"street2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// Offer represents a buyer offer
type Offer struct {
	OfferID       string    `json:"offerId"`
	ItemID        string    `json:"itemId"`
	ItemTitle     string    `json:"itemTitle"`
	BuyerUsername string    `json:"buyerUsername"`
	OfferPrice    float64   `json:"offerPrice"`
	ListPrice     float64   `json:"listPrice"`
	Currency      string    `json:"currency"`
	CreatedDate   time.Time `json:"createdDate"`
	Status        string    `json:"status"`
}

// ShippingLabel represents a shipping label
type ShippingLabel struct {
	LabelID        string    `json:"labelId"`
	OrderID        string    `json:"orderId"`
	TrackingNumber string    `json:"trackingNumber"`
	LabelURL       string    `json:"labelUrl"`
	Cost           float64   `json:"cost"`
	CreatedDate    time.Time `json:"createdDate"`
}
