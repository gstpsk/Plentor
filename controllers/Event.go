package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/go-zepto/zepto/web"
	"github.com/gstpsk/Plentor/db"
	"github.com/gstpsk/Plentor/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewEventController(ctx web.Context) error { // /api/event/new
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
