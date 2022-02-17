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

	// Pages
	z.Get("/", controllers.IndexPage)
	z.Get("/login", controllers.LoginPage)
	z.Get("/logout", controllers.LogoutPage)
	z.Get("/register", controllers.RegisterPage)
	z.Get("/forgot", controllers.ForgotPage)
	z.Get("/dashboard", controllers.DashboardPage)
	z.Get("/event/new", controllers.NewEventPage)
	z.Get("/event/{id}", controllers.EventPage)
	z.Get("/event/{id}/signup", controllers.EventSignupPage)
	z.Get("/event/signup/{id}/success", controllers.EventSignupSuccessPage)

	// API routes

	// Auth
	z.Post("/api/login", controllers.LoginController)
	z.Post("/api/register", controllers.RegisterController)

	// Events
	z.Post("/api/event/new", controllers.NewEventController)
	z.Post("/api/event/{id}", controllers.UpdateEventController)
	z.Get("/api/events", controllers.EventsController)
	z.Get("/api/event/{id}", controllers.EventController)

	// Registrations
	z.Post("/api/registration/new", controllers.NewRegistrationController)
	z.Get("/api/registrations/{event_id}", controllers.RegistrationsController)
	z.Get("/api/registration/{id}", controllers.RegistrationController)

	z.Start()
}
