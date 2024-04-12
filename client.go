package yookassa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type Client struct {
	shopId string
	apiKey string
}

func NewClient(shopId, apiKey string) *Client {
	return &Client{shopId: shopId, apiKey: apiKey}
}

type (
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	}
	Confirmation struct {
		Type      string `json:"type"`
		ReturnUrl string `json:"return_url"`
	}
	NewPaymentRequest struct {
		Amount       Amount       `json:"amount"`
		Capture      bool         `json:"capture"`
		Confirmation Confirmation `json:"confirmation"`
		Description  string       `json:"description"`
	}

	NewPaymentResponse struct {
		Id           string `json:"id"`
		Status       string `json:"status"`
		Сonfirmation struct {
			Type            string `json:"type"`
			ConfirmationUrl string `json:"confirmation_url"`
		} `json:"confirmation"`
	}
)

// NewPayment - create new payment
func (o *Client) NewPayment(ctx context.Context, payment NewPaymentRequest) (
	result NewPaymentResponse, err error) {

	payload, err := json.Marshal(payment)
	if err != nil {
		return
	}

	req, _ := http.NewRequestWithContext(
		ctx, "POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(payload))
	req.SetBasicAuth(o.shopId, o.apiKey)
	req.Header.Add("Idempotence-Key", uuid.New().String())
	req.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("fail to send response to server. Status code: %d, info: %+v", resp.StatusCode, cast.ToString(body))
		return
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return
	}

	return
}

type CheckPaymentStatusResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	IncomeAmount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"income_amount"`
	PaymentMethod struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Title string `json:"title"`
		Card  struct {
			First6      string `json:"first6"`
			Last4       string `json:"last4"`
			ExpiryYear  string `json:"expiry_year"`
			ExpiryMonth string `json:"expiry_month"`
			CardType    string `json:"card_type"`
			CardProduct struct {
				Code string `json:"code"`
			} `json:"card_product"`
		} `json:"card"`
		IssuerCountry string `json:"issuer_country"`
	} `json:"payment_method"`
}

// CheckPaymentStatus - check payment status
func (o *Client) CheckPaymentStatus(ctx context.Context, paymentId string) (result CheckPaymentStatusResponse, err error) {

	req, _ := http.NewRequestWithContext(
		ctx, "GET", fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s", paymentId), nil)
	req.SetBasicAuth(o.shopId, o.apiKey)
	req.Header.Add("Idempotence-Key", uuid.New().String())
	req.Header.Add("Content-Type", "application/json")

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("fail to send response to server. Status code: %d, info: %+v", resp.StatusCode, cast.ToString(body))
		return
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return
	}

	return
}
