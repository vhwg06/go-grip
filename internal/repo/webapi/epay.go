package webapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type PaymentVerifier interface {
	Verify(ctx context.Context, payload map[string]string) (bool, error)
	PaymentID(payload map[string]string) string
}

type EpayVerifier struct {
	secret string
}

func NewEpayVerifier(secret string) *EpayVerifier {
	return &EpayVerifier{secret: secret}
}

func (v *EpayVerifier) Verify(ctx context.Context, payload map[string]string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("EpayVerifier.Verify: %w", err)
	}

	if v.secret == "" {
		return false, nil
	}

	orderID := payload["order_id"]
	paymentID := payload["payment_id"]
	signature := strings.TrimSpace(payload["sign"])
	if orderID == "" || paymentID == "" || signature == "" {
		return false, nil
	}

	mac := hmac.New(sha256.New, []byte(v.secret))
	mac.Write([]byte(orderID + ":" + paymentID))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(strings.ToLower(signature)), []byte(strings.ToLower(expected))), nil
}

func (v *EpayVerifier) PaymentID(payload map[string]string) string {
	if id := payload["payment_id"]; id != "" {
		return id
	}
	return payload["trade_no"]
}
