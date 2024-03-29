package util

import (
	"errors"

	integration "fajurion.com/node-integration"
)

type AppToken struct {
	Node   uint   `json:"node"` // Node ID
	Domain string `json:"domain"`
	Token  string `json:"token"`
}

func ConnectToApp(account string, session string, app uint, cluster uint) (AppToken, error) {

	res, err := integration.PostRequest("/node/get_lowest", map[string]interface{}{
		"account": account,
		"session": session,
		"app":     app,
		"cluster": cluster,
		"node":    integration.NODE_ID,
		"token":   integration.NODE_TOKEN,
	})
	if err != nil {
		return AppToken{}, err
	}

	if !res["success"].(bool) {
		return AppToken{}, errors.New(res["error"].(string))
	}

	return AppToken{
		Node:   uint(res["id"].(float64)),
		Domain: res["domain"].(string),
		Token:  res["token"].(string),
	}, nil
}
