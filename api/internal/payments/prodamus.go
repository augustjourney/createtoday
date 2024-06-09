package payments

import (
	"context"
	"createtodayapi/internal/logger"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

type ProdamusInitPayload struct {
	Amount      uint64 `json:"amount"`
	OrderId     int64  `json:"orderId"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type Prodamus struct{}

func (t *Prodamus) GetPaymentLink(ctx context.Context, payload GetPaymentLinkPayload) (*GetPaymentLinkResult, error) {
	var result GetPaymentLinkResult

	params := t.generateQueryParams(ProdamusInitPayload{
		Amount:      payload.Amount,
		OrderId:     payload.OrderId,
		Email:       payload.Email,
		Phone:       payload.Phone,
		Description: payload.Description,
	})

	base, err := url.Parse(fmt.Sprintf("https://%s.payform.ru", payload.Login))
	if err != nil {
		logger.Error(ctx, "error parsing url", "err", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		logger.Error(ctx, "error creating prodamus payment link", "err", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "text/plain")
	req.URL.RawQuery = params

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(ctx, "error requesting prodamus payment link", "err", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, "error reading prodamus response body", "err", err)
		return nil, err
	}

	paymentId, err := uuid.NewRandom()
	if err != nil {
		logger.Error(ctx, "error generating prodamus payment id", "err", err)
		return nil, err
	}

	result.PaymentURL = string(body)
	result.PaymentID = paymentId.String()
	result.OrderID = payload.OrderId

	return &result, nil
}

func (t *Prodamus) generateQueryParams(data ProdamusInitPayload) string {
	q := url.Values{}

	q.Add("do", "link")
	q.Add("order_id", strconv.FormatInt(data.OrderId, 10))
	q.Add("customer_email", data.Email)
	q.Add("customer_phone", data.Phone)
	q.Add("products[0][name]", data.Description)
	q.Add("products[0][price]", strconv.FormatUint(data.Amount, 10))
	q.Add("products[0][quantity]", "1")
	q.Add("callbackType", "json")

	return q.Encode()
}

func NewProdamus() *Prodamus {
	return &Prodamus{}
}
