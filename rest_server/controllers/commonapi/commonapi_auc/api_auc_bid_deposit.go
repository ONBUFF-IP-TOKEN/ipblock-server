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

	if bid, err := model.GetDB().GetAucBidDeposit(bid.AucId, bid.BidAttendeeWalletAddr); err != nil {
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
func PostAucBidDeposit(bidDeposit *context_auc.BidDepositSubmit, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// auc_auctions 테이블에서 경매 정보 불러오기
	auction, _, err := model.GetDB().GetAucAuction(bidDeposit.AucId)
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
	if getBidDeposit, err := model.GetDB().GetAucBidDeposit(bidDeposit.AucId, bidDeposit.BidAttendeeWalletAddr); err != nil {
		log.Error("GetAucBidDeposit :", err)
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	} else {
		if getBidDeposit != nil &&
			strings.EqualFold(getBidDeposit.BidAttendeeWalletAddr, bidDeposit.BidAttendeeWalletAddr) &&
			getBidDeposit.DepositState != context_auc.Deposit_state_fail {
			log.Error("Alreay deposit submit auc_id:", getBidDeposit.AucId, " wallet:", bidDeposit.BidAttendeeWalletAddr)
			resp.SetReturn(resultcode.Result_Auc_Bid_AlreadyDeposit)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
	}

	// 지불 hash가 기존에 지불한적이 있었던 해쉬정보인지 체크
	if exist, err := model.GetDB().GetAucBidByTxhash(bidDeposit.DepositTxHash); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	} else {
		if exist {
			log.Error("PostAucBidDeposit errorCode:", resultcode.Result_Reused_Txhash, " message:", base.ReturnCodeText(resultcode.Result_Reused_Txhash))
			resp.SetReturn(resultcode.Result_Reused_Txhash)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}
	}

	bidDeposit.ProductId = auction.ProductId
	bidDeposit.DepositTs = datetime.GetTS2MilliSec()
	bidDeposit.DepositAmount = float64(auction.BidDeposit)
	bidDeposit.DepositState = context_auc.Deposit_state_checking
	bidDeposit.TokenType = auction.ProductInfo.Prices[0].TokenType

	// auc_bids_deposit table에 최초 저장
	if id, err := model.GetDB().InsertAucBidDepsit(bidDeposit); err != nil {
		log.Error("InsertAucBidDepsit :", err)
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
