# go-yookassa

Неофициальный клиент для создания платежей в сервисе ЮКасса.

# Использование
```
    client := yookassa.NewClient(os.Getenv("SHOP_ID"), os.Getenv("API_KEY"))

    // Создание платежа
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

    // Получение статуса платежа
    paymentStatus, err = client.CheckPaymentStatus(context.TODO(), paymentObject.Id)

```

Возможные варианты статуса платежа:

- ```pending``` - ожидается оплата

- ```succeeded``` - успешно

- ```canceled``` - что-то пошло не так