package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestInventoryItemJSONSerialization tests that InventoryItem can be properly
// serialized to and from JSON with optional price field
func TestInventoryItemJSONSerialization(t *testing.T) {
	testCases := []struct {
		name     string
		item     InventoryItem
		expected string
	}{
		{
			name: "Item with price",
			item: InventoryItem{
				ID:       "1",
				Item:     "Widget",
				Location: "Seattle",
				Priority: "Standard",
				Price: &Price{
					Value:    29.99,
					Currency: "USD",
				},
			},
			expected: `{"id":"1","item":"Widget","location":"Seattle","priority":"Standard","price":{"value":29.99,"currency":"USD"}}`,
		},
		{
			name: "Item without price",
			item: InventoryItem{
				ID:       "2",
				Item:     "Gadget",
				Location: "Portland",
				Priority: "High",
				Price:    nil,
			},
			expected: `{"id":"2","item":"Gadget","location":"Portland","priority":"High"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test serialization
			data, err := json.Marshal(tc.item)
			if err != nil {
				t.Fatalf("Failed to marshal item: %v", err)
			}

			// Compare JSON (ignore whitespace differences)
			var got, want interface{}
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("Failed to unmarshal result: %v", err)
			}
			if err := json.Unmarshal([]byte(tc.expected), &want); err != nil {
				t.Fatalf("Failed to unmarshal expected: %v", err)
			}

			gotJSON, _ := json.Marshal(got)
			wantJSON, _ := json.Marshal(want)
			if !bytes.Equal(gotJSON, wantJSON) {
				t.Errorf("JSON mismatch:\ngot:  %s\nwant: %s", gotJSON, wantJSON)
			}
		})
	}
}

// TestInventoryItemJSONDeserialization tests that InventoryItem can be properly
// deserialized from JSON
func TestInventoryItemJSONDeserialization(t *testing.T) {
	testCases := []struct {
		name     string
		jsonStr  string
		wantErr  bool
		validate func(*testing.T, InventoryItem)
	}{
		{
			name:    "Valid item with price",
			jsonStr: `{"id":"1","item":"Widget","location":"Seattle","priority":"Standard","price":{"value":29.99,"currency":"USD"}}`,
			wantErr: false,
			validate: func(t *testing.T, item InventoryItem) {
				if item.ID != "1" {
					t.Errorf("Expected ID=1, got %s", item.ID)
				}
				if item.Price == nil {
					t.Fatal("Expected price to be present")
				}
				if item.Price.Value != 29.99 {
					t.Errorf("Expected price value=29.99, got %f", item.Price.Value)
				}
				if item.Price.Currency != "USD" {
					t.Errorf("Expected currency=USD, got %s", item.Price.Currency)
				}
			},
		},
		{
			name:    "Valid item without price",
			jsonStr: `{"id":"2","item":"Gadget","location":"Portland","priority":"High"}`,
			wantErr: false,
			validate: func(t *testing.T, item InventoryItem) {
				if item.ID != "2" {
					t.Errorf("Expected ID=2, got %s", item.ID)
				}
				if item.Price != nil {
					t.Errorf("Expected price to be nil, got %+v", item.Price)
				}
			},
		},
		{
			name:    "Item with explicit null price",
			jsonStr: `{"id":"3","item":"Thing","location":"Denver","priority":"Low","price":null}`,
			wantErr: false,
			validate: func(t *testing.T, item InventoryItem) {
				if item.ID != "3" {
					t.Errorf("Expected ID=3, got %s", item.ID)
				}
				if item.Price != nil {
					t.Errorf("Expected price to be nil, got %+v", item.Price)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var item InventoryItem
			err := json.Unmarshal([]byte(tc.jsonStr), &item)

			if tc.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, item)
			}
		})
	}
}

// TestBackwardsCompatibility ensures that old clients without price field can still work
func TestBackwardsCompatibility(t *testing.T) {
	// Old format without price field
	oldFormatJSON := `{
		"id": "old-1",
		"item": "Legacy Widget",
		"location": "Chicago",
		"priority": "Medium"
	}`

	var item InventoryItem
	err := json.Unmarshal([]byte(oldFormatJSON), &item)
	if err != nil {
		t.Fatalf("Failed to unmarshal old format: %v", err)
	}

	// Validate that all fields are preserved
	if item.ID != "old-1" {
		t.Errorf("Expected ID=old-1, got %s", item.ID)
	}
	if item.Item != "Legacy Widget" {
		t.Errorf("Expected Item='Legacy Widget', got %s", item.Item)
	}
	if item.Location != "Chicago" {
		t.Errorf("Expected Location=Chicago, got %s", item.Location)
	}
	if item.Priority != "Medium" {
		t.Errorf("Expected Priority=Medium, got %s", item.Priority)
	}
	if item.Price != nil {
		t.Errorf("Expected Price to be nil, got %+v", item.Price)
	}

	// Ensure old format can be serialized back
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal back: %v", err)
	}

	// The serialized data should not contain price field
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if _, exists := result["price"]; exists {
		t.Error("Price field should not be present in serialized data when nil")
	}
}
