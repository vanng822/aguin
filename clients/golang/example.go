
package main

import (
	"github.com/vanng822/aguin/clients/golang/aguin"
	"fmt"
)

func main() {
	t := aguin.New("545e0716f2fea0c7a9c46c74", "545e0716f2fea0c7a9c46c74fec46c71", "545e0716f2fea0c7a9c46c74fec46c71", "http://127.0.0.1:8080/")
	fmt.Println(t.Status())
	fmt.Println(t.Get("something", nil))
}