package message_routes

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"

	integration "fajurion.com/node-integration"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Request for deleting a message
type deleteMessageRequest struct {
	TokenID     string `json:"id"`          // Conversation token id
	Token       string `json:"token"`       // Conversation token (token)
	Certificate string `json:"certificate"` // Message certificate
}

// Route: /conversations/message/delete
func deleteMessage(c *fiber.Ctx) error {

	// Parse request
	var req deleteMessageRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get conversation token
	token, err := caching.ValidateToken(req.TokenID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "invalid conversation token")
	}

	// Get claims from message certificate
	claims, valid := conversations.GetCertificateClaims(req.Certificate)
	if !valid {
		return integration.InvalidRequest(c, "invalid certificate claims")
	}

	// Check if certificate is valid for the provided conversation token
	if !claims.Valid(claims.Message, token.Conversation, token.ID) {
		return integration.InvalidRequest(c, "no permssion to delete message")
	}

	// Delete the message in the database
	if err := database.DBConn.Where("id = ?", claims.Message).Delete(&conversations.Message{}).Error; err != nil && err != gorm.ErrRecordNotFound {
		return integration.FailedRequest(c, "server.error", err)
	}

	// Send a system message to delete the message on all clients who are storing it
	if err := SendSystemMessage(claims.Conversation, DeletedMessage, []string{claims.Message}); err != nil {
		return integration.FailedRequest(c, "server.error", err)
	}

	return integration.SuccessfulRequest(c)
}
