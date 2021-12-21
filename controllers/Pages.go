package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-zepto/zepto/web"
)

func LoginPage(ctx web.Context) error {
	return ctx.Render("pages/login")
}

func RegisterPage(ctx web.Context) error {
	return ctx.Render("pages/register")
}

func ForgotPage(ctx web.Context) error {
	return ctx.Render("pages/forgot")
}

func DashboardPage(ctx web.Context) error {
	fmt.Println(len(ctx.Request().Cookies()))
	_, err := ctx.Request().Cookie("SESSION-ID")
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, http.ErrNoCookie) {
			return ctx.Redirect("/login")
		}
		log.Printf("Uncommon error trying to read cookie: %s", err)
	}
	return ctx.Render("pages/dashboard")
}

func NewEventPage(ctx web.Context) error {
	fmt.Println(len(ctx.Request().Cookies()))
	_, err := ctx.Request().Cookie("SESSION-ID")
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, http.ErrNoCookie) {
			return ctx.Redirect("/login")
		}
		log.Printf("Uncommon error trying to read cookie: %s", err)
	}
	return ctx.Render("pages/new_event")
}
