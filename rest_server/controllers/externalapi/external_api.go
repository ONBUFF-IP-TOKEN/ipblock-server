package externalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/labstack/echo"
)

type ExternalAPI struct {
	base.BaseController

	conf    *config.ServerConfig
	apiConf *baseconf.APIServer
	echo    *echo.Echo
}

func PreCheck(c echo.Context) base.PreCheckResponse {
	conf := config.GetInstance()
	if err := base.SetContext(c, &conf.Config, context.NewIPBlockServerContext); err != nil {
		log.Error(err)
		return base.PreCheckResponse{
			IsSucceed: false,
		}
	}

	// auth token 검증
	walletAddr, isValid := auth.GetIAuth().IsValidAuthToken(c.Request().Header["Authorization"][0][7:])
	if conf.Auth.AuthEnable && !isValid {
		// auth token 오류 리턴
		res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

		return base.PreCheckResponse{
			IsSucceed: false,
			Response:  res,
		}
	}
	base.GetContext(c).(*context.IPBlockServerContext).SetWalletAddr(*walletAddr)

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

func (o *ExternalAPI) Init(e *echo.Echo) error {
	o.echo = e
	o.BaseController.PreCheck = PreCheck

	if err := o.MapRoutes(o, e, o.apiConf.Routes); err != nil {
		return err
	}

	// serving documents for RESTful APIs
	if o.conf.IPServer.APIDocs {
		e.Static("/docs", "docs/ext")
	}

	return nil
}

func (o *ExternalAPI) GetConfig() *baseconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[1]
	return o.apiConf
}

func NewAPI() *ExternalAPI {
	return &ExternalAPI{}
}

func (o *ExternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

func (o *ExternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}

func (o *ExternalAPI) PostLogin(c echo.Context) error {
	return commonapi.PostLogin(c)
}

func (o *ExternalAPI) PostRegisterItem(c echo.Context) error {
	return commonapi.PostRegisterItem(c)
}

func (o *ExternalAPI) DeleteUnregisterItem(c echo.Context) error {
	return commonapi.DeleteUnregisterItem(c)
}

func (o *ExternalAPI) GetItemList(c echo.Context) error {
	return commonapi.GetItemList(c)
}

func (o *ExternalAPI) PostPurchaseItem(c echo.Context) error {
	return commonapi.PostPurchaseItem(c)
}

func (o *ExternalAPI) GetHistoryTransferItem(c echo.Context) error {
	return commonapi.GetHistoryTransferItem(c)
}

func (o *ExternalAPI) GetHistoryTransferMe(c echo.Context) error {
	return commonapi.GetHistoryTransferMe(c)
}
