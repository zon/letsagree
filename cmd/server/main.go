//go:build !test

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var version = "dev"

func flagAddr() string {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		return os.Getenv("TEST_ADDR")
	}
	addr := ":8080"
	for i := 1; i < len(os.Args)-1; i++ {
		if os.Args[i] == "--addr" {
			addr = os.Args[i+1]
			break
		}
	}
	return addr
}

func main() {
	addr := flagAddr()

	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}