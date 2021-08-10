package internalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi/commonapi_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
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

// product apis(m1.1)
// product 등록
func (o *InternalAPI) PostRegisterProduct(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewProductInfo()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostRegisterProduct(params, ctx)
}

// product 삭제
func (o *InternalAPI) DeleteUnregisterProduct(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewUnregisterProduct()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.DeleteUnregisterProduct(params, ctx)
}

// product 업데이트
func (o *InternalAPI) PostUpdateProduct(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewProductInfo()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostUpdateProduct(params, ctx)
}

// product state만 업데이트
func (o *InternalAPI) PostUpdateProductState(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewProductUpdateState()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostUpdateProductState(params, ctx)
}

// product list 조회
func (o *InternalAPI) GetProductList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewProductList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetProductList(params, ctx)
}

// nft list 조회
func (o *InternalAPI) GetNftList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewNftList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetNftList(params, ctx)
}

// auc 물품 등록
func (o *InternalAPI) PostAucProductRegister(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewProductInfo()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostAucProductRegister(params, ctx)
}

// auc 물품 삭제
func (o *InternalAPI) DeleteAucProductRemove(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewRemoveProduct()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.DeleteAucProductRemove(params, ctx)
}

func (o *InternalAPI) GetAucProductList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewProductList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucProductList(params, ctx)
}
