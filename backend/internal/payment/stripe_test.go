package payment

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v82"
)

func TestMapPaymentIntentStatus(t *testing.T) {
	tests := []struct {
		input stripe.PaymentIntentStatus
		want  string
	}{
		{stripe.PaymentIntentStatusSucceeded, "paid"},
		{stripe.PaymentIntentStatusCanceled, "failed"},
		{stripe.PaymentIntentStatusProcessing, "pending"},
		{stripe.PaymentIntentStatusRequiresAction, "pending"},
		{stripe.PaymentIntentStatusRequiresCapture, "pending"},
		{stripe.PaymentIntentStatusRequiresConfirmation, "pending"},
		{stripe.PaymentIntentStatusRequiresPaymentMethod, "pending"},
	}
	for _, tt := range tests {
		got := mapPaymentIntentStatus(tt.input)
		if got != tt.want {
			t.Errorf("mapPaymentIntentStatus(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMapRefundStatus(t *testing.T) {
	tests := []struct {
		input stripe.RefundStatus
		want  string
	}{
		{stripe.RefundStatusSucceeded, "success"},
		{stripe.RefundStatusFailed, "failed"},
		{stripe.RefundStatusPending, "pending"},
		{stripe.RefundStatusRequiresAction, "pending"},
	}
	for _, tt := range tests {
		got := mapRefundStatus(tt.input)
		if got != tt.want {
			t.Errorf("mapRefundStatus(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestCentsFromCNY(t *testing.T) {
	tests := []struct {
		cny   string
		cents int64
	}{
		{"1.00", 100},
		{"0.01", 1},
		{"99.99", 9999},
		{"100.50", 10050},
	}
	for _, tt := range tests {
		d, _ := decimal.NewFromString(tt.cny)
		got := centsFromCNY(d)
		if got != tt.cents {
			t.Errorf("centsFromCNY(%s) = %d, want %d", tt.cny, got, tt.cents)
		}
	}
}

func TestDecimalFromCents(t *testing.T) {
	tests := []struct {
		cents int64
		want  string
	}{
		{100, "1"},
		{1, "0.01"},
		{9999, "99.99"},
		{10050, "100.5"},
	}
	for _, tt := range tests {
		got := decimalFromCents(tt.cents)
		want, _ := decimal.NewFromString(tt.want)
		if !got.Equal(want) {
			t.Errorf("decimalFromCents(%d) = %s, want %s", tt.cents, got, tt.want)
		}
	}
}
