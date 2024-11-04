package monkey

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

func TestCompute(t *testing.T) {
	patches := gomonkey.ApplyFunc(networkCompute, func(a, b int) (int, error) {
		return 2, nil
	})
	defer patches.Reset()

	sum, err := Compute(1, 2)
	if err != nil {
		t.Error(err)
	}
	if sum != 3 {
		t.Errorf("sum is %d, want 3", sum)
	}
}
