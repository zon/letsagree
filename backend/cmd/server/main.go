package main

import (
	"context"
	"log"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server/internal/auth"
	"server/internal/oidc"
	"server/internal/store"
)

var version = "dev"

var cli struct {
	Addr      string `name:"addr" default:":8080" help:"HTTP listen address."`
	ConfigDir string `name:"config" default:"config" help:"Path to config directory."`
}

func RegisterRoutes(r *gin.Engine, o *auth.Orchestration) {
	authGroup := r.Group("/auth")
	authGroup.GET("/login", o.Login)
	authGroup.GET("/callback", o.Callback)
	authGroup.POST("/logout", o.Logout)

	protected := r.Group("")
	protected.Use(o.RequireAuth)
	protected.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

func main() {
	kong.Parse(&cli)

	cfg, err := oidc.LoadConfig(cli.ConfigDir + "/humanity-protocol.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	var provider auth.OIDCProvider
	if cfg.IssuerURL != "" {
		realProvider, err := oidc.NewProvider(ctx, cfg)
		if err != nil {
			log.Fatalf("failed to create OIDC provider: %v", err)
		}
		provider = realProvider
	} else {
		provider = oidc.NewStubProvider(oidc.AnyIDToken())
	}

	o := auth.NewOrchestration(provider, store.StubSessions(), store.StubUsers())

	var db *gorm.DB
	db, err = store.NewDB(cli.ConfigDir + "/postgres.json")
	if err != nil {
		log.Printf("warning: could not connect to database, using stub stores: %v", err)
	} else {
		if err := db.AutoMigrate(&store.User{}, &store.Session{}); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
		o = auth.NewOrchestration(provider, store.New(db), store.New(db))
	}

	r := gin.Default()
	RegisterRoutes(r, o)

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, version)
	})

	log.Printf("starting server on %s", cli.Addr)
	if err := r.Run(cli.Addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}