package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func GetMyOrderList(reqInfo *context.OrderList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	orders, totalCount, err := model.GetDB().GetMyOrderList(reqInfo)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context.PageInfoResponse{
			PageOffset: reqInfo.PageOffset,
			PageSize:   int64(len(orders)),
			TotalSize:  totalCount,
		}
		resp.Value = context.OrderListResponse{
			PageInfo: pageInfo,
			Orders:   orders,
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
