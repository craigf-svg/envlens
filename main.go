package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No local .env found")
	}
	envList := os.Environ()
	fmt.Println("Environment Variables Found: ", envList)
}
