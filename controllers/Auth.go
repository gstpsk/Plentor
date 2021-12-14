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
	buf, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Fatalf("Failed to read body to buffer: %s", err)
	}

	var jdata map[string]string
	json.Unmarshal(buf, &jdata)
	username := jdata["username"]
	password := jdata["password"]

	fmt.Printf("username: %s\npassword: %s\n", username, password)
	hash, err := util.GetBcrypt(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %s", err)
	}
	fmt.Println("hash: " + hash)

	return ctx.RenderJson(fmt.Sprintf("{'message': '%s'}", "peter"))
}
