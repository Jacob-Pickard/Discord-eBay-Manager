package ebay

import (
	"testing"
	"time"
)

func TestOrderValidation(t *testing.T) {
	order := Order{
		OrderID:       "12345",
		BuyerUsername: "testbuyer",
		LineItems: []LineItem{
			{
				LineItemID: "item-1",
				Title:      "Test Item",
				Quantity:   1,
				Price:      50.00,
				SKU:        "TEST-SKU",
			},
		},
		TotalPrice:        50.00,
		Currency:          "USD",
		FulfillmentStatus: "FULFILLED",
		CreationDate:      time.Now(),
	}

	if order.OrderID == "" {
		t.Error("OrderID should not be empty")
	}

	if len(order.LineItems) == 0 {
		t.Error("Order should have at least one line item")
	}

	if order.LineItems[0].Price != 50.00 {
		t.Errorf("Expected price 50.00, got %f", order.LineItems[0].Price)
	}

	if order.Currency != "USD" {
		t.Errorf("Expected currency USD, got %s", order.Currency)
	}
}

func TestOfferValidation(t *testing.T) {
	offer := Offer{
		OfferID:       "offer-123",
		ItemID:        "item-456",
		ItemTitle:     "Test Item",
		BuyerUsername: "offerbuyer",
		OfferPrice:    45.00,
		ListPrice:     60.00,
		Currency:      "USD",
		Status:        "PENDING",
		CreatedDate:   time.Now(),
	}

	if offer.OfferID == "" {
		t.Error("OfferID should not be empty")
	}

	if offer.OfferPrice <= 0 {
		t.Error("Offered amount should be greater than 0")
	}

	if offer.Status != "PENDING" {
		t.Errorf("Expected status PENDING, got %s", offer.Status)
	}

	if offer.Currency != "USD" {
		t.Errorf("Expected currency USD, got %s", offer.Currency)
	}
}

func TestShippingLabelValidation(t *testing.T) {
	label := ShippingLabel{
		LabelID:        "label-123",
		OrderID:        "order-456",
		TrackingNumber: "1Z999AA10123456789",
		LabelURL:       "https://example.com/label.pdf",
		Cost:           8.50,
		CreatedDate:    time.Now(),
	}

	if label.LabelID == "" {
		t.Error("LabelID should not be empty")
	}

	if label.TrackingNumber == "" {
		t.Error("TrackingNumber should not be empty")
	}

	if label.Cost <= 0 {
		t.Error("Cost should be greater than 0")
	}

	if label.LabelURL == "" {
		t.Error("LabelURL should not be empty")
	}
}

func TestAddressValidation(t *testing.T) {
	address := Address{
		Name:       "John Doe",
		Street1:    "123 Main St",
		City:       "New York",
		State:      "NY",
		PostalCode: "10001",
		Country:    "US",
	}

	if address.Name == "" {
		t.Error("Name should not be empty")
	}

	if address.Street1 == "" {
		t.Error("Street1 should not be empty")
	}

	if address.City == "" {
		t.Error("City should not be empty")
	}

	if address.PostalCode == "" {
		t.Error("PostalCode should not be empty")
	}
}
