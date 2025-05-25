package tool

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func RegisterSupportTools(registry *Registry) {
	registry.Register(Tool{
		Name:        "customer_info",
		Description: "Get information about a customer order or account",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query_type": map[string]interface{}{
					"type":        "string",
					"description": "Type of information to look up (order, shipping, return)",
					"enum":        []string{"order", "shipping", "return"},
				},
				"order_id": map[string]interface{}{
					"type":        "string",
					"description": "ID of the order to check",
				},
			},
			"required": []string{"query_type", "order_id"},
		},
		Handler: handleCustomerInfo,
	})
}

func handleCustomerInfo(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	queryType, _ := args["query_type"].(string)
	orderID, _ := args["order_id"].(string)

	if queryType == "" || orderID == "" {
		return nil, fmt.Errorf("query_type and order_id are required")
	}

	order := map[string]interface{}{
		"order_id":       orderID,
		"date":           time.Now().AddDate(0, 0, -rand.Intn(30)).Format("2006-01-02"),
		"status":         randomChoice([]string{"Delivered", "Shipped", "Processing", "Cancelled"}),
		"items":          rand.Intn(5) + 1,
		"total_amount":   fmt.Sprintf("$%.2f", 10.0+rand.Float64()*200.0),
		"customer_email": "customer@example.com",
	}

	switch queryType {
	case "order":

	case "shipping":
		status := randomChoice([]string{"Order Placed", "Processing", "Shipped", "Out for Delivery", "Delivered"})
		order["shipping_info"] = map[string]interface{}{
			"status":             status,
			"tracking_number":    fmt.Sprintf("TRK%d", rand.Intn(1000000)),
			"carrier":            randomChoice([]string{"FedEx", "UPS", "USPS", "DHL"}),
			"estimated_delivery": time.Now().AddDate(0, 0, rand.Intn(5)).Format("2006-01-02"),
		}

	case "return":
		orderDate := time.Now().AddDate(0, 0, -rand.Intn(60))
		daysSinceOrder := int(time.Since(orderDate).Hours() / 24)
		isEligible := daysSinceOrder <= 30

		order["return_info"] = map[string]interface{}{
			"is_eligible":           isEligible,
			"days_since_purchase":   daysSinceOrder,
			"return_window_days":    30,
			"days_left":             max(0, 30-daysSinceOrder),
			"return_policy_details": "Items must be in original packaging and undamaged.",
		}
	}

	return order, nil
}

func randomChoice(options []string) string {
	return options[rand.Intn(len(options))]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
