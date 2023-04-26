package main

import (
	"final_project_easycash/delivery"

	_ "github.com/lib/pq"
)

func main() {

	// f, _ := os.Create(".log")
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	delivery.Server().Run()
}
