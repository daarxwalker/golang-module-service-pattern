package main

import (
	"fmt"

	"example/app/module"
	"example/app/module/categoryModule"
	"example/app/service"
	"example/core"
)

func main() {
	app := core.NewApp()

	// Set services container
	app.SetContainer(service.NewContainer())

	// Register modules
	app.RegisterModule(module.Category, categoryModule.Module)

	appConfig := app.GetConfig().GetApp()
	fmt.Printf("%s running - http://localhost:%d\n", appConfig.Name, appConfig.Port)
	app.Start()
}
