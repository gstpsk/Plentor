package main

import (
	"github.com/go-zepto/zepto"
	"github.com/gstpsk/Plentor/controllers"
	"github.com/gstpsk/Plentor/util"
)

func main() {
	util.LoadConfig()

	// Create Zepto
	z := zepto.NewZepto()

	// Routes
	z.Get("/", controllers.HelloIndex)
	z.Post("/api/login", controllers.LoginController)

	z.Start()
}
