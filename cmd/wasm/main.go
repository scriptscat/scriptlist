//+build wasm

package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"hash"

	"github.com/gopherjs/gopherjs/js"
)

import (
	"syscall/js"
)

// Secret wasm,给前端使用
const Secret = "NQ3kDBBjRmBpRHSX3"

func main() {
	js.Global().Set("statistics", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//id := args[0].Int()
		token := args[1].String()
		h := hmac.New(func() hash.Hash {
			return sha1.New()
		}, []byte(Secret))
		//h.Write([]byte(strconv.Itoa(id)))
		h.Write([]byte(token))
		return js.ValueOf(fmt.Sprintf("%x", h.Sum(nil)))
	}))
}
