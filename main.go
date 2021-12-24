package main

import (
	"github.com/go-zepto/zepto"
	"github.com/gstpsk/Plentor/controllers"
	"github.com/gstpsk/Plentor/util"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load the config file
	util.LoadConfig()

	// Create Zepto
	z := zepto.NewZepto()

	// Routes
	z.Get("/", controllers.HelloIndex)
	z.Get("/login", controllers.LoginPage)
	z.Get("/logout", controllers.LogoutPage)
	z.Get("/register", controllers.RegisterPage)
	z.Get("/forgot", controllers.ForgotPage)
	z.Get("/dashboard", controllers.DashboardPage)
	z.Get("/event/new", controllers.NewEventPage)
	z.Get("/event/{id}", controllers.EventPage)

	// API routes
	z.Post("/api/login", controllers.LoginController)
	z.Post("/api/register", controllers.RegisterController)
	z.Post("/api/event/new", controllers.NewEventController)
	z.Get("/api/events", controllers.EventsController)
	z.Get("/api/event/{id}", controllers.EventController)

	z.Start()
}
