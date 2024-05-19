package payments

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	TinkoffBaseURL = "https://securepay.tinkoff.ru/v2"
)

type ReceiptItem struct {
	Name     string `json:"Name"`
	Price    uint64 `json:"Price"`
	Quantity int    `json:"Quantity"`
	Amount   uint64 `json:"Amount"`
	Tax      string `json:"Tax"`
}

type Receipt struct {
	Email    string        `json:"Email"`
	Taxation string        `json:"Taxation"`
	Items    []ReceiptItem `json:"Items"`
}

type TinkoffInitPayload struct {
	TerminalKey string            `json:"TerminalKey"`
	Amount      uint64            `json:"Amount"`
	OrderId     string            `json:"OrderId"`
	Description string            `json:"Description"`
	DATA        map[string]string `json:"DATA"`
	Receipt     *Receipt          `json:"Receipt"`
	Token       string            `json:"Token"`
}

type TinkoffInitResponse struct {
	PaymentId   string `json:"PaymentId"`
	PaymentURL  string `json:"PaymentURL"`
	Success     bool   `json:"Success"`
	ErrorCode   string `json:"ErrorCode"`
	TerminalKey string `json:"TerminalKey"`
	Status      string `json:"Status"`
	OrderId     string `json:"OrderId"`
	Amount      int    `json:"Amount"`
	Message     string `json:"Message"`
	Details     string `json:"Details"`
}

func (p *TinkoffInitPayload) getValuesForToken() map[string]string {
	return map[string]string{
		"TerminalKey": p.TerminalKey,
		"Amount":      strconv.FormatUint(p.Amount, 10),
		"OrderId":     p.OrderId,
		"Description": p.Description,
	}
}

func (p *TinkoffInitPayload) sortValuesForToken(values map[string]string) []string {
	keys := make([]string, 0, len(values))

	for k, _ := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	sortedValues := make([]string, 0, len(values))

	for _, k := range keys {
		sortedValues = append(sortedValues, values[k])
	}

	return sortedValues
}

func (p *TinkoffInitPayload) GenerateToken(password string) string {
	// Чтобы сгенерировать токен — нужно проделать несколько шагов
	// Шаг 1. Из payload убрать объекты и массивы
	values := p.getValuesForToken()

	// Шаг 2. Добавить в эти полученные значения пароль от терминала
	values["Password"] = password

	// Шаг 3. Отсортировать по ключу в алфавитном порядке
	sortedValues := p.sortValuesForToken(values)

	// Шаг 4. Сделать конкатенацию значений в одну строку без пробелов
	hashValue := strings.Join(sortedValues, "")

	// Шаг 5. Сделать хэш sha256
	h := sha256.New()
	h.Write([]byte(hashValue))
	p.Token = hex.EncodeToString(h.Sum(nil))

	return p.Token
}

func (p *TinkoffInitPayload) updateAmount() {
	// тинькофф эквайринг проводит платежи в копейках
	p.Amount = p.Amount * 100
}

type Tinkoff struct{}

func (t *Tinkoff) GetPaymentLink(ctx context.Context, payload GetPaymentLinkPayload) (*GetPaymentLinkResult, error) {
	initPayload := TinkoffInitPayload{
		TerminalKey: payload.Login,
		Amount:      payload.Amount,
		OrderId:     strconv.Itoa(int(payload.OrderId)),
		Description: payload.Description,
		DATA: map[string]string{
			"email": payload.Email,
		},
	}

	initPayload.updateAmount()

	initPayload.GenerateToken(payload.Password)

	if payload.SendReceipt {
		initPayload.Receipt = &Receipt{
			Email: payload.Email,
			// TODO: fix taxation
			Taxation: "osn",
			Items: []ReceiptItem{
				{
					Name:     initPayload.Description,
					Price:    initPayload.Amount,
					Amount:   initPayload.Amount,
					Quantity: 1,
					Tax:      "none",
				},
			},
		}
	}

	jsonBody, err := json.Marshal(initPayload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(TinkoffBaseURL+"/Init", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	result := TinkoffInitResponse{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		fmt.Println(result)
		return nil, errors.New(result.Details)
	}

	return &GetPaymentLinkResult{
		PaymentID:  result.PaymentId,
		PaymentURL: result.PaymentURL,
		OrderID:    payload.OrderId,
	}, nil
}

func NewTinkoff() *Tinkoff {
	return &Tinkoff{}
}
