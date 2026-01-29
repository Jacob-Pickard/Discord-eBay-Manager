package ebay

import "time"

// Order represents an eBay order
type Order struct {
	OrderID          string    `json:"orderId"`
	CreationDate     time.Time `json:"creationDate"`
	BuyerUsername    string    `json:"buyerUsername"`
	TotalPrice       float64   `json:"totalPrice"`
	Currency         string    `json:"currency"`
	ShippingAddress  Address   `json:"shippingAddress"`
	LineItems        []LineItem `json:"lineItems"`
	FulfillmentStatus string   `json:"fulfillmentStatus"`
}

// LineItem represents a single item in an order
type LineItem struct {
	LineItemID     string  `json:"lineItemId"`
	Title          string  `json:"title"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
	SKU            string  `json:"sku"`
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
	OfferID      string    `json:"offerId"`
	ItemID       string    `json:"itemId"`
	ItemTitle    string    `json:"itemTitle"`
	BuyerUsername string   `json:"buyerUsername"`
	OfferPrice   float64   `json:"offerPrice"`
	ListPrice    float64   `json:"listPrice"`
	Currency     string    `json:"currency"`
	CreatedDate  time.Time `json:"createdDate"`
	Status       string    `json:"status"`
}

// Listing represents an eBay listing
type Listing struct {
	SKU          string   `json:"sku"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	Quantity     int      `json:"quantity"`
	Category     string   `json:"category"`
	Condition    string   `json:"condition"`
	ImageURLs    []string `json:"imageUrls"`
	ShippingCost float64  `json:"shippingCost"`
}

// ShippingLabel represents a shipping label
type ShippingLabel struct {
	LabelID      string    `json:"labelId"`
	OrderID      string    `json:"orderId"`
	TrackingNumber string  `json:"trackingNumber"`
	LabelURL     string    `json:"labelUrl"`
	Cost         float64   `json:"cost"`
	CreatedDate  time.Time `json:"createdDate"`
}
