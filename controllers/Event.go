package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/go-zepto/zepto/web"
	"github.com/gstpsk/Plentor/models"
	"github.com/gstpsk/Plentor/util"
)

func EventController(ctx web.Context) error {
	// Read request body to buffer and unmarshall
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	newEvent := models.Event{}

	json.Unmarshal(buf, &newEvent)

	// Save to the database
	db := util.DatabaseConnect()
	defer db.Close() // good practice mate, can't let 'em linger
	//err = util.InsertUser(db, username, email, hash)

	/*


		// Form response struct
		type Message struct {
			Message string `json:"message"`
		}
		var respMsg = Message{}
		respMsg.Message = "success"

		if err != nil {
			log.Printf("Failed to register user: %s", err.Error())
			if err.Error() == "UNIQUE constraint failed: users.email" { // bad way to check i know but idek how else i could do it
				respMsg.Message = "User already exists!"
			}
			ctx.SetStatus(500) // internal server error bro
			respMsg.Message = "Internal server error"
		}

		str, err := json.Marshal(respMsg)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %s", err)
		}*/
	fmt.Println(newEvent.Timeblocks[0].Date)
	return ctx.RenderJson("{}")
}
