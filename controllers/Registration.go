package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/go-zepto/zepto/web"
	"github.com/gstpsk/Plentor/db"
	"github.com/gstpsk/Plentor/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a new registration
func NewRegistrationController(ctx web.Context) error { // POST: /api/registration/new
	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	// Create a new Registration object
	newRegistration := models.Registration{}

	// Unmarshall json into the Registration object
	err = json.Unmarshal(buf, &newRegistration)
	if err != nil {
		log.Fatalf("Failed to unmarshall JSON: %s", err)
	}

	// Insert the Registration object into the database
	col := db.GetRegistrationCol()
	res, err := col.InsertOne(context.Background(), newRegistration)
	if err != nil {
		log.Fatalf("Failed to insert registration into database: %s", err)
	}

	// Form a response struct
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

// Get all registrations for a certain event
func RegistrationsController(ctx web.Context) error { // GET: /api/registrations/{event_id}
	// Get registration collection
	col := db.GetRegistrationCol()

	// Form filter
	var filter bson.M = bson.M{"eventid": ctx.Params()["event_id"]}

	// Search
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to query for registration with certain eventid: %s", err)
	}

	var results []models.Registration
	cur.All(context.Background(), &results)

	return ctx.RenderJson(results)
}

// Retrieve one registration
func RegistrationController(ctx web.Context) error { // GET: /api/registration/{id}
	// Get registration collection
	col := db.GetRegistrationCol()

	// Convert hex to ObjectID
	objId, err := primitive.ObjectIDFromHex(ctx.Params()["id"])
	if err != nil {
		log.Fatalf("Failed to convert hex to ObjectId: %s", objId)
	}

	// Form filter
	var filter bson.M = bson.M{"_id": objId}

	// Search
	res := col.FindOne(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to query for registration: %s", err)
	}

	// Create registration object
	var result models.Registration
	// Unmarshall into object
	res.Decode(&result)

	return ctx.RenderJson(result)
}

func ICalController(ctx web.Context) error {
	return nil
}
