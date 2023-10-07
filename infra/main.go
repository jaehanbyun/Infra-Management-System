package main

import (
	"net/http"

	"github.com/jaehanbyun/infra/app"
)

func main() {
	a := app.MakeHandler()

	http.ListenAndServe(":5000", a)
}
