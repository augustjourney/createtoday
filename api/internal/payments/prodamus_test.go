package payments

import (
	"context"
	"createtodayapi/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func newProdamusSystem() (*Prodamus, *config.Config) {
	prodamus := Prodamus{}
	conf := config.New("../../.env")
	return &prodamus, conf
}

func TestProdamusGenerateQueryParams(t *testing.T) {
	t.Parallel()
	prodamus, _ := newProdamusSystem()
	payload := ProdamusInitPayload{
		Amount:      100,
		Description: "Тестовый продукт",
		OrderId:     1337,
		Email:       "test@test.com",
		Phone:       "18123801133",
	}

	params := prodamus.generateQueryParams(payload)
	vals, err := url.ParseQuery(params)

	require.NoError(t, err)

	do, ok := vals["do"]
	require.True(t, ok)
	assert.Equal(t, "link", do[0])

	orderId, ok := vals["order_id"]
	require.True(t, ok)
	assert.Equal(t, "1337", orderId[0])

	customerPhone, ok := vals["customer_phone"]
	require.True(t, ok)
	assert.Equal(t, "18123801133", customerPhone[0])

	customerEmail, ok := vals["customer_email"]
	require.True(t, ok)
	assert.Equal(t, "test@test.com", customerEmail[0])

	productPrice, ok := vals["products[0][price]"]
	require.True(t, ok)
	assert.Equal(t, "100", productPrice[0])

	productName, ok := vals["products[0][name]"]
	require.True(t, ok)
	assert.Equal(t, "Тестовый продукт", productName[0])

	productQuantity, ok := vals["products[0][quantity]"]
	require.True(t, ok)
	assert.Equal(t, "1", productQuantity[0])

	callbackType, ok := vals["callbackType"]
	require.True(t, ok)
	assert.Equal(t, "json", callbackType[0])

}

func TestProdamusGetPaymentLink(t *testing.T) {
	prodamus, _ := newProdamusSystem()
	result, err := prodamus.GetPaymentLink(context.Background(), GetPaymentLinkPayload{
		Login:       "andreyandreev",
		Amount:      100,
		Description: "Тестовый продукт",
		OrderId:     1337,
		Email:       "test@test.com",
		Phone:       "18123801133",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	t.Log(result.PaymentID, result.PaymentURL)
}
