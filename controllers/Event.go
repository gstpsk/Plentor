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

	// Create event object
	newEvent := models.Event{}
	// Unmarhall json into event object
	json.Unmarshal(buf, &newEvent)
	// set user id
	cookie, err := ctx.Request().Cookie("SESSION-ID")
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
		log.Fatalf("Failed to read body to buffer: %s", err)
	}
	// Create event object
	updatedEvent := models.Event{}
	// Unmarhall json into event object
	json.Unmarshal(buf, &updatedEvent)
	// set user id
	cookie, err := ctx.Request().Cookie("SESSION-ID")
	updatedEvent.UserId = Sessions[cookie.Value]["id"]

	// Update event object in database
	col := db.GetEventCol()
	objId, err := primitive.ObjectIDFromHex(ctx.Params()["id"])
	var filter bson.M = bson.M{"_id": objId}
	res, err := col.ReplaceOne(context.Background(), filter, updatedEvent)
	if err != nil {
		log.Fatalf("Failed to update event in database: %s", err)
	}
	if !(res.ModifiedCount > 0) {
		respMsg.Message = "Failed to update event"
	}

	return ctx.RenderJson(respMsg)
}

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
	if !RequestIsAuthorized(ctx) {
		ctx.SetStatus(403)
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
