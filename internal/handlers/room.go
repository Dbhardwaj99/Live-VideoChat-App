package handlers

import (
	"fmt"
	"os"
	"time"

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

	uuid, suuid,_ := createOrGetRoom(uuid)
	return c.Render("room", fiber.Map{
		"uuid": uuid,
		"room": room,
	})
}

func RoomWs(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	_,_, room := createOrGetRoom(uuid)
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

	_,_, room := createOrGetRoom(uuid)
}

func roomViewerConn(c *websocket.Conn, p *w.Peers) {

}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}	