package controllers

import "github.com/go-zepto/zepto/web"

func LoginPage(ctx web.Context) error {
	return ctx.Render("pages/login")
}
