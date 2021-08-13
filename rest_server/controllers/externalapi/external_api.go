package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi/commonapi_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
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

	if conf.Auth.AuthEnable {
		author, ok := c.Request().Header["Authorization"]
		if !ok {
			// auth token 오류 리턴
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		if !conf.Auth.InternalAuth {
			walletAddr, isValid := auth.GetIAuth().IsValidAuthToken(author[0][7:])
			if !isValid {
				// auth token 오류 리턴
				res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

				return base.PreCheckResponse{
					IsSucceed: false,
					Response:  res,
				}
			}
			base.GetContext(c).(*context.IPBlockServerContext).SetWalletAddr(*walletAddr)
		} else {
			// membership server 인증 진행
			walletAddr, _, isValid := auth.GetIAuth().GetAuthInfo(author[0][7:])
			if !isValid {
				// auth token 오류 리턴
				res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

				return base.PreCheckResponse{
					IsSucceed: false,
					Response:  res,
				}
			}

			if ret, err := auth.CheckAuthToken(walletAddr, author[0][7:]); err != nil || !ret {
				res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)
				return base.PreCheckResponse{
					IsSucceed: false,
					Response:  res,
				}
			}

			base.GetContext(c).(*context.IPBlockServerContext).SetWalletAddr(walletAddr)
		}
	} else {
		//base.GetContext(c).(*context.IPBlockServerContext).SetWalletAddr("0x9Ec7EDE9204E17dfa34e1d381ED5f49A0D578e96")
		base.GetContext(c).(*context.IPBlockServerContext).SetWalletAddr("0x38f998d033990a315b08afc0f78059fb7d11dc4d")
	}

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

// product apis(v1.1)
func (o *ExternalAPI) GetProductList(c echo.Context) error {
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

func (o *ExternalAPI) PostProductOrder(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewOrderProduct()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostProductOrder(params, ctx)
}

func (o *ExternalAPI) GetMyOrderList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context.NewOrderList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.GetMyOrderList(params, ctx)
}

// 경매 리스트 요청
func (o *ExternalAPI) GetAucAuctionList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewAuctionList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucAuctionList(params, ctx)
}

// 경매 입찰 보증금 확인
func (o *ExternalAPI) GetAucBidDeposit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidDepositVerify()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucBidDeposit(params, ctx)
}

// 경매 입찰 보증금정보 전송
func (o *ExternalAPI) PostAucBidDeposit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidDeposit()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostAucBidDeposit(params, ctx)
}

// 경매 입찰 진행
func (o *ExternalAPI) PostAucBidSubmit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidSubmit()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostAucBidSubmit(params, ctx)
}

// 경매 입찰 리스트 요청
func (o *ExternalAPI) GetAucBidList(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidAttendeeList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetAucBidList(params, ctx)
}

// 낙찰 받기
func (o *ExternalAPI) PostAucBidWinnerSubmit(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidSuccess()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostAucBidWinnerSubmit(params, ctx)
}

// 낙찰 포기
func (o *ExternalAPI) PostAucBidWinnerGiveUp(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewBidSuccessGiveup()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.NewBidSuccessGiveup(params, ctx)
}
