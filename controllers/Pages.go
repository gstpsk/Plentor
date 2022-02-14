package controllers

import (
	"github.com/go-zepto/zepto/web"
)

func IndexPage(ctx web.Context) error {
	return ctx.Render("pages/index")
}

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
	if !RequestIsAuthorized(ctx) {
		ctx.Redirect("/login")
	}
	return ctx.Render("pages/dashboard")
}

func NewEventPage(ctx web.Context) error {
	if !RequestIsAuthorized(ctx) {
		ctx.Redirect("/login")
	}
	return ctx.Render("pages/new_event")
}

func EventPage(ctx web.Context) error {
	if !RequestIsAuthorized(ctx) {
		ctx.Redirect("/login")
	}
	return ctx.Render("pages/event")
}

func EventSignupPage(ctx web.Context) error {
	return ctx.Render("pages/event_signup")
}
