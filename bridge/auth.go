package bridge

type ConnectionToken struct {
	UserID   string
	Session  string
	Username string
	Tag      string
}

var ConnectionTokens = make(map[string]ConnectionToken)

func CheckToken(token string) ConnectionToken {
	return ConnectionTokens[token]
}

func RemoveToken(token string) {
	delete(ConnectionTokens, token)
}

func AddToken(token string, id string, session string, username string, tag string) {
	ConnectionTokens[token] = ConnectionToken{
		UserID:   id,
		Session:  session,
		Username: username,
		Tag:      tag,
	}
}
