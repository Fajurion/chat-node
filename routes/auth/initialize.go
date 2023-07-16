package auth

import (
	"chat-node/database"
	"chat-node/database/fetching"
	"chat-node/util"
	"chat-node/util/requests"
	"log"

	"github.com/Fajurion/pipesfiber"
	"github.com/gofiber/fiber/v2"
)

type intializeRequest struct {
	NodeToken  string   `json:"node_token"`
	Session    string   `json:"session"`
	UserID     string   `json:"user_id"`
	Username   string   `json:"username"`
	Tag        string   `json:"tag"`
	SessionIds []string `json:"session_ids"`
}

func initializeConnection(c *fiber.Ctx) error {

	// Parse the request
	var req intializeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if util.NODE_TOKEN != req.NodeToken {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	log.Println(req.UserID, "|", req.SessionIds)
	database.DBConn.Where("account = ?", req.UserID).Not("id IN ?", req.SessionIds).Delete(&fetching.Session{})

	tk := util.GenerateToken(200)

	// Check if there are too many users
	if pipesfiber.GetConnections(req.UserID) >= 3 {
		return requests.FailedRequest(c, "too.many.connections", nil)
	}

	pipesfiber.AddToken(tk, pipesfiber.ConnectionToken{
		UserID:  req.UserID,
		Session: req.Session,
		Data: util.UserData{
			Username: req.Username,
			Tag:      req.Tag,
		},
	})

	return c.JSON(fiber.Map{
		"success": true,
		"load":    0,
		"token":   tk,
	})
}
