package calls

import (
	"chat-node/database/credentials"
	"time"

	"github.com/livekit/protocol/auth"
)

func GetJoinToken(room, identity string) (string, error) {

	// Generate token
	at := auth.NewAccessToken(credentials.LIVEKIT_KEY, credentials.LIVEKIT_SECRET)

	// Add permissions
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
