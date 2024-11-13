package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	portNum        string = ":8080"
	originDomain   string = "https://www.example.com"
	successURL     string = "https://www.example.com/success"
	failureURL     string = "https://www.example.com/failure"
	cancelURL      string = "https://www.example.com/cancel"
	urlProd        string = "https://restapi.payplus.co.il/api/v1.0/PaymentPages/generateLink"
	urlDev         string = "https://restapidev.payplus.co.il/api/v1.0/PaymentPages/generateLink"
	paymentPageUID string = "" // Payment page UID
	apiKey         string = "" // API key
	secretKey      string = "" // Secret key
)

func main() {
	http.HandleFunc("/get_payment_page", getPaymentPage)

	log.Println("Started on port", portNum)

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getPaymentPage(w http.ResponseWriter, r *http.Request) {
	url := urlDev
	// url := urlProd for production

	var rObject = struct {
		Amount         int    `json:"amount"`
		CurrencyCode   string `json:"currency_code"`
		PaymentPageUID string `json:"payment_page_uid"`
		RefURLSuccess  string `json:"refURL_success"`
		RefURLOrigin   string `json:"refURL_origin"`
		RefURLFailure  string `json:"refURL_failure"`
		RefURLCancel   string `json:"refURL_cancel"`
		CreateToken    bool   `json:"create_token"`
		HostedFields   bool   `json:"hosted_fields"`
		ChargeMethod   int    `json:"charge_method"`
		Customer       struct {
			CustomerName string `json:"customer_name"`
			Email        string `json:"email"`
			Phone        string `json:"phone"`
		} `json:"customer"`
		Items []struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
			Price    int    `json:"price"`
			VatType  int    `json:"vat_type"`
		} `json:"items"`
	}{
		Amount:         100,
		CurrencyCode:   "ILS",
		PaymentPageUID: paymentPageUID,
		RefURLSuccess:  successURL,
		RefURLOrigin:   originDomain,
		RefURLFailure:  failureURL,
		RefURLCancel:   cancelURL,
		CreateToken:    true,
		HostedFields:   true,
		ChargeMethod:   1,
		Customer: struct {
			CustomerName string `json:"customer_name"`
			Email        string `json:"email"`
			Phone        string `json:"phone"`
		}{
			CustomerName: "John Doe",
			Email:        "test@@example.com",
			Phone:        "1234567890",
		},
		Items: []struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
			Price    int    `json:"price"`
			VatType  int    `json:"vat_type"`
		}{
			{
				Name:     "item1",
				Quantity: 1,
				Price:    100,
				VatType:  1,
			},
		},
	}

	jsonStr, err := json.Marshal(rObject)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	auth := struct {
		APIKey    string `json:"api_key"`
		SecretKey string `json:"secret_key"`
	}{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}
	authStr, err := json.Marshal(auth)
	req.Header.Set("Authorization", string(authStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(body))
}
