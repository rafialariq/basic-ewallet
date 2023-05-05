package main

import (
	"final_project_easycash/delivery"

	_ "github.com/lib/pq"
)

func main() {
	delivery.Server().Run()
}
