package main

import "github.com/adelrosarioh/tarantulas/collectors/transportation"

func main() {

	if err := transportation.Run(); err != nil {
		panic(err)
	}
}
