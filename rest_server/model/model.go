package model

import "github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"

type AuthInfo struct {
	AuthToken  string            `json:"auth_token"`
	ExpireDate int64             `json:"expire_date"`
	WalletAuth context.LoginAuth `json:"wallet_auth"`
}
