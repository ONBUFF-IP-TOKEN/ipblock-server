package internalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi/commonapi_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/labstack/echo"
)

// 공지 사항 등록
func (o *InternalAPI) PostNoticeRegister(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewNoticeRegister()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostNoticeRegister(params, ctx)
}

// 공지 사항 리스트 요청
func (o *InternalAPI) GetNotice(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewNoticeList()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.GetNotice(params, ctx.EchoContext)
}

// 공지 사항 삭제
func (o *InternalAPI) DeleteNoticeRemove(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewNoticeRemove()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.DeleteNoticeRemove(params, ctx)
}

// 공지 사항 수정
func (o *InternalAPI) PostNoticeUpdate(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)

	params := context_auc.NewNoticeUpdate()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi_auc.PostNoticeUpdate(params, ctx)
}
