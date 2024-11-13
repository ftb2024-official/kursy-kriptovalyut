package entities

import (
	"testing"
	"time"
)

type newCoinTest struct {
	arg1 string
	arg2 float64
	res1 string
	res2 float64
}

var newCoinTests = []newCoinTest{
	{"", 0, "BTC", 0},
	{"BTC", -1, "BTC", 0},
	{"ETH", 1, "ETH", 1},
	{"USDT", 2, "BTC", 2},
	{"smth else", 3, "BTC", 3},
}

func TestNewCoin(t *testing.T) {
	for _, test := range newCoinTests {
		out, _ := NewCoin(test.arg1, test.arg2)
		if out.title != test.res1 || out.price != test.res2 {
			t.Errorf("got: %v, %v, want: %v, %v", out.title, out.price, test.res1, test.res2)
		}

		// Проверка, что время создания актуально (с небольшим допустимым отклонением)
		if time.Since(out.actualAt) > time.Second {
			t.Errorf("actualAt is too old: got %v", out.actualAt)
		}
	}
}

// for _, test := range []struct {
// 	arg1 string
// 	arg2 float64
// 	res1 string
// 	res2 float64
// }{
// 	{"", 0, "BTC", 0},
// 	{"BTC", -1, "BTC", 0},
// 	{"ETH", 1, "ETH", 1},
// 	{"USDT", 2, "BTC", 2},
// 	{"smth else", 3, "BTC", 3},
// } {
// 	out := NewCoin(test.arg1, test.arg2)
// 	if out.title != test.res1 || out.price != test.res2 {
// 		t.Errorf("got: %v, %v, want: %v, %v", out.title, out.price, test.res1, test.res2)
// 	}

// 	// Проверка, что время создания актуально (с небольшим допустимым отклонением)
// 	if time.Since(out.actualAt) > time.Second {
// 		t.Errorf("actualAt is too old: got %v", out.actualAt)
// 	}
// }
