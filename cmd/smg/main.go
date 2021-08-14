package main

import (
	"fmt"

	"github.com/jarviliam/smg"
)

func main() {
	spider := smg.NewSpider()
	err := spider.Fetch("https://www.reddit.com/r/neovim")
	if err != nil {
		fmt.Println(err)
	}
}
