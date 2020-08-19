package calico

import (
	"fmt"
	"testing"
)

func TestShow(t *testing.T) {
	result, err := Show()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}
