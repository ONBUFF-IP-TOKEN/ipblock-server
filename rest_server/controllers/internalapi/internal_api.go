package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/labstack/echo"
)

type InternalAPI struct {
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

	// ctx := base.GetContext(c).(*context.LinkViewContext)
	// cerberus := libcerberus.GetInstance()
	// resp, err := cerberus.AuthInternal(ctx.ApplicationKey(), ctx.ApplicationSecret())
	// if err != nil {
	// 	log.Error(err)
	// 	return base.PreCheckResponse{
	// 		IsSucceed: false,
	// 		Response:  base.ResponseInternalServerError(),
	// 	}
	// }
	// if resp.Result != "000" {
	// 	return base.PreCheckResponse{
	// 		IsSucceed: false,
	// 		Response:  resp,
	// 	}
	// }

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

func (o *InternalAPI) Init(e *echo.Echo) error {
	o.echo = e
	o.BaseController.PreCheck = PreCheck

	if err := o.MapRoutes(o, e, o.apiConf.Routes); err != nil {
		return err
	}

	// // serving documents for RESTful APIs
	// if o.conf.LinkView.APIDocs {
	// 	e.Static("/docs", "docs/int")
	// }

	return nil
}

func (o *InternalAPI) GetConfig() *baseconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[0]
	return o.apiConf
}

func NewAPI() *InternalAPI {
	return &InternalAPI{}
}

func (o *InternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

func (o *InternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}
