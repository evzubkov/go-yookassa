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
		Type            string `json:"type,omitempty"`
		ReturnUrl       string `json:"return_url,omitempty"`
		ConfirmationUrl string `json:"confirmation_url,omitempty"`
	}
	CardProduct struct {
		Code string `json:"code,omitempty"`
	}
	Card struct {
		First6      string      `json:"first6,omitempty"`
		Last4       string      `json:"last4,omitempty"`
		ExpiryYear  string      `json:"expiry_year,omitempty"`
		ExpiryMonth string      `json:"expiry_month,omitempty"`
		CardType    string      `json:"card_type,omitempty"`
		CardProduct interface{} `json:"card_product,omitempty"`
	}
	PaymentMethod struct {
		Id            string      `json:"id,omitempty"` //id метода оплаты для проведения автоплатежей
		Type          string      `json:"type,omitempty"`
		Title         string      `json:"title,omitempty"`
		Saved         bool        `json:"saved,omitempty"`
		Card          interface{} `json:"card,omitempty"`
		IssuerCountry string      `json:"issuer_country,omitempty"`
	}
	IncomeAmount struct {
		Value    string `json:"value,omitempty"`
		Currency string `json:"currency,omitempty"`
	}
)

type (
	NewPaymentRequest struct {
		Amount            Amount      `json:"amount,omitempty"`
		Capture           bool        `json:"capture,omitempty"` // true - мгновенное списание, false - двухстадийная оплата с холдированием средств
		Description       string      `json:"description,omitempty"`
		Confirmation      interface{} `json:"confirmation,omitempty"`
		PaymentMethod     interface{} `json:"payment_method_data,omitempty"`
		SavePaymentMethod string      `json:"save_payment_method,omitempty"` //true - для сохранения метода и оплаты и проведения автоплатежей
	}

	NewPaymentResponse struct {
		Id            string        `json:"id,omitempty"`
		Amount        Amount        `json:"amount,omitempty"`
		Status        string        `json:"status,omitempty"`
		PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
		Сonfirmation  Confirmation  `json:"confirmation,omitempty"`
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
	Id            string        `json:"id,omitempty"`
	Status        string        `json:"status,omitempty"`
	Amount        interface{}   `json:"amount,omitempty"`
	IncomeAmount  IncomeAmount  `json:"income_amount,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
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

type CapturePaymentResponse struct {
	Id            string        `json:"id,omitempty"`
	Amount        Amount        `json:"amount,omitempty"`
	Status        string        `json:"status,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
}

// CapturePayment  - move holding payment
//
// Transfer money with a two-stage character
func (o *Client) CapturePayment(ctx context.Context, amount Amount, paymentId string) (result CapturePaymentResponse, err error) {

	payload, err := json.Marshal(amount)
	if err != nil {
		return
	}

	req, _ := http.NewRequestWithContext(
		ctx, "POST", fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s/capture", paymentId), bytes.NewBuffer(payload))
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

type CancelPaymentResponse struct {
	Id            string        `json:"id,omitempty"`
	Amount        Amount        `json:"amount,omitempty"`
	Status        string        `json:"status,omitempty"`
	Paid          bool          `json:"paid,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
}

// CancelPayment  - cancel holding payment
//
// Cancel payment for two-step payment
func (o *Client) CancelPayment(ctx context.Context, paymentId string) (result CancelPaymentResponse, err error) {

	req, _ := http.NewRequestWithContext(
		ctx, "POST", fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s/cancel", paymentId), nil)
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
