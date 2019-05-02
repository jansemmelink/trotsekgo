package main

import (
	"net/http"

	"github.com/jansemmelink/log"
	"github.com/jansemmelink/trotsek/lib/auth"
	"github.com/jansemmelink/trotsek/lib/router"
)

func main() {
	log.DebugOn()
	a := auth.New()
	r := router.New()
	r.Sub("/auth").
		Sub("register").Post(a.RegisterOper()).
		Sub("reset").Post(a.ResetOper()).
		Sub("activate").Post(a.ActivateOper()).
		Sub("login").Post(a.LoginOper()).
		Sub("logout").Post(a.LogoutOper())

	//authApp.r.Post("/auth/register", authApp.a.Register)
	err := http.ListenAndServe("localhost:8070", r)
	if err != nil {
		panic(err)
	}
}
