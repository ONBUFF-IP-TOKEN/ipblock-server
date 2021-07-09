package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/labstack/echo"
)

func GetHistoryTransferItem(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewGetHistoryTransferItem()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(base.BaseResponse)
	historys, totalCount, err := model.GetDB().GetHistoryTransferItem(params)

	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context.PageInfoResponse{
			PageOffset: params.PageOffset,
			PageSize:   int64(len(historys)),
			TotalSize:  totalCount,
		}
		resp.Value = context.GetHistoryTransferItemResponse{
			PageInfo: pageInfo,
			Historys: historys,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetHistoryTransferMe(c echo.Context) error {
	ctx := base.GetContext(c).(*context.IPBlockServerContext)
	params := context.NewGetHistoryTransferMe()
	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	if err := params.CheckValidate(ctx); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	resp := new(base.BaseResponse)
	historys, totalCount, err := model.GetDB().GetHistoryTransferMe(params)

	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context.PageInfoResponse{
			PageOffset: params.PageOffset,
			PageSize:   int64(len(historys)),
			TotalSize:  totalCount,
		}
		resp.Value = context.GetHistoryTransferItemResponse{
			PageInfo: pageInfo,
			Historys: historys,
		}
	}

	return c.JSON(http.StatusOK, resp)
}
