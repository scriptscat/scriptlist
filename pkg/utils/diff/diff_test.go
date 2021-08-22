package diff

import (
	"fmt"
	"testing"
)

func TestDiff(t *testing.T) {
	ret := Diff("ok\nqwe123", "oo\nqwe123")
	fmt.Println(ret)
}
