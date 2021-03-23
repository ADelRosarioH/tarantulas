package main

import (
	"github.com/adelrosarioh/tarantulas/collectors"
)

func main() {

	basicbasket := &collectors.BasicBasket{}

	if err := basicbasket.Run(); err != nil {
		panic(err)
	}
}
