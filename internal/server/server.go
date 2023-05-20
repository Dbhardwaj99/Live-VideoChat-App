package server

import (
	"flag"
	"os"
	"time"

	// import handlers "internal/handlers"
	"videochat-app/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
)

var (
	addr = flag.String("addr, ':'", os.Getenv("PORT"), "")
	cert = flag.String("cert", "", "")
	key  = flag.String("key", "", "")
)

func Run() error {
	flag.Parse()

	if *addr == ":" {
		*addr = ":8080"
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes and handlers
	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.RoomCreate)
	app.Get("/room/:uuid", handlers.Room)

	app.Get("/room/:uuid/ws", websocket.New(handlers.RoomWs, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))

	app.Get("/room/:uuid/chat", handlers.RoomChat)
	app.Get("/room/:uuid/chat/ws", websocket.New(handlers.RoomChatWs))
	app.Get("/room/:uuid/viewer/ws", websocket.New(handlers.RoomViewerWs))

	app.Get("/stream/:ssuid", handlers.Stream)

	app.Get("/stream/:ssuid/ws", websocket.New(handlers.StreamWs))
	app.Get("/stream/:ssuid/chat/ws", handlers.StreamChatWs)
	app.Get("/stream/:ssuid/viewer/ws", websocket.New(handlers.StreamViewerWs))

}
