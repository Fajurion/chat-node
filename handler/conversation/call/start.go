package call

import (
	"chat-node/caching"
	"log"

	"github.com/Fajurion/pipesfiber/wshandler"
)

func start(message wshandler.Message) {

	if message.ValidateForm("id", "token") {
		wshandler.ErrorResponse(message, "invalid")
		log.Println("invalid form")
		return
	}

	_, err := caching.ValidateToken(message.Data["id"].(string), message.Data["token"].(string))
	if err != nil {
		wshandler.ErrorResponse(message, "invalid")
		println("invalid token")
		return
	}

	// TODO: Finish

}
