package commonapi_auc

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func DeleteSystemRedisRemove(systemRedis *context_auc.SystemRedisRemove, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if strings.EqualFold(systemRedis.AuctionList, context_auc.TRUE) {
		model.GetDB().DeleteAuctionList()
		model.GetDB().DeleteAuctionCacheAll()
	}
	if strings.EqualFold(systemRedis.ProductList, context_auc.TRUE) {
		model.GetDB().DeleteProductList()
	}
	if len(systemRedis.BidList) > 0 {
		aucId, _ := strconv.ParseInt(systemRedis.BidList, 10, 64)
		model.GetDB().DeleteBidList(aucId)
	}
	if len(systemRedis.AuctionId) > 0 {
		aucId, _ := strconv.ParseInt(systemRedis.AuctionId, 10, 64)
		model.GetDB().CacheDelProduct(aucId)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
