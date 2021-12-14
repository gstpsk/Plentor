package main

import (
	"github.com/gstpsk/Plentor/controllers"
	"github.com/go-zepto/zepto"
)

func main() {
	// Create Zepto
	z := zepto.NewZepto()

	// Routes
	z.Get("/", controllers.HelloIndex)

	z.Start()
}
