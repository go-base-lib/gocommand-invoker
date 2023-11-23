//go:build unix || (js && wasm) || plan9 || wasip1

package gocommandinvoker

import (
	"fmt"
	"testing"
)

func TestFindCommandPath(t *testing.T) {
	p, err := FindCommandPath("~/applications/bin/suwellCoverRpcAndConfig")
	fmt.Println(p, err)

	fmt.Println(0656 & 0111)
}
