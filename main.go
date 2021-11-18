package main

import (
	"os"
)

func main() {
	domain := os.Getenv("quodd_domain")
	username := os.Getenv("quodd_username")
	password := os.Getenv("quodd_password")
	gateway := os.Getenv("quodd_gateway")

	if domain == "" || username == "" || password == "" || gateway == "" {
		panic("Missing environment variables")
	}

}
