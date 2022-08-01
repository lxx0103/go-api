package main

import (
	"go-api/cmd"
	"os"
)

// @title go-api API
// @version 1.0
// @description API for go-api.

// @contact.name Lewis
// @contact.email lxx0103@yahoo.com

// @host 0.0.0.0:8080
// @BasePath /
func main() {
	cmd.Run(os.Args)
}
