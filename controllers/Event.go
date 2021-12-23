package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/go-zepto/zepto/web"
	"github.com/gstpsk/Plentor/db"
	"github.com/gstpsk/Plentor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewEventController(ctx web.Context) error { // POST: /api/event/new
	cookie, err := ctx.Request().Cookie("SESSION-ID")

	// ErrNoCookie aka they not allowed
	if err != nil {
		ctx.SetStatus(403)
		fmt.Println("no cookie found")
		return ctx.RenderJson("Forbidden")
	}

	// Check if cookie value is in session manager
	if Sessions[cookie.Value] == nil {
		ctx.SetStatus(403)
		fmt.Println("no session found")
		return ctx.RenderJson("Forbidden")
	}

	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	// Create event object
	newEvent := models.Event{}
	// Unmarhall json into event object
	json.Unmarshal(buf, &newEvent)
	// set user id
	newEvent.UserId = Sessions[cookie.Value]["id"]

	// Insert event object into database
	col := db.GetEventCol()
	res, err := col.InsertOne(context.Background(), newEvent)
	if err != nil {
		log.Fatalf("Failed to insert event into database: %s", err)
	}

	// Form response struct
	type Message struct {
		Message string `json:"message"`
		Id      string `json:"id"`
	}
	var respMsg = Message{}
	respMsg.Message = "success"
	respMsg.Id = res.InsertedID.(primitive.ObjectID).Hex()

	// return
	return ctx.RenderJson(respMsg)
}

func EventsController(ctx web.Context) error { // GET: /api/events
	cookie, err := ctx.Request().Cookie("SESSION-ID")

	// ErrNoCookie aka they not allowed
	if err != nil {
		ctx.SetStatus(403)
		fmt.Println("no cookie found")
		return ctx.RenderJson("Forbidden")
	}

	// Check if cookie value is in session manager
	if Sessions[cookie.Value] == nil {
		ctx.SetStatus(403)
		fmt.Println("no session found")
		return ctx.RenderJson("Forbidden")
	}

	// Get all events for user
	col := db.GetEventCol()
	var filter bson.M = bson.M{"userid": Sessions[cookie.Value]["id"]}
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to query events for user %s: %s", Sessions[cookie.Value]["id"], err)
	}
	var results []models.Event
	cur.All(context.Background(), &results)

	return ctx.RenderJson(results)
}

// EventController: returns data from one specific event
func EventController(ctx web.Context) error { // GET: /api/event/{id}
	cookie, err := ctx.Request().Cookie("SESSION-ID")

	// ErrNoCookie aka they not allowed
	if err != nil {
		ctx.SetStatus(403)
		fmt.Println("no cookie found")
		return ctx.RenderJson("Forbidden")
	}

	// Check if cookie value is in session manager
	// tbh this is kinda bad for security but
	// we'll fix it later. Not like ppl gonna
	// bruteforce the session id ;)
	// yes they will.
	if Sessions[cookie.Value] == nil {
		ctx.SetStatus(403)
		fmt.Println("no session found")
		return ctx.RenderJson("Forbidden")
	}

	col := db.GetEventCol()
	objId, err := primitive.ObjectIDFromHex(ctx.Params()["id"])
	if err != nil {
		log.Fatalf("Failed to convert hex to ObjectId: %s", objId)
	}
	var filter bson.M = bson.M{"_id": objId}
	res := col.FindOne(context.Background(), filter)
	var event models.Event
	err = res.Decode(&event)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.Redirect("/dashboard")
		} else {
			log.Fatalf("Failed to decode event: %s", err)
		}
	}

	return ctx.RenderJson(event)
}
