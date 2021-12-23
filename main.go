package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/go-zepto/zepto"
	"github.com/gstpsk/Plentor/controllers"
	"github.com/gstpsk/Plentor/util"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load the config file
	util.LoadConfig()

	// Initialize the database
	//os.Remove("./db/dev.db")
	file, err := os.Open("./db/dev.db")
	if errors.Is(err, os.ErrNotExist) {
		db, err := sql.Open("sqlite3", "./db/dev.db")
		util.Check(err)
		defer db.Close() // good practice mate

		sqlStmt := `
			CREATE TABLE users (id integer not null primary key, username text, email text not null unique, hash text not null)
		
		`
		_, err = db.Exec(sqlStmt)
		util.Check(err)
	}
	file.Close()

	// Create Zepto
	z := zepto.NewZepto()

	// Routes
	z.Get("/", controllers.HelloIndex)
	z.Get("/login", controllers.LoginPage)
	z.Get("/register", controllers.RegisterPage)
	z.Get("/forgot", controllers.ForgotPage)
	z.Get("/dashboard", controllers.DashboardPage)
	z.Get("/event/new", controllers.NewEventPage)

	// API routes
	z.Post("/api/login", controllers.LoginController)
	z.Post("/api/register", controllers.RegisterController)
	z.Post("/api/event/new", controllers.NewEventController)

	z.Start()
}
