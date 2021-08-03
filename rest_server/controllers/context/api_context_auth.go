package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

///////////// API ///////////////////////////////
// login
type LoginAuth struct {
	WalletAddr string `json:"wallet_address" validate:"required"`
	Message    string `json:"message" validate:"required"`
	Sign       string `json:"sign" validate:"required"`
}
type LoginParam struct {
	WalletType string    `json:"wallet_type" validate:"required"`
	WalletAuth LoginAuth `json:"wallet_auth" validate:"required"`
}

func NewLoginParam() *LoginParam {
	return new(LoginParam)
}

func (o *LoginParam) CheckValidate() *base.BaseResponse {
	if len(o.WalletType) == 0 && (Wallet_type_metamask != o.WalletType) {
		return base.MakeBaseResponse(resultcode.Result_Auth_InvalidWalletType)
	}
	if len(o.WalletAuth.WalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireWalletAddress)
	}
	if len(o.WalletAuth.Message) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_RequireMessage)
	}
	if len(o.WalletAuth.Sign) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_RequireSign)
	}
	return nil
}

type LoginResponse struct {
	AuthToken  string `json:"auth_token" validate:"required"`
	ExpireDate int64  `json:"expire_date" validate:"required"`
}

/////////////////////////

// verify auth token
type VerifyAuthToken struct {
	WalletAddr string `json:"wallet_address" validate:"required"`
	AuthToken  string `json:"auth_token" validate:"required"`
}

func NewVerifyAuthToken() *VerifyAuthToken {
	return new(VerifyAuthToken)
}

/////////////////////////
