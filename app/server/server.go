package server

import (
	appWebsocker "alekseikormski.com/server-status-monitoring/app/web-socket"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

type Server struct {
	port int
	wss  *appWebsocker.WebSocketServer
	res  embed.FS
}

func NewServer(port int, wss *appWebsocker.WebSocketServer, res embed.FS) *Server {
	return &Server{
		port: port,
		wss:  wss,
		res:  res,
	}
}

func (s *Server) RegisterRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", Application(s.res))
	r.Static("/static", filepath.Join("front-end", "build", "static"))

	r.GET("/healthz", Healthz)
	r.GET("/websocket", WebSocket(s.wss))
	return r
}

func (s *Server) Start(engine *gin.Engine) error {
	if err := engine.Run(fmt.Sprintf(":%d", s.port)); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
