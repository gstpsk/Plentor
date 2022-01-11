package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-zepto/zepto/web"
	"github.com/gstpsk/Plentor/db"
	"github.com/gstpsk/Plentor/models"
	"github.com/gstpsk/Plentor/util"
	"go.mongodb.org/mongo-driver/bson"
)

// Initialize session manager
var Sessions = make(map[string]map[string]string)

func LoginController(ctx web.Context) error {
	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	// Unmarshall JSON data
	var jdata map[string]string
	json.Unmarshal(buf, &jdata)
	email := jdata["email"]
	password := jdata["password"]

	// Form response struct
	type Message struct {
		Message string `json:"message"`
	}
	var respMsg = Message{}

	// Get hash from database
	col := db.GetUserCol()
	var filter bson.M = bson.M{"email": email}
	res := col.FindOne(context.Background(), filter)
	var user models.User
	err = res.Decode(&user)

	if err != nil { // no results
		respMsg.Message = "Invalid email or password!"
		return ctx.RenderJson(respMsg)

	}

	// Validate the hash
	if !util.HashIsValid(user.Hash, password) {
		respMsg.Message = "Invalid email or password!"
		return ctx.RenderJson(respMsg)
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
	Sessions[SessionId]["id"] = user.Id
	// Set session cookie
	cookie := http.Cookie{
		Name:   "SESSION-ID",
		Value:  SessionId,
		MaxAge: 0,
		Path:   "/",
		Secure: true, // prevent warning: “SameSite” attribute set to “None” or an invalid value, without the “secure” attribute.
	}
	http.SetCookie(ctx.Response(), &cookie)
	respMsg.Message = "success"
	return ctx.RenderJson(respMsg)
}

func RegisterController(ctx web.Context) error {
	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	// Form response struct
	type Message struct {
		Message string `json:"message"`
	}
	var respMsg = Message{}
	respMsg.Message = "success"

	// Unmarshall json
	var jdata map[string]string
	json.Unmarshal(buf, &jdata)
	password := jdata["password"]

	// Hash the plaintext password
	hash, err := util.GetBcrypt(password)
	util.Check(err)

	// Create object
	var user models.User
	user.Name = jdata["username"]
	user.Email = jdata["email"]
	user.Hash = hash

	// Check if user already exists
	col := db.GetUserCol()
	var filter bson.M = bson.M{"email": user.Email}
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to query database for user: %s", err)
	}
	var results []bson.M
	cur.All(context.Background(), &results)
	if len(results) > 0 { // User already exists
		respMsg.Message = "User already exists!"
	} else {
		// Insert user into database
		_, err = col.InsertOne(context.Background(), user)
		if err != nil {
			log.Fatalf("Failed to insert user into database: %s", err)
		}
	}

	return ctx.RenderJson(respMsg)
}

func LogoutPage(ctx web.Context) error {
	cookie, err := ctx.Request().Cookie("SESSION-ID")
	if err != nil {
		return ctx.Redirect("/login")
	}

	// CLear session
	Sessions[cookie.Value] = nil

	// Expire the cookie
	newCookie := http.Cookie{
		Name:    "SESSION-ID",
		Value:   cookie.Value,
		Expires: time.Unix(0, 0), // Jan 1st 1970s
		Path:    "/",
		Secure:  true, // prevent warning: “SameSite” attribute set to “None” or an invalid value, without the “secure” attribute.
	}
	http.SetCookie(ctx.Response(), &newCookie)

	return ctx.Redirect("/login")
}

func RequestIsAuthorized(ctx web.Context) bool {
	cookie, err := ctx.Request().Cookie("SESSION-ID")

	// ErrNoCookie aka they not allowed
	if err != nil {
		ctx.SetStatus(403)
		log.Println("No cookie found")
		return false
	}

	// TODO
	// Check if cookie value is in session manager
	// tbh this is kinda bad for security but
	// we'll fix it later. Not like ppl gonna
	// bruteforce the session id ;)
	// yes they will.
	if Sessions[cookie.Value] == nil {
		ctx.SetStatus(403)
		log.Printf("No session found for: %s\n", cookie.Value)
		return false
	}

	return true
}
