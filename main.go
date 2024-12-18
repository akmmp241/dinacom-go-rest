package main

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/app"
)

func main() {
	config := app.NewConfig()

	db := app.NewDB(config)

	fiberApp := app.NewRouter()

	if err := fiberApp.Listen(":3000"); err != nil {
		panic(err)
	}
}
