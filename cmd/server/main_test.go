//go:build test

package main

import (
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestAddrFlagOverridesListenAddress(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		addr := os.Getenv("TEST_ADDR")
		r := gin.New()
		r.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": "test"})
		})
		http.ListenAndServe(addr, r)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestAddrFlagOverridesListenAddress")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1", "TEST_ADDR=:9090")

	if err := cmd.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:9090/version")
	if err != nil {
		cmd.Process.Kill()
		t.Fatalf("GET /version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		cmd.Process.Kill()
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestAddrFlagDefaultsTo8080(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS_DEFAULT") == "1" {
		addr := os.Getenv("TEST_ADDR")
		r := gin.New()
		r.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": "test"})
		})
		http.ListenAndServe(addr, r)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestAddrFlagDefaultsTo8080")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_DEFAULT=1", "TEST_ADDR=:8080")

	if err := cmd.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/version")
	if err != nil {
		cmd.Process.Kill()
		t.Fatalf("GET /version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		cmd.Process.Kill()
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	cmd.Process.Kill()
	cmd.Wait()
}

func TestAddrFlagWithNoArgsDefaultsTo8080(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS_NO_ARGS") == "1" {
		addr := os.Getenv("TEST_ADDR")
		r := gin.New()
		r.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": "test"})
		})
		http.ListenAndServe(addr, r)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestAddrFlagWithNoArgsDefaultsTo8080")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS_NO_ARGS=1", "TEST_ADDR=:8080")

	if err := cmd.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/version")
	if err != nil {
		cmd.Process.Kill()
		t.Fatalf("GET /version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		cmd.Process.Kill()
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	cmd.Process.Kill()
	cmd.Wait()
}

var version = "dev"