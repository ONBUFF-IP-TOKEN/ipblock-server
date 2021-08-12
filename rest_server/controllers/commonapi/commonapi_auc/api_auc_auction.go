package commonapi_auc

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func PostAucAuctionRegister(auction *context_auc.AucAuctionRegister, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_product table에 저장
	auction.CreateTs = datetime.GetTS2MilliSec()
	if id, err := model.GetDB().InsertAucAuction(auction); err != nil {
		log.Error("InsertProduct :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		auction.Id = id
		resp.Value = auction
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostAucAuctionUpdate(auction *context_auc.AucAuctionUpdate, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_product table에 업데이트
	if id, err := model.GetDB().UpdateAucAuction(auction); err != nil {
		log.Error("UpdateAucAuction :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		auction.Id = id
		resp.Value = auction
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetAucAuctionList(auctionList *context_auc.AuctionList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)

	//redis exist check
	if pageInfo, auctions, err := model.GetDB().GetAuctionListCache(&auctionList.PageInfo); err == nil {
		resp.Success()
		resp.Value = context_auc.AuctionListResponse{
			PageInfo:    *pageInfo,
			AucAuctions: *auctions,
		}
	} else {
		// cache 에 없다면 db에서 직접 로드
		auctions, totalCount, err := model.GetDB().GetAucAuctionList(auctionList)
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.Success()
			pageInfo := context_auc.PageInfoResponse{
				PageOffset: auctionList.PageOffset,
				PageSize:   int64(len(auctions)),
				TotalSize:  totalCount,
			}
			resp.Value = context_auc.AuctionListResponse{
				PageInfo:    pageInfo,
				AucAuctions: auctions,
			}
			model.GetDB().SetAuctionListCache(&auctionList.PageInfo, &pageInfo, &auctions)
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func DeleteAucAuctiontRemove(product *context_auc.RemoveAuction, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_products table 에서 삭제
	if ret, err := model.GetDB().DeleteAucAuction(product.Id); err != nil {
		log.Error("DeleteAucAuctiontRemove :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if !ret {
			resp.SetReturn(resultcode.Result_DBNotExistProduct)
		}
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
