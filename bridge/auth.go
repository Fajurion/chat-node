package bridge

type ConnectionToken struct {
	UserID  int64
	Session string
}

var ConnectionTokens = make(map[string]ConnectionToken)

func CheckToken(token string) int64 {
	return ConnectionTokens[token].UserID
}

func AddToken(token string, id int64, session string) {
	ConnectionTokens[token] = ConnectionToken{
		UserID:  id,
		Session: session,
	}
}
