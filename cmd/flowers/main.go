package main

import "github.com/adelrosarioh/tarantulas/collectors/flowers"

func main() {

	if err := flowers.Run(); err != nil {
		panic(err)
	}
}
