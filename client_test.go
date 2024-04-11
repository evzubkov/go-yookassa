package yookassa_test

import (
	"context"
	"os"
	"testing"

	"github.com/evzubkov/go-yookassa"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestСClient(t *testing.T) {
	godotenv.Load()
	client := yookassa.NewClient(os.Getenv("SHOP_ID"), os.Getenv("API_KEY"))

	paymentObject, err := client.NewPayment(context.TODO(), yookassa.NewPaymentRequest{
		Description: "Услуга",
		Capture:     true,
		Amount: yookassa.Amount{
			Value:    "100.00",
			Currency: "RUB",
		},
		Confirmation: yookassa.Confirmation{
			Type:      "redirect",
			ReturnUrl: "https://example.com",
		},
	})
	assert.NoError(t, err)

	_, err = client.CheckPaymentStatus(context.TODO(), paymentObject.Id)
	assert.NoError(t, err)
}
