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
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/schedule"
	"github.com/labstack/echo"
)

// 경매 등록
func PostAucAuctionRegister(auction *context_auc.AucAuctionRegister, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. 해당 product이 존재하는지 체크
	product, err := model.GetDB().GetAucProductById(auction.ProductId)
	if err != nil {
		log.Error("GetAucProductById :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if product == nil {
			log.Error("GetAucProductById invalid product id ")
			resp.SetReturn(resultcode.Result_Auc_Auction_RequireProductId)
		} else {
			//2. 상품 가격 정보 복사
			auction.BidStartAmount = product.Prices[0].Price
			auction.BidCurAmount = 0
			auction.BidDeposit = context_auc.CheckDepositPrice(product.Prices[0].Price)
			auction.TokenType = product.Prices[0].TokenType
			auction.Price = product.Prices[0].Price

			//3. auc_product table에 저장
			auction.CreateTs = datetime.GetTS2MilliSec()
			if id, err := model.GetDB().InsertAucAuction(auction); err != nil {
				log.Error("InsertProduct :", err)
				resp.SetReturn(resultcode.Result_DBError)
			} else {
				auction.Id = id
				resp.Value = auction

				// 스케줄러 리셋
				schedule.GetScheduler().ResetAuctionScheduler()
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

//경매 정보 업데이트
func PostAucAuctionUpdate(auction *context_auc.AucAuctionUpdate, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 금액 check
	auction.BidStartAmount = context_auc.CheckPrice(auction.BidStartAmount)
	auction.BidDeposit = context_auc.CheckDepositPrice(auction.BidStartAmount)
	auction.Price = context_auc.CheckPrice(auction.Price)

	// 2. 존재하는 경매인지 check
	if _, cnt, err := model.GetDB().GetAucAuction(auction.Id); err != nil {
		log.Error("GetAucAuction :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if cnt == 0 {
			resp.SetReturn(resultcode.Result_DBNotExistAuction)
		} else {
			// 3. auc_product table에 업데이트
			if id, err := model.GetDB().UpdateAucAuction(auction); err != nil {
				log.Error("UpdateAucAuction :", err)
				resp.SetReturn(resultcode.Result_DBError)
			} else {
				if id == 0 {
					resp.SetReturn(resultcode.Result_DBNotExistAuction)
				} else {
					// 가격도 product 테이블에 업데이트 처리
					model.GetDB().UpdateAucProductForPrice(auction.ProductId, auction.TokenType, auction.Price)
					resp.Value = auction
					// 스케줄러 리셋
					schedule.GetScheduler().ResetAuctionScheduler()
				}

			}
		}
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

// 경매 정보 리스트 요청 (경매 상태)
func GetAucAuctionListByAucState(auctionList *context_auc.AuctionListByAucState, c echo.Context) error {
	resp := new(base.BaseResponse)

	// active 경매 정보만 redis에서 가져온다.
	if auctionList.ActiveState == context_auc.Auction_active_state_active {
		//redis exist check
		if pageInfo, auctions, err := model.GetDB().GetAuctionListByAucStateCache(&auctionList.PageInfo, auctionList.AucState); err == nil {
			resp.Success()
			resp.Value = context_auc.AuctionListByAucStateResponse{
				PageInfo:    *pageInfo,
				AucAuctions: *auctions,
			}
			return c.JSON(http.StatusOK, resp)
		}
	}

	auctions, totalCount, err := model.GetDB().GetAucAuctionListByAucState(auctionList)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context_auc.PageInfoResponse{
			PageOffset: auctionList.PageOffset,
			PageSize:   int64(len(auctions)),
			TotalSize:  totalCount,
		}
		resp.Value = context_auc.AuctionListByAucStateResponse{
			PageInfo:    pageInfo,
			AucAuctions: auctions,
		}

		// active 경매 정보만 redis에 남긴다.
		if auctionList.ActiveState == context_auc.Auction_active_state_active {
			model.GetDB().SetAuctionListByAucStateCache(&auctionList.PageInfo, &pageInfo, &auctions, auctionList.AucState)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// 경매 정보 리스트 요청 (추천 경매)
func GetAucAuctionListByRecommand(auctionList *context_auc.AuctionListByRecommand, c echo.Context) error {
	resp := new(base.BaseResponse)

	// active 경매 정보만 redis에서 가져온다.
	if auctionList.ActiveState == context_auc.Auction_active_state_active {
		//redis exist check
		if pageInfo, auctions, err := model.GetDB().GetAuctionListByRecommandCache(&auctionList.PageInfo); err == nil {
			resp.Success()
			resp.Value = context_auc.AuctionListRecommandResponse{
				PageInfo:    *pageInfo,
				AucAuctions: *auctions,
			}
			return c.JSON(http.StatusOK, resp)
		}
	}

	auctions, totalCount, err := model.GetDB().GetAucAuctionListByRecommand(auctionList)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context_auc.PageInfoResponse{
			PageOffset: auctionList.PageOffset,
			PageSize:   int64(len(auctions)),
			TotalSize:  totalCount,
		}
		resp.Value = context_auc.AuctionListRecommandResponse{
			PageInfo:    pageInfo,
			AucAuctions: auctions,
		}

		// active 경매 정보만 redis에 남긴다.
		if auctionList.ActiveState == context_auc.Auction_active_state_active {
			model.GetDB().SetAuctionListByRecommandCache(&auctionList.PageInfo, &pageInfo, &auctions)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// 경매 삭제
func DeleteAucAuctiontRemove(auction *context_auc.RemoveAuction, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_products table 에서 삭제
	if ret, err := model.GetDB().DeleteAucAuction(auction.Id); err != nil {
		log.Error("DeleteAucAuctiontRemove :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if !ret {
			resp.SetReturn(resultcode.Result_DBNotExistProduct)
		}
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 단일 경매 정보 요청
func GetAucAuction(auction *context_auc.GetAuction, c echo.Context) error {
	resp := new(base.BaseResponse)

	//redis exist check
	if auc, err := model.GetDB().GetAuctionCache(auction.Id); err == nil {
		resp.Success()
		resp.Value = auc
		return c.JSON(http.StatusOK, resp)
	}

	auc, count, err := model.GetDB().GetAucAuction(auction.Id)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else if err == nil && count == 0 {
		// 존재하지 않은 경매
		resp.SetReturn(resultcode.Result_DBNotExistAuction)
	} else {
		resp.Success()
		resp.Value = auc

		model.GetDB().SetAuctionCache(auc)
	}

	return c.JSON(http.StatusOK, resp)
}

func PostAucAuctionFinish(auctionFinish *context_auc.AuctionFinish, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 경매 테이블 종료 업데이트
	if affected, err := model.GetDB().UpdateAucAuctionAucState(auctionFinish.Id, context_auc.Auction_auc_state_finish, true); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else if err == nil && affected == 0 {
		resp.SetReturn(resultcode.Result_DBNotExistAuction)
	} else {
		// 스케줄러 리셋
		schedule.GetScheduler().ResetAuctionScheduler()

		// 2. 기존 최고 입찰자 정보 가져오기
		bid, err := model.GetDB().GetAucBidBestAttendee(auctionFinish.Id)
		if err != nil {
			log.Error("PostAucBidSubmit :", err)
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			// 3. 입찰자 리스트 낙찰 업데이트
			if _, err := model.GetDB().UpdateAucBidFinish(bid, context_auc.Bid_state_success); err != nil {
				log.Error("PostAucBidSubmit :", err)
				resp.SetReturn(resultcode.Result_DBError)
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// local func
// 경매 기간중인지 확인
func IsAuctionPeriod(auction *context_auc.AucAuction, aucId int64) bool {
	if auction == nil {
		auc, _, err := model.GetDB().GetAucAuction(aucId)
		if err != nil || auc == nil {
			return false
		}

		auction = auc
	}

	curT := datetime.GetTS2MilliSec()
	if auction.AucStartTs > curT || auction.AucEndTs < curT {
		return false
	}

	return true
}

// 경매 종료되었는지 확인
func IsAuctionEnd(auction *context_auc.AucAuction, aucId int64) bool {
	if auction == nil {
		auc, _, err := model.GetDB().GetAucAuction(aucId)
		if err != nil || auc == nil {
			return false
		}

		auction = auc
	}

	if auction.AucState == context_auc.Auction_auc_state_ready ||
		auction.AucState == context_auc.Auction_auc_state_start ||
		auction.AucState == context_auc.Auction_auc_state_paused {
		return false
	}
	return true
}
