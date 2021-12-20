package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-zepto/zepto/web"
	"github.com/gstpsk/Plentor/util"
)

// Initialize session manager
var Sessions = make(map[string]map[string]string)

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

	// Form response struct
	type Message struct {
		Message string `json:"message"`
	}
	var respMsg = Message{}

	// Validate the hash
	if !util.HashIsValid(hash, password) {
		respMsg.Message = "Invalid username or password!"
		str, err := json.Marshal(respMsg)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %s", err)
		}
		return ctx.RenderJson(string(str))
	}

	// Generate session id that doesn't exist
	var SessionId string
	for {
		SessionId = util.RandomString(32)
		if Sessions[SessionId] == nil {
			break
		}
	}

	// Initialize session map
	Sessions[SessionId] = make(map[string]string)
	// Add expiry timestamp
	Sessions[SessionId]["expires"] = time.Now().Add(time.Hour * 2).String()
	// Set session cookie
	cookie := http.Cookie{
		Name:   "SESSION-ID",
		Value:  SessionId,
		MaxAge: 300,
		Secure: true, // prevent warning: “SameSite” attribute set to “None” or an invalid value, without the “secure” attribute.
	}
	http.SetCookie(ctx.Response(), &cookie)

	respMsg.Message = "success"
	str, err := json.Marshal(respMsg)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}
	return ctx.RenderJson(str)
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
