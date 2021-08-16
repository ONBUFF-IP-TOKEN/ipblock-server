package commonapi_auc

import (
	"net/http"
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
)

// 입찰 보증금 지불 여부 확인
func GetAucBidDeposit(bid *context_auc.BidDepositVerify, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if bid, err := model.GetDB().GetAucBidAttendee(bid.AucId, bid.BidAttendeeWalletAddr); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if bid == nil || bid.DepositState == context_auc.Deposit_state_fail {
			//보증금을 지불한적이 없음
			resp.SetReturn(resultcode.Result_Auc_Bid_RequireDeposit)
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 입찰 보증금 지불
func PostAucBidDeposit(bidDeposit *context_auc.BidDeposit, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// auc_auctions 테이블에서 경매 정보 불러오기
	auction, err := model.GetDB().GetAucAuction(bidDeposit.AucId)
	if err != nil {
		log.Error("GetAucAuction :", err)
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// 경매 진행중인지 확인
	if !IsAuctionPeriod(auction, bidDeposit.AucId) {
		log.Error("Not Auction period")
		resp.SetReturn(resultcode.Result_Auc_Auction_NotPeriod)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// 이미 보증금을 지불하였는지 확인
	if bid, err := model.GetDB().GetAucBidAttendee(bidDeposit.AucId, bidDeposit.BidAttendeeWalletAddr); err != nil {
		log.Error("GetAucBidAttendee :", err)
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	} else {
		if bid != nil &&
			strings.EqualFold(bid.BidAttendeeWalletAddr, bidDeposit.BidAttendeeWalletAddr) &&
			bid.DepositState != context_auc.Deposit_state_fail {
			log.Error("Alreay deposit submit auc_id:", bid.AucId, " wallet:", bidDeposit.BidAttendeeWalletAddr)
			resp.SetReturn(resultcode.Result_Auc_Bid_AlreadyDeposit)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
	}

	// 지불 hash가 기존에 지불한적이 있었던 해쉬정보인지 체크
	if exist, err := model.GetDB().GetAucBidBestAttendeeByTxhash(bidDeposit.DepositTxHash); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	} else {
		if exist {
			resp.SetReturn(resultcode.Result_Reused_Txhash)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
	}

	bidDeposit.ProductId = auction.ProductId
	bidDeposit.BidState = context_auc.Bid_state_ready
	bidDeposit.BidTs = datetime.GetTS2MilliSec()
	bidDeposit.DepositAmount = float64(auction.BidDeposit)
	bidDeposit.DepositState = context_auc.Deposit_state_checking
	bidDeposit.TokenType = auction.ProductInfo.Prices[0].TokenType

	// auc_bids table에 최초 저장
	if id, err := model.GetDB().InsertAucBid(bidDeposit); err != nil {
		log.Error("InsertAucBid :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		bidDeposit.Id = id
		resp.Value = bidDeposit

		// 영수증 확인 go 채널로 전달
		data := &basenet.CommandData{
			CommandType: token.TokenCmd_Bid_Deposit_CheckReceipt,
			Data:        bidDeposit,
		}
		commonapi.GetTokenProc(data)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 입찰 전송
func PostAucBidSubmit(bidSubmit *context_auc.BidSubmit, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 경매 진행중인지 확인
	if !IsAuctionPeriod(nil, bidSubmit.AucId) {
		log.Error("Not Auction period")
		resp.SetReturn(resultcode.Result_Auc_Auction_NotPeriod)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// 1. 기존 최고 입찰자 정보 가져오기
	bid, err := model.GetDB().GetAucBidBestAttendee(bidSubmit.AucId)
	if err != nil {
		log.Error("PostAucBidSubmit :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if bid == nil {
			log.Error("not exist bids data")
			resp.SetReturn(resultcode.Result_Auc_Bid_RequireDeposit)
		} else {
			// 2. 이미 내가 최고 입찰자 인지 확인
			if strings.EqualFold(bid.BidAttendeeWalletAddr, ctx.WalletAddr()) {
				resp.SetReturn(resultcode.Result_Auc_Bid_AlreadyBestAttendee)
			} else {
				// 3. 내가 제시한 입찰가가 최고가 인지 확인
				if bidSubmit.BidAmount <= bid.BidAmount {
					resp.SetReturn(resultcode.Result_Auc_Bid_NotBestBidAmount)
				} else {
					// 4. 입찰 성공 (auc_bids table update)
					bidSubmit.BidState = context_auc.Bid_state_submit
					if _, err := model.GetDB().UpdateAucBidSubmit(bidSubmit); err != nil {
						log.Error("UpdateAucBidSubmit :", err)
						resp.SetReturn(resultcode.Result_DBError)
					}
				}
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 입찰자 리스트
func GetAucBidList(pageInfo *context_auc.BidAttendeeList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// cache 에 없다면 db에서 직접 로드
	bids, totalCount, err := model.GetDB().GetAucBidBestAttendeeList(pageInfo)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context_auc.PageInfoResponse{
			PageOffset: pageInfo.PageOffset,
			PageSize:   int64(len(bids)),
			TotalSize:  totalCount,
		}
		resp.Value = context_auc.BidListResponse{
			PageInfo: pageInfo,
			Bids:     bids,
		}
		//model.GetDB().SetProductListCache(&productList.PageInfo, &pageInfo, &products)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 낙찰 받기
func PostAucBidWinnerSubmit(bid *context_auc.BidWinner, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 경매 종료되었는지 확인
	if !IsAuctionEnd(nil, bid.AucId) {
		log.Error("Auction is not over yet")
		resp.SetReturn(resultcode.Result_Auc_Auction_NotOverYet)
	} else {
		// 2. 낙찰자가 맞는지 확인
		if successBid, err := model.GetDB().GetAucBidAttendee(bid.AucId, bid.BidAttendeeWalletAddr); err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			if successBid != nil && successBid.BidState == context_auc.Bid_state_success {
				// 낙찰자 확인 완료
				// 기존 입찰 정보 채우기
				tempBid := *bid
				bid.Bid = *successBid
				bid.Bid.BidWinnerTxHash = tempBid.BidWinnerTxHash

				bid.BidWinnerState = context_auc.Bid_winner_state_submit_checking

				// 지불 hash가 기존에 지불한적이 있었던 해쉬정보인지 체크
				if exist, err := model.GetDB().GetAucBidBestAttendeeByTxhash(bid.BidWinnerTxHash); err != nil {
					resp.SetReturn(resultcode.Result_DBError)
					return ctx.EchoContext.JSON(http.StatusOK, resp)
				} else {
					if exist {
						resp.SetReturn(resultcode.Result_Reused_Txhash)
						return ctx.EchoContext.JSON(http.StatusOK, resp)
					}
				}

				// 2. 낙찰 정보 db 업데이트
				if _, err := model.GetDB().UpdateAucBidWinner(bid); err != nil {
					resp.SetReturn(resultcode.Result_DBError)
				} else {
					// 3. 영수증 확인 go 채널로 전달
					data := &basenet.CommandData{
						CommandType: token.TokenCmd_Bid_Winner_CheckReceipt,
						Data:        bid,
					}
					commonapi.GetTokenProc(data)
				}
			} else {
				// 낙찰자 아님
				resp.SetReturn(resultcode.Result_Auc_Bid_NotWinner)
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 낙찰 포기
func NewBidSuccessGiveup(bid *context_auc.BidWinnerGiveup, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 경매 종료되었는지 확인
	if !IsAuctionEnd(nil, bid.AucId) {
		log.Error("Auction is not over yet")
		resp.SetReturn(resultcode.Result_Auc_Auction_NotOverYet)
	} else {
		// 2. 낙찰자가 맞는지 확인
		if successBid, err := model.GetDB().GetAucBidAttendee(bid.AucId, bid.BidAttendeeWalletAddr); err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			if successBid != nil && successBid.BidState == context_auc.Bid_state_success {
				// 낙찰자 확인 완료
				bid.BidWinnerState = context_auc.Bid_winner_state_giveup

				// 2. 낙찰 포기 정보 db 업데이트
				if _, err := model.GetDB().UpdateAucBidWinnerState(&bid.Bid, context_auc.Bid_winner_state_giveup); err != nil {
					resp.SetReturn(resultcode.Result_DBError)
				}
			} else {
				// 낙찰자 아님
				resp.SetReturn(resultcode.Result_Auc_Bid_NotWinner)
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// local func
// 경매 기간중인지 확인
func IsAuctionPeriod(auction *context_auc.AucAuction, aucId int64) bool {
	if auction == nil {
		auc, err := model.GetDB().GetAucAuction(aucId)
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
		auc, err := model.GetDB().GetAucAuction(aucId)
		if err != nil || auc == nil {
			return false
		}

		auction = auc
	}

	return auction.AucEndTs < datetime.GetTS2MilliSec()
}