package main

import "log"

func main() {

	r := InitWebServer()
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
