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
	"github.com/labstack/echo"
)

// 경매 등록
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

//경매 정보 업데이트
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

// 경매 정보 리스트 요청
func GetAucAuctionList(auctionList *context_auc.AuctionList, c echo.Context) error {
	resp := new(base.BaseResponse)

	// active 경매 정보만 redis에서 가져온다.
	if auctionList.ActiveState == context_auc.Auction_active_state_active {
		//redis exist check
		if pageInfo, auctions, err := model.GetDB().GetAuctionListCache(&auctionList.PageInfo); err == nil {
			resp.Success()
			resp.Value = context_auc.AuctionListResponse{
				PageInfo:    *pageInfo,
				AucAuctions: *auctions,
			}
			return c.JSON(http.StatusOK, resp)
		}
	}

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

		// active 경매 정보만 redis에 남긴다.
		if auctionList.ActiveState == context_auc.Auction_active_state_active {
			model.GetDB().SetAuctionListCache(&auctionList.PageInfo, &pageInfo, &auctions)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// 경매 삭제
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
