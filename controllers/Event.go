package controllers

import (
	"context"
	"encoding/json"
	"errors"
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
	if !RequestIsAuthorized(ctx) {
		return ctx.RenderJson("Forbidden")
	}

	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	// Create the Event object
	newEvent := models.Event{}

	// Unmarhall json into the Event object
	err = json.Unmarshal(buf, &newEvent)
	if err != nil {
		log.Fatalf("Failed to unmarshall JSON: %s", err)
	}

	// Set the user id
	cookie, err := ctx.Request().Cookie("SESSION-ID")
	newEvent.UserId = Sessions[cookie.Value]["id"]

	// Insert the Event object into the database
	col := db.GetEventCol()
	res, err := col.InsertOne(context.Background(), newEvent)
	if err != nil {
		log.Fatalf("Failed to insert event into database: %s", err)
	}

	// Form a response struct
	type Message struct {
		Message string `json:"message"`
		Id      string `json:"id"`
	}
	var respMsg = Message{}
	respMsg.Message = "success"
	// Give it the ObjectID of the previously created Event.
	respMsg.Id = res.InsertedID.(primitive.ObjectID).Hex()

	// return
	return ctx.RenderJson(respMsg)
}

func UpdateEventController(ctx web.Context) error { // POST: /api/event/{id}
	if !RequestIsAuthorized(ctx) {
		return ctx.RenderJson("Forbidden")
	}

	// Form response struct
	type Message struct {
		Message string `json:"message"`
	}
	var respMsg = Message{}
	respMsg.Message = "success"

	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Printf("Failed to read body to buffer: %s", err)
		ctx.SetStatus(500)
		respMsg.Message = "Server error"
		return ctx.RenderJson(respMsg)
	}

	// Create event object
	updatedEvent := models.Event{}
	// Unmarhall json into event object
	json.Unmarshal(buf, &updatedEvent)

	// Get events
	col := db.GetEventCol()

	// Create objectid
	objId, err := primitive.ObjectIDFromHex(ctx.Params()["id"])
	if err != nil {
		log.Printf("Failed to convert hex to ObjectId: %s", objId)
		ctx.SetStatus(500)
		respMsg.Message = "Server error"
		return ctx.RenderJson(respMsg)
	}

	// Create filter
	var filter bson.M = bson.M{"_id": objId}

	// Query
	sres := col.FindOne(context.Background(), filter)

	// Create event object and demarshall
	var event models.Event
	err = sres.Decode(&event)
	if err != nil {
		ctx.SetStatus(403)
		respMsg.Message = "Forbidden"
		return ctx.RenderJson(respMsg)
	}

	// Check if user id of original event matches user trying to update the event
	cookie, err := ctx.Request().Cookie("SESSION-ID")
	if event.UserId != Sessions[cookie.Value]["id"] {
		ctx.SetStatus(403)
		respMsg.Message = "Forbidden"
		return ctx.RenderJson(respMsg)
	}

	// Set user id on updated event
	updatedEvent.UserId = event.UserId

	// Update event object in database
	ures, err := col.ReplaceOne(context.Background(), filter, updatedEvent)
	if err != nil {
		log.Printf("Failed to update event in database: %s", err)
		ctx.SetStatus(500)
		respMsg.Message = "Failed to update event"
	}
	if !(ures.ModifiedCount > 0) {
		ctx.SetStatus(500)
		respMsg.Message = "Failed to update event"
	}

	return ctx.RenderJson(respMsg)
}

// Returns all events for a specific user
func EventsController(ctx web.Context) error { // GET: /api/events
	if !RequestIsAuthorized(ctx) {
		ctx.SetStatus(403)
		return ctx.RenderJson("Forbidden")
	}

	// Get cookie
	cookie, err := ctx.Request().Cookie("SESSION-ID")

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
	// TODO
	// This also rly bad for security
	// But people registering need to know
	// the event details
	if false {
		ctx.SetStatus(403)
		return ctx.RenderJson("Forbidden")
	}

	// Get events
	col := db.GetEventCol()

	// Create objectid
	objId, err := primitive.ObjectIDFromHex(ctx.Params()["id"])
	if err != nil {
		log.Printf("Failed to convert hex to ObjectId: %s", objId)
		ctx.SetStatus(500)
		return ctx.RenderJson("Server error")
	}

	// Create filter
	var filter bson.M = bson.M{"_id": objId}

	// Query
	res := col.FindOne(context.Background(), filter)

	// Create event objet and demarshall
	var event models.Event
	err = res.Decode(&event)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.Redirect("/dashboard")
		} else {
			log.Printf("Failed to decode event: %s", err)
			ctx.SetStatus(500)
			return ctx.RenderJson("Server error")
		}
	}

	return ctx.RenderJson(event)
}
