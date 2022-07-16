package eth

import (
	"math/big"
	"testing"
)

func TestStringToWei(t *testing.T) {
	qty := "0.11"
	qtyInWei := EthStringToWei(qty)
	expected := big.NewInt(110000000000000000)
	if qtyInWei.Cmp(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, qtyInWei)
	}
}
