package main

import (
	"context"
	"log"
	"os"

	"github.com/yyh-gl/ofukuro/butler"
	"github.com/yyh-gl/ofukuro/task"
)

func main() {
	b := butler.CallButler()

	humidityNotification := task.NewHumidityNotification()
	healthCheckNotification := task.NewHealthCheckNotification()

	ctx := context.Background()
	b.AddTask(ctx, humidityNotification)
	b.AddTask(ctx, healthCheckNotification)

	if err := b.StartWorking(ctx); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
