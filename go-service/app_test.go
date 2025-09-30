package main

import (
	"testing"
)

func TestValidateInventoryItem_ValidPrice(t *testing.T) {
	item := &InventoryItem{
		ID:       "1",
		Item:     "Widget",
		Location: "Seattle",
		Priority: "Standard",
		Price: &Price{
			Value:    29.99,
			Currency: "USD",
		},
	}

	err := validateInventoryItem(item)
	if err != nil {
		t.Errorf("Expected no error for valid price, got: %v", err)
	}
}

func TestValidateInventoryItem_NoPriceIsValid(t *testing.T) {
	item := &InventoryItem{
		ID:       "1",
		Item:     "Widget",
		Location: "Seattle",
		Priority: "Standard",
		Price:    nil,
	}

	err := validateInventoryItem(item)
	if err != nil {
		t.Errorf("Expected no error when price is nil, got: %v", err)
	}
}

func TestValidateInventoryItem_NegativePrice(t *testing.T) {
	item := &InventoryItem{
		ID:       "1",
		Item:     "Widget",
		Location: "Seattle",
		Priority: "Standard",
		Price: &Price{
			Value:    -10.0,
			Currency: "USD",
		},
	}

	err := validateInventoryItem(item)
	if err == nil {
		t.Error("Expected error for negative price, got nil")
	}
}

func TestValidateInventoryItem_InvalidCurrency(t *testing.T) {
	testCases := []struct {
		currency string
		name     string
	}{
		{"US", "too short"},
		{"USDD", "too long"},
		{"us", "lowercase"},
		{"Us", "mixed case"},
		{"U$D", "special char"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item := &InventoryItem{
				ID:       "1",
				Item:     "Widget",
				Location: "Seattle",
				Priority: "Standard",
				Price: &Price{
					Value:    29.99,
					Currency: tc.currency,
				},
			}

			err := validateInventoryItem(item)
			if err == nil {
				t.Errorf("Expected error for invalid currency '%s', got nil", tc.currency)
			}
		})
	}
}

func TestValidateInventoryItem_ValidCurrencies(t *testing.T) {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD"}

	for _, currency := range currencies {
		t.Run(currency, func(t *testing.T) {
			item := &InventoryItem{
				ID:       "1",
				Item:     "Widget",
				Location: "Seattle",
				Priority: "Standard",
				Price: &Price{
					Value:    29.99,
					Currency: currency,
				},
			}

			err := validateInventoryItem(item)
			if err != nil {
				t.Errorf("Expected no error for valid currency '%s', got: %v", currency, err)
			}
		})
	}
}

func TestValidateInventoryItem_ZeroPrice(t *testing.T) {
	item := &InventoryItem{
		ID:       "1",
		Item:     "Widget",
		Location: "Seattle",
		Priority: "Standard",
		Price: &Price{
			Value:    0.0,
			Currency: "USD",
		},
	}

	err := validateInventoryItem(item)
	if err != nil {
		t.Errorf("Expected no error for zero price, got: %v", err)
	}
}
