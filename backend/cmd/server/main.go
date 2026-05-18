package main

import (
	"log"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/gin-gonic/gin"
)

var version = "dev"

var cli struct {
	Addr      string `name:"addr" default:":8080" help:"HTTP listen address."`
	ConfigDir string `name:"config" default:"config" help:"Path to config directory."`
}

func main() {
	kong.Parse(&cli)

	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, version)
	})

	log.Printf("starting server on %s", cli.Addr)
	if err := r.Run(cli.Addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
