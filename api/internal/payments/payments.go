package payments

import (
	"context"
)

type GetPaymentLinkPayload struct {
	Login           string      `json:"login"`
	Password        string      `json:"password"`
	Amount          uint64      `json:"amount"`
	Email           string      `json:"email"`
	Phone           string      `json:"phone"`
	Description     string      `json:"description"`
	OrderId         int64       `json:"order_id"`
	SendReceipt     bool        `json:"send_receipt"`
	ReceiptSettings interface{} `json:"receipt_settings"`
}

type GetPaymentLinkResult struct {
	PaymentURL string `json:"payment_url"`
	PaymentID  string `json:"payment_id"`
	OrderID    int64  `json:"order_id"`
}

type PaymentSystem interface {
	GetPaymentLink(ctx context.Context, payload GetPaymentLinkPayload) (*GetPaymentLinkResult, error)
}

func NewPaymentSystem(paymentSystemType string) PaymentSystem {
	switch paymentSystemType {
	case "tinkoff":
		return NewTinkoff()
	case "prodamus":
		return NewProdamus()
	}
	return nil
}
