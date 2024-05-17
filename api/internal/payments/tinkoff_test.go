package payments

import (
	"context"
	"createtodayapi/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func newTinkoffSystem() (*Tinkoff, *config.Config) {
	tinkoff := Tinkoff{}
	conf := config.New("../../.env")
	return &tinkoff, conf
}

func TestGetValuesForToken(t *testing.T) {
	t.Parallel()
	t.Run("should not have data and receipt in values for token", func(t *testing.T) {
		payload := TinkoffInitPayload{
			TerminalKey: "98234234DEMO",
			DATA: map[string]string{
				"email": "test@example.com",
			},
			Description: "test order",
			Amount:      2500,
			OrderId:     "101231",
			Receipt: &Receipt{
				Email:    "test@example.com",
				Taxation: "usn",
				Items: []ReceiptItem{
					ReceiptItem{
						Name:   "Test Product",
						Price:  100,
						Amount: 100,
						Tax:    "none",
					},
				},
			},
		}

		values := payload.getValuesForToken()
		_, DataFieldExist := values["DATA"]
		_, ReceiptFieldExist := values["Receipt"]

		assert.False(t, DataFieldExist)
		assert.False(t, ReceiptFieldExist)
	})
}

func TestGenerateToken(t *testing.T) {
	t.Parallel()
	t.Run("should generate token", func(t *testing.T) {
		payload := TinkoffInitPayload{
			TerminalKey: "98234234DEMO",
			DATA: map[string]string{
				"email": "test@example.com",
			},
			Description: "test order",
			Amount:      2500,
			OrderId:     "101231",
			Receipt: &Receipt{
				Email:    "test@example.com",
				Taxation: "usn",
				Items: []ReceiptItem{
					ReceiptItem{
						Name:   "Test Product",
						Price:  2500,
						Amount: 2500,
						Tax:    "none",
					},
				},
			},
		}
		payload.GenerateToken("secret-123")
		token := "258e262db188f3e6bd13cb7231742392d78de6f4c8dcfdd98f6f11f0f934bdea"

		assert.Equal(t, payload.Token, token)
	})
}

func TestAmountUpdating(t *testing.T) {
	t.Parallel()
	t.Run("should update amount correctly", func(t *testing.T) {
		var amount uint64 = 2400

		payload := TinkoffInitPayload{
			TerminalKey: "98234234DEMO",
			DATA: map[string]string{
				"email": "test@example.com",
			},
			Description: "test order",
			Amount:      amount,
			OrderId:     "101231",
			Receipt: &Receipt{
				Email:    "test@example.com",
				Taxation: "usn",
				Items: []ReceiptItem{
					ReceiptItem{
						Name:   "Test Product",
						Price:  amount,
						Amount: amount,
						Tax:    "none",
					},
				},
			},
		}

		payload.updateAmount()

		assert.Equal(t, payload.Amount, amount*100, "amount should be updated correctly")
	})
}

func TestGetPaymentLink(t *testing.T) {
	t.Parallel()
	tinkoff, conf := newTinkoffSystem()
	t.Run("should not get payment link", func(t *testing.T) {
		result, err := tinkoff.GetPaymentLink(context.Background(), GetPaymentLinkPayload{
			Email:       "test@example.com",
			Description: "Test Product",
			Amount:      2500,
			OrderId:     "101231",
			Phone:       "13812345678",
			SendReceipt: false,
			Login:       "123123123DEMO",
			Password:    "secret-123",
		})
		require.Error(t, err)
		assert.Nil(t, result, "result should be nil")
	})

	t.Run("should get payment link", func(t *testing.T) {
		result, err := tinkoff.GetPaymentLink(context.Background(), GetPaymentLinkPayload{
			Email:       "test@example.com",
			Description: "Test Product",
			Amount:      2500,
			OrderId:     "101231",
			Phone:       "13812345678",
			SendReceipt: false,
			Login:       conf.TinkoffTestLogin,
			Password:    conf.TinkoffTestPassword,
		})
		require.NoError(t, err)
		assert.NotNil(t, result, "should get link")
		t.Log(result.PaymentID, result.PaymentURL)
		if result != nil {
			assert.NotEqual(t, result.PaymentURL, "", "payment url should not be empty")
			assert.NotEqual(t, result.PaymentID, "", "payment id should not be empty")
		}
	})
}
