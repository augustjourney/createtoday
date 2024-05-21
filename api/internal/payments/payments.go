package payments

import (
	"context"
	"strings"
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

const StatusSucceeded = "succeeded"
const StatusCanceled = "canceled"
const StatusRejected = "rejected"
const StatusPending = "pending"
const StatusExpired = "expired"

var statuses map[string]string = map[string]string{
	StatusSucceeded: "CONFIRMED,Completed,succeeded,success",
	StatusCanceled:  "canceled,order_canceled",
	StatusRejected:  "REJECTED,Declined,order_denied",
	StatusExpired:   "DEADLINE_EXPIRED",
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

func FormatStatus(status string) string {
	// У всех платежных систем свои статусы
	// Здесь приводим их в одну единую систему

	if status == "" {
		return StatusPending
	}

	for needStatus, possibleStatus := range statuses {
		if strings.Contains(strings.ToLower(possibleStatus), strings.ToLower(status)) {
			return needStatus
		}
	}

	return StatusPending
}
