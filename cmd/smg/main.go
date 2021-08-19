package main

import (
	"fmt"

	"github.com/jarviliam/smg"
)

func main() {
	smg := smg.NewSMG()
	err := smg.Run("http://localhost:8080/pokesp/servlet/test/sp/index.html")
	if err != nil {
		fmt.Println(err)
	}
}
