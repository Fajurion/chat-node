package bridge

type ConnectionToken struct {
	UserID  int64
	Session uint64
}

var ConnectionTokens = make(map[string]ConnectionToken)

func CheckToken(token string) ConnectionToken {
	return ConnectionTokens[token]
}

func RemoveToken(token string) {
	delete(ConnectionTokens, token)
}

func AddToken(token string, id int64, session uint64) {
	ConnectionTokens[token] = ConnectionToken{
		UserID:  id,
		Session: session,
	}
}
