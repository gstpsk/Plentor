package main

import (
	"database/sql"
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
	os.Remove("./db/dev.db")
	db, err := sql.Open("sqlite3", "./db/dev.db")
	util.Check(err)
	defer db.Close() // good practice mate

	sqlStmt := `
		CREATE TABLE users (id integer not null primary key, username text not null unique, hash text not null)
	
	`
	_, err = db.Exec(sqlStmt)
	util.Check(err)

	//tx, err := db.Begin()
	util.Check(err)

	// Create Zepto
	z := zepto.NewZepto()

	// Routes
	z.Get("/", controllers.HelloIndex)
	z.Get("/login", controllers.LoginPage)
	z.Post("/api/login", controllers.LoginController)
	z.Post("/api/register", controllers.RegisterController)

	z.Start()
}
