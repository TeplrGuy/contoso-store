# Inventory Price Feature Documentation

## Overview
This document describes the optional price field added to inventory items. This is a **non-breaking change** that allows inventory items to optionally include price information.

## Feature Summary
- **Status**: Implemented
- **Backwards Compatible**: Yes ✅
- **Required**: No (price is optional)

## API Changes

### InventoryItem Schema

```json
{
  "id": "string (required)",
  "item": "string (required)",
  "location": "string (required)",
  "priority": "string (required)",
  "price": {
    "value": "number (>= 0)",
    "currency": "string (3-letter uppercase, e.g., USD)"
  } // optional
}
```

### Example Payloads

**With Price:**
```json
{
  "id": "widget-001",
  "item": "Premium Widget",
  "location": "Seattle",
  "priority": "High",
  "price": {
    "value": 29.99,
    "currency": "USD"
  }
}
```

**Without Price (Backwards Compatible):**
```json
{
  "id": "widget-002",
  "item": "Basic Widget",
  "location": "Portland",
  "priority": "Standard"
}
```

## API Endpoints

### POST /inventory
Create or update an inventory item.

**Request Body:** InventoryItem (see schema above)

**Responses:**
- `201 Created` - Item created successfully
- `400 Bad Request` - Invalid JSON or missing required fields
- `422 Unprocessable Entity` - Validation error (e.g., invalid price)

### GET /inventory?id={id}
Retrieve an inventory item by ID.

**Parameters:**
- `id` (required) - Inventory item ID

**Responses:**
- `200 OK` - Returns InventoryItem
- `404 Not Found` - Item not found

## Price Validation Rules

When `price` is provided, it must satisfy:
1. **value**: Must be a number >= 0 (zero is allowed)
2. **currency**: Must be exactly 3 uppercase letters (e.g., USD, EUR, GBP, JPY)

### Valid Examples
- `{"value": 29.99, "currency": "USD"}`
- `{"value": 0, "currency": "EUR"}`
- `{"value": 1000.50, "currency": "JPY"}`

### Invalid Examples
- `{"value": -10, "currency": "USD"}` ❌ (negative value)
- `{"value": 29.99, "currency": "usd"}` ❌ (lowercase)
- `{"value": 29.99, "currency": "US"}` ❌ (too short)
- `{"value": 29.99, "currency": "USDD"}` ❌ (too long)

## Backwards Compatibility

### Existing Clients
Clients that don't send the `price` field will continue to work exactly as before:
```json
// This still works!
{
  "id": "old-item",
  "item": "Legacy Widget",
  "location": "Chicago",
  "priority": "Medium"
}
```

### Old Records
- Records created without `price` will not have the field when retrieved
- No migration is needed for existing data
- The `price` field will only appear in responses for items that have it

## Testing

### Unit Tests
Comprehensive unit tests are included in `go-service/app_test.go` and `go-service/integration_test.go`:
- ✅ Price validation (positive, negative, zero)
- ✅ Currency format validation (uppercase, length)
- ✅ JSON serialization/deserialization
- ✅ Backwards compatibility
- ✅ Optional field handling

Run tests:
```bash
cd go-service
go test -v ./...
```

### Manual Testing

#### Test with cURL (when services are running):

**Create item with price:**
```bash
curl -X POST http://localhost:3000/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-001",
    "item": "Premium Widget",
    "location": "Seattle",
    "priority": "High",
    "price": {
      "value": 29.99,
      "currency": "USD"
    }
  }'
```

**Create item without price:**
```bash
curl -X POST http://localhost:3000/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-002",
    "item": "Basic Widget",
    "location": "Portland",
    "priority": "Standard"
  }'
```

**Retrieve item:**
```bash
curl http://localhost:3000/inventory?id=test-001
```

#### Test with UI:
1. Navigate to http://localhost:3000
2. Scroll to "Create inventory item" section
3. Fill in required fields (ID, Item, Location, Priority)
4. Optionally fill in Price and Currency fields
5. Click "Create"
6. Success message will show the created item with price (if provided)

## Files Modified

1. **go-service/app.go**
   - Added `Price` and `InventoryItem` structs
   - Implemented POST/PUT handler with validation
   - Implemented GET handler with Dapr state store integration
   - Added `validateInventoryItem()` function

2. **go-service/app_test.go**
   - Added comprehensive unit tests for price validation

3. **go-service/integration_test.go**
   - Added JSON serialization/deserialization tests
   - Added backwards compatibility tests

4. **node-service/routes/inventory.js**
   - Updated to return JSON instead of plain text
   - Added POST endpoint for creating inventory items
   - Added error handling

5. **node-service/views/index.jade**
   - Added "Create inventory item" form
   - Added optional price input fields
   - Added client-side JavaScript for form submission
   - Added success/error message display

6. **node-service/open-api/swagger.json**
   - Added `Price` definition
   - Added `InventoryItem` definition
   - Updated `/inventory` GET and POST documentation

7. **dapr-components/statestore.yaml**
   - Added statestore component configuration for go-app

## Architecture Notes

### Data Flow
1. **Frontend** (node-service UI) → User enters inventory data with optional price
2. **Node Proxy** (node-service) → Forwards request to Go service via Dapr
3. **Go Service** (go-service) → Validates and stores in Dapr state store
4. **Dapr State Store** → Persists to Redis (locally) or Cosmos DB (in Azure)

### State Management
- Uses Dapr state management API
- State store name: `statestore`
- Each inventory item is stored with its ID as the key
- Price is stored as part of the JSON document (if present)

## Deployment Notes

### Local Development
1. Ensure Redis is running (required for Dapr state store)
2. Start services with Dapr:
   ```bash
   # Terminal 1: Go service
   dapr run --app-id go-app --app-port 8050 --dapr-http-port 3502 \
     --resources-path ./dapr-components -- go run .

   # Terminal 2: Node service
   dapr run --app-id node-app --app-port 3000 --dapr-http-port 3501 \
     --resources-path ./dapr-components -- npm run start
   ```

### Azure Deployment
- No changes required to deployment scripts
- The Cosmos DB state store in Azure will automatically handle the new schema
- Environment variables remain unchanged

## Future Enhancements

Potential improvements (not implemented in this version):
- [ ] Price history tracking
- [ ] Currency conversion
- [ ] Price range queries
- [ ] Bulk inventory operations
- [ ] Price alerts/notifications
- [ ] Support for multiple prices (e.g., wholesale, retail)

## Support

For questions or issues:
1. Check unit tests for expected behavior
2. Review Swagger documentation at `/swagger.json`
3. Verify Dapr state store configuration
4. Check logs for validation errors

## Change Log

### v0.1.0 - Initial Implementation
- Added optional price field to inventory items
- Implemented backend validation (Go service)
- Updated frontend UI (Node service)
- Added comprehensive tests
- Updated API documentation
- Maintained full backwards compatibility
