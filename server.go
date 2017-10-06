package main

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olebedev/staticbin"
)

// Server route handlers
type Server struct {
	ginEngine *gin.Engine
	port      int
	debug     bool
	store     *Store
}

// NewServer creates a new server
func NewServer(debug bool, port int, store *Store) *Server {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	s := &Server{
		ginEngine: r,
		port:      port,
		debug:     debug,
		store:     store,
	}

	if debug {
		r.Static("/static", "./static")
	} else {
		r.Use(staticbin.Static(Asset, staticbin.Options{
			SkipLogging: true,
		}))
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://pagaraprata.netlify.com"},
		AllowMethods:     []string{"GET", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", s.index)
	r.GET("/version", s.version)
	r.POST("/api/debts", s.createDebt)
	r.GET("/api/debts/:id", s.getDebt)
	r.PUT("/api/debts/:id/:secret", s.updateDebt)

	return s
}

// Just for testing
func (s *Server) routes() *gin.Engine {
	return s.ginEngine
}

// Run the server
func (s *Server) Run() {
	s.ginEngine.Run("0.0.0.0:" + strconv.Itoa(s.port))
}

func (s *Server) loadFile(filePath string) ([]byte, error) {
	if s.debug {
		return ioutil.ReadFile(filePath)
	}

	return Asset(filePath)
}
