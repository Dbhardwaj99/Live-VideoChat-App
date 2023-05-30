package handlers

import (
	"fmt"
	"os"
	w "videochat-app/pkg/webrtc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	guuid "github.com/google/uuid"
)

func RoomCreate(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("/room/%s", guuid.New().String()))
}

func Room(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		c.Status(404)
		return c.Redirect("/")
	}

	ws := "ws"
	if os.Getenv("APP_ENV") == "production" {
		ws = "wss"
	}

	uuid, suuid, _ := createOrGetRoom(uuid)
	return c.Render("peer", fiber.Map{
		"RoomWSAddr":   fmt.Sprintf("%s://%s/room/%s/ws", ws, c.Hostname(), uuid),
		"RoomLink":     fmt.Sprintf("%s://%s/room/%s", c.Protocol(), c.Hostname(), uuid),
		"ChatWSAddr":   fmt.Sprintf("%s://%s/room/%s/chat/ws", ws, c.Hostname(), uuid),
		"StreamLink":   fmt.Sprintf("%s://%s/stream/%s", c.Protocol(), c.Hostname(), suuid),
		"ViewerWSAddr": fmt.Sprintf("%s://%s/room/%s/viewer/ws", ws, c.Hostname(), uuid),
		"Type":         "room",
	}, "layouts/main")
}

func RoomWs(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	_, _, room := createOrGetRoom(uuid)
	w.RoomConn(c, room.Peers)
}

func createOrGetRoom(uuid string) (string, string, *w.Room) {
	room, ok := rooms[uuid]
	if !ok {
		room = NewRoom(uuid)
		rooms[uuid] = room
	}

	return uuid, room.Suuid, room
}

func RoomViewerWs(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	w.RoomLock.Lock()
	if peer, ok := w.Rooms[uuid]; ok {
		w.RoomLock.Unlock()
		roomViewerConn(c, peer.Peers)
		return
	}
	w.RoomLock.Unlock()

}

func roomViewerConn(c *websocket.Conn, p *w.Peers) {

}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
