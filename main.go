package main

import (
	"assignment2/api/controllers"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Starting assignment2")
	fmt.Println("running server")

	app := controllers.App{}
	app.Initialize()
	app.RunServer()

}
