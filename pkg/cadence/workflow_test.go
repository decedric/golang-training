package cadence

import (
	"math/big"
	"testing"
)

func TestFibonacci(t *testing.T) {
	n_arr := [5]uint{0, 1, 5, 10, 100}
	want_arr := [5]big.Int{}
	want_arr[0].SetString("0", 10)
	want_arr[0].SetString("1", 10)
	want_arr[0].SetString("5", 10)
	want_arr[0].SetString("55", 10)
	want_arr[0].SetString("354224848179261915075", 10)
	for i := 0; i < 5; i++ {
		actual := fibonacci(n_arr[i])
		if want_arr[i].Cmp(actual) == 0 {
			t.Errorf("got %s, wanted %s", actual.Text(10), want_arr[i].Text(10))
		}
	}
}
