package server

import (
	appWebsocker "alekseikormski.com/server-status-monitoring/app/web-socket"
	"embed"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Healthz(c *gin.Context) {
	c.Status(http.StatusOK)
}

func Application(res embed.FS) func(c *gin.Context) {
	return func(c *gin.Context) {
		content, err := res.ReadFile("front-end/build/index.html")
		if err != nil {
			log.Printf("cannot return index.html: %v", err)
			c.Status(500)
			return
		}

		c.Writer.Write(content)
	}
}

func WebSocket(wss *appWebsocker.WebSocketServer) func(c *gin.Context) {
	return func(c *gin.Context) {
		wss.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		// upgrade this connection to a WebSocketServer
		// connection
		ws, err := wss.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
		}

		defer ws.Close()
		go wss.Handle(ws)
		wss.Reader(ws)
	}
}
