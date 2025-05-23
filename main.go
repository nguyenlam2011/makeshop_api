package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type PaypayAttributes struct {
	PaymentID string `json:"paymentId"`
}

type Order struct {
	ShopID            string           `json:"shopId"`
	SystemOrderNumber string           `json:"systemOrderNumber"`
	Amount            int              `json:"amount"`
	PayMethod         string           `json:"payMethod"`
	Status            string           `json:"status"`
	OrderedAt         string           `json:"orderedAt"`
	PaypayAttributes  PaypayAttributes `json:"paypayAttributes"`
}

type OrdersResponse struct {
	Orders map[string][]Order `json:"orders"`
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	payMethod := r.URL.Query().Get("payMethod")
	shopIDs := r.URL.Query()["shopID"] // shopID as a list

	log.Printf("Received /v1/orders call: month=%s, payMethod=%s, shopIDs=%v", month, payMethod, shopIDs)

	ordersMap := make(map[string][]Order)
	for i, shopID := range shopIDs {
		order := Order{
			ShopID:            shopID,
			SystemOrderNumber: "SO123456-" + shopID,
			Amount:            1000 + i*100,
			PayMethod:         payMethod,
			Status:            "paid",
			OrderedAt:         time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			PaypayAttributes:  PaypayAttributes{PaymentID: "PP1234567890-" + shopID},
		}
		ordersMap[shopID] = append(ordersMap[shopID], order)
	}
	resp := OrdersResponse{
		Orders: ordersMap,
	}

	// Log the output response as JSON
	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
	} else {
		log.Printf("Response: %s", respJSON)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/v1/orders", ordersHandler)
	http.ListenAndServe(":3000", nil)
}
