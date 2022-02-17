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
)

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
	_, err = col.InsertOne(context.Background(), newRegistration)
	if err != nil {
		log.Fatalf("Failed to insert registration into database: %s", err)
	}

	// Form a response struct
	type Message struct {
		Message string `json:"message"`
	}
	var respMsg = Message{}
	respMsg.Message = "success"

	// return
	return ctx.RenderJson(respMsg)
}

func RegistrationsController(ctx web.Context) error { // GET: /api/registrations/{event_id}
	// Get registration collection
	col := db.GetRegistrationCol()

	// Create objectid
	// objId, err := primitive.ObjectIDFromHex(ctx.Params()["event_id"])
	// if err != nil {
	// 	log.Printf("Failed to convert to ObjectID: %s", err)
	// 	return err
	// }
	// log.Printf("objid: %s", objId.Hex())
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
