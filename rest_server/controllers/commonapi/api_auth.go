package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/constant"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
	"github.com/labstack/echo"
)

func PostLogin(c echo.Context) error {
	params := context.NewLoginParam()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(context.IpBlockBaseResponse)
	// 1. verify sign check
	if !token.GetToken().VerifySign(params.WalletAuth.WalletAddr, params.WalletAuth.Message, params.WalletAuth.Sign) {
		// invalid sign info
		resp.SetResult(constant.Result_Auth_InvalidLoginInfo)
		return c.JSON(http.StatusOK, resp)
	}

	// 2. redis duplicate check
	if authInfo, err := model.GetDB().GetAuthInfo(params.WalletAuth.WalletAddr); err == nil {
		// redis에 기존 정보가 있다면 기존에 발급된 토큰으로 응답한다.
		resp.Success()
		resp.Value = context.LoginResponse{
			AuthToken:  authInfo.AuthToken,
			ExpireDate: authInfo.ExpireDate,
		}
	} else {
		// 3. create auth token
		authToken, expireDate, err := auth.GetIAuth().EncryptJwt(params.WalletAuth.WalletAddr)
		if err != nil {
			resp.SetResult(constant.Result_Auth_DontEncryptJwt)
		} else {
			resp.Success()
			resp.Value = context.LoginResponse{
				AuthToken:  authToken,
				ExpireDate: expireDate,
			}

			// 3. redis save
			authInfo := model.AuthInfo{
				AuthToken:  authToken,
				ExpireDate: expireDate,
				WalletAuth: params.WalletAuth,
			}
			err = model.GetDB().SetAuthInfo(&authInfo)
			if err != nil {
				return base.BaseJSONInternalServerError(c, err)
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}
