package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/adelrosarioh/tarantulas/collectors/agricultural"
	"github.com/adelrosarioh/tarantulas/collectors/basicbaskets"
	"github.com/adelrosarioh/tarantulas/collectors/dairy"
	"github.com/adelrosarioh/tarantulas/collectors/flowers"
	"github.com/adelrosarioh/tarantulas/collectors/hardware"
	"github.com/adelrosarioh/tarantulas/collectors/medicines"
	"github.com/adelrosarioh/tarantulas/collectors/sirenado"
	"github.com/adelrosarioh/tarantulas/collectors/textbooks"
	"github.com/adelrosarioh/tarantulas/collectors/transportation"
)

type ScheduleEventInput struct {
	Collector string
}

func handler(ctx context.Context, event ScheduleEventInput) {

	collector := event.Collector

	switch collector {
	case "agricultural":
		if err := agricultural.Run(); err != nil {
			panic(err)
		}
	case "basicbaskets":
		if err := basicbaskets.Run(); err != nil {
			panic(err)
		}
	case "dairy":
		if err := dairy.Run(); err != nil {
			panic(err)
		}
	case "flowers":
		if err := flowers.Run(); err != nil {
			panic(err)
		}
	case "hardware":
		if err := hardware.Run(); err != nil {
			panic(err)
		}
	case "medicines":
		if err := medicines.Run(); err != nil {
			panic(err)
		}
	case "textbooks":
		if err := textbooks.Run(); err != nil {
			panic(err)
		}
	case "transportation":
		if err := transportation.Run(); err != nil {
			panic(err)
		}
	case "sirenado":
		if err := sirenado.Run(); err != nil {
			panic(err)
		}
	default:
		panic(errors.New("collector not found"))
	}

	fmt.Printf("%s collector is done.\n", collector)
}

func main() {
	lambda.Start(handler)
}
