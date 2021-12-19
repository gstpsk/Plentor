package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/go-zepto/zepto/web"
	"github.com/gstpsk/Plentor/util"
)

func LoginController(ctx web.Context) error {
	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	var jdata map[string]string
	json.Unmarshal(buf, &jdata)
	username := jdata["username"]
	password := jdata["password"]

	// Get hash from database
	db := util.DatabaseConnect()
	defer db.Close()
	hash, err := util.GetHashByUsername(db, username)
	if err != nil {
		log.Fatalf("Failed to retrieve hash from database: %s", err)
	}
	// Validate the hash
	if !util.HashIsValid(hash, password) {
		return ctx.RenderJson(fmt.Sprintf("{'message': '%s'}", "Invalid username or password!"))
	}

	return ctx.RenderJson(fmt.Sprintf("{'message': '%s'}", "success"))
}

func RegisterController(ctx web.Context) error {
	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	var jdata map[string]string
	json.Unmarshal(buf, &jdata)
	username := jdata["username"]
	password := jdata["password"]

	// Hash the plaintext password
	hash, err := util.GetBcrypt(password)
	util.Check(err)

	// Save to the database
	db := util.DatabaseConnect()
	defer db.Close() // good practice mate, can't let 'em linger
	err = util.InsertUser(db, username, hash)

	if err != nil {
		log.Printf("Failed to register user: %s", err.Error())
		if err.Error() == "UNIQUE constraint failed: users.username" { // bad way to check i know but idek how else i could do it
			return ctx.RenderJson(fmt.Sprintf("{'message': '%s'}", "User already exists!"))
		}
		ctx.SetStatus(500) // internal server error bro
		return ctx.RenderJson(fmt.Sprintf("{'message': '%s'}", "failed"))
	}

	return ctx.RenderJson(fmt.Sprintf("{'message': '%s'}", "success"))
}
