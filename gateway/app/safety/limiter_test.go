package safety

import (
	"fmt"
	"testing"
)

func TestLimiter_Pass(t *testing.T) {
	l := Init(100, 5).Start()
	for i := 0; i < 10000; i++ {
		if !l.Pass() {

		} else {
			fmt.Println("token.....................",len(l.tokenChan))
		}
	}
}
