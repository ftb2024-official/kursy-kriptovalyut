package entities

import (
	"errors"
	"testing"
)

func TestNewCoin(t *testing.T) {
	// testcase 1: valid input
	coin, err := NewCoin("BTC", 1000)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if coin.title != "BTC" || coin.price != 1000 {
		t.Fatalf("unexpected Coin values: %v", coin)
	}

	// testcase 2: empty title
	_, err = NewCoin("", 1000)
	if err == nil {
		t.Fatalf("expected error for empty title, got nil")
	}
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected wrapped error %v, got %v", ErrEmptyTitle, err)
	}

	// testcase 3: negative price
	_, err = NewCoin("ETH", -1000)
	if err == nil {
		t.Fatalf("expected error for negative price, got nil")
	}
	if !errors.Is(err, ErrNegativePrice) {
		t.Fatalf("expected wrapped error '%v', got '%v'", ErrNegativePrice, err)
	}

	// testcase 4: zero price
	_, err = NewCoin("BTC", 0)
	if err == nil {
		t.Fatalf("expected error for zero price, got nil")
	}
	if !errors.Is(err, ErrZeroPrice) {
		t.Fatalf("expected wrapped error '%v', got '%v'", ErrZeroPrice, err)
	}
}
