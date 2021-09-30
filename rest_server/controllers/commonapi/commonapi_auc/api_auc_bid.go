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
	"github.com/labstack/echo"
)

// 입찰 전송
func PostAucBidSubmit(bidSubmit *context_auc.BidSubmit, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// auc_auctions 테이블에서 경매 정보 불러오기
	auction, _, err := model.GetDB().GetAucAuction(bidSubmit.AucId)
	if err != nil {
		log.Error("GetAucAuction :", err)
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// 경매 진행중인지 확인
	if !IsAuctionPeriod(auction, bidSubmit.AucId) {
		log.Error("Not Auction period")
		resp.SetReturn(resultcode.Result_Auc_Auction_NotPeriod)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// 보증금 납부 확인
	if bid, err := model.GetDB().GetAucBidDeposit(bidSubmit.AucId, bidSubmit.BidAttendeeWalletAddr); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if bid == nil || bid.DepositState == context_auc.Deposit_state_fail {
			//보증금을 지불한적이 없음
			resp.SetReturn(resultcode.Result_Auc_Bid_RequireDeposit)
		}
	}

	// 1. 기존 최고 입찰자 정보 가져오기
	bid, err := model.GetDB().GetAucBidBestAttendee(bidSubmit.AucId)
	if err != nil {
		log.Error("PostAucBidSubmit :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		bidSubmit.ProductId = auction.ProductId
		bidSubmit.BidState = context_auc.Bid_state_submit
		bidSubmit.BidTs = datetime.GetTS2MilliSec()
		bidSubmit.TokenType = auction.ProductInfo.Prices[0].TokenType

		// 2. 입찰 최고가 체크 기존 최고가의 5%이하 10배가 넘지 않도록 한다.
		if bid == nil {
			// 최초 입찰자인경우에는 시작가의 5%이하 10배가 넘지 않도록한다.
			if bidSubmit.BidAmount < auction.BidStartAmount*1.05 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMin)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMin) //최소가 범위 이탈
			} else if bidSubmit.BidAmount > auction.BidStartAmount*10 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMax)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMax) // 최대가 범위 이탈
			} else {
				// 입찰 한 사람이 아무도 없기때문에 바로 입찰 정보 저장
				if _, err := model.GetDB().InsertAucBidSubmit(bidSubmit); err != nil {
					log.Error("InsertAucBidSubmit :", err)
					resp.SetReturn(resultcode.Result_DBError)
				} else {
					// 최고가 정보 업데이트
					model.GetDB().UpdateAucAuctionBestBid(bidSubmit.AucId, bidSubmit.BidAmount)
				}
			}
		} else {
			// 기존 최고가 시작가의 5%이하 10배가 넘지 않도록한다.
			if bidSubmit.BidAmount < bid.BidAmount*1.05 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMin)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMin) //최소가 범위 이탈
			} else if bidSubmit.BidAmount > bid.BidAmount*10 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMax)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMax) // 최대가 범위 이탈
			} else {
				// 이미 내가 최고 입찰자 인지 확인
				if strings.EqualFold(bid.BidAttendeeWalletAddr, ctx.WalletAddr()) {
					resp.SetReturn(resultcode.Result_Auc_Bid_AlreadyBestAttendee)
				} else {
					// 입찰 정보 저장
					bidSubmit.BidState = context_auc.Bid_state_submit
					if _, err := model.GetDB().InsertAucBidSubmit(bidSubmit); err != nil {
						log.Error("InsertAucBidSubmit :", err)
						resp.SetReturn(resultcode.Result_DBError)
					} else {
						// 최고가 정보 업데이트
						model.GetDB().UpdateAucAuctionBestBid(bidSubmit.AucId, bidSubmit.BidAmount)
					}
				}
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostAucBidSubmit2(bidSubmit *context_auc.BidSubmit, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// auc_auctions 테이블에서 경매 정보 불러오기
	auction, _, err := model.GetDB().GetAucAuction(bidSubmit.AucId)
	if err != nil {
		log.Error("GetAucAuction :", err)
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// 경매 진행중인지 확인
	if !IsAuctionPeriod(auction, bidSubmit.AucId) {
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
		bidSubmit.ProductId = auction.ProductId
		bidSubmit.BidState = context_auc.Bid_state_submit
		bidSubmit.BidTs = datetime.GetTS2MilliSec()
		bidSubmit.TokenType = auction.ProductInfo.Prices[0].TokenType

		// 2. 입찰 최고가 체크 기존 최고가의 5%이하 10배가 넘지 않도록 한다.
		if bid == nil {
			// 최초 입찰자인경우에는 시작가의 5%이하 10배가 넘지 않도록한다.
			if bidSubmit.BidAmount < auction.BidStartAmount*1.05 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMin)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMin) //최소가 범위 이탈
			} else if bidSubmit.BidAmount > auction.BidStartAmount*10 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMax)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMax) // 최대가 범위 이탈
			} else {
				// 입찰 한 사람이 아무도 없기때문에 바로 입찰 정보 저장
				if _, err := model.GetDB().InsertAucBidSubmit(bidSubmit); err != nil {
					log.Error("InsertAucBidSubmit :", err)
					resp.SetReturn(resultcode.Result_DBError)
				} else {
					// 최고가 정보 업데이트
					model.GetDB().UpdateAucAuctionBestBid(bidSubmit.AucId, bidSubmit.BidAmount)
				}
			}
		} else {
			// 기존 최고가 시작가의 5%이하 10배가 넘지 않도록한다.
			if bidSubmit.BidAmount < bid.BidAmount*1.05 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMin)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMin) //최소가 범위 이탈
			} else if bidSubmit.BidAmount > bid.BidAmount*10 {
				log.Error(base.ReturnCodeText(resultcode.Result_Auc_Bid_OutofRangeMax)+" : ", bidSubmit.BidAmount)
				resp.SetReturn(resultcode.Result_Auc_Bid_OutofRangeMax) // 최대가 범위 이탈
			} else {
				// 이미 내가 최고 입찰자 인지 확인
				if strings.EqualFold(bid.BidAttendeeWalletAddr, bidSubmit.BidAttendeeWalletAddr) {
					resp.SetReturn(resultcode.Result_Auc_Bid_AlreadyBestAttendee)
				} else {
					// 입찰 정보 저장
					bidSubmit.BidState = context_auc.Bid_state_submit
					if _, err := model.GetDB().InsertAucBidSubmit(bidSubmit); err != nil {
						log.Error("InsertAucBidSubmit :", err)
						resp.SetReturn(resultcode.Result_DBError)
					} else {
						// 최고가 정보 업데이트
						model.GetDB().UpdateAucAuctionBestBid(bidSubmit.AucId, bidSubmit.BidAmount)
					}
				}
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 입찰자 리스트
func GetAucBidList(bidList *context_auc.BidAttendeeList, c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if pageInfo, bids, err := model.GetDB().GetBidListCache(bidList.AucId, &bidList.PageInfo); err == nil {
		resp.Success()
		resp.Value = context_auc.BidListResponse{
			PageInfo: *pageInfo,
			Bids:     *bids,
		}
	} else {
		// cache 에 없다면 db에서 직접 로드
		bids, totalCount, err := model.GetDB().GetAucBidAttendeeList(bidList)
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.Success()
			pageInfo := context_auc.PageInfoResponse{
				PageOffset: bidList.PageOffset,
				PageSize:   int64(len(bids)),
				TotalSize:  totalCount,
			}
			resp.Value = context_auc.BidListResponse{
				PageInfo: pageInfo,
				Bids:     bids,
			}
			model.GetDB().SetBidListCache(bidList.AucId, &bidList.PageInfo, &pageInfo, &bids)
		}
	}

	return c.JSON(http.StatusOK, resp)
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
				if exist, err := model.GetDB().GetAucBidByTxhash(bid.BidWinnerTxHash); err != nil {
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
func PostAucBidWinnerGiveUp(bid *context_auc.BidWinnerGiveup, ctx *context.IPBlockServerContext) error {
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

// 입찰 삭제
func DeleteAucBidRemove(bid *context_auc.BidRemove, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_products table 에서 삭제
	if ret, err := model.GetDB().DeleteAucBid(bid); err != nil {
		log.Error("DeleteAucBid :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if !ret {
			resp.SetReturn(resultcode.Result_DBNotExistProduct)
		}
	}

	if ret, err := model.GetDB().DeleteAucBidDeposit(bid); err != nil {
		log.Error("DeleteAucBidDeposit :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if !ret {
			resp.SetReturn(resultcode.Result_DBNotExistProduct)
		}
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 입찰 보증금 반환 리스트
func GetAucBidDepositRefund(req *context_auc.BidDepositRefundList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 경매 종료되었는지 확인
	if !IsAuctionEnd(nil, req.AucId) {
		log.Error("Auction is not over yet")
		resp.SetReturn(resultcode.Result_Auc_Auction_NotOverYet)
	} else {
		// 2. 낙찰자 제외한 반환 리스트 추출
		// 낙찰자가 낙찰을 취소했을때도 제외시킨다.
		bids, totalCount, err := model.GetDB().GetAucBidDepositRefund(req)
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.Success()
			pageInfo := context_auc.PageInfoResponse{
				PageOffset: req.PageOffset,
				PageSize:   int64(len(bids)),
				TotalSize:  totalCount,
			}
			resp.Value = context_auc.BidDepositRefundListResponse{
				PageInfo: pageInfo,
				Bids:     bids,
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 낙찰 확인
func GetAucBidWinnerVerify(req *context_auc.BidWinnerVerify, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// auc_auctions 테이블에서 경매 정보 불러오기
	auction, _, err := model.GetDB().GetAucAuction(req.AucId)
	if err != nil {
		log.Error("GetAucAuction :", err)
		resp.SetReturn(resultcode.Result_DBError)
		return ctx.EchoContext.JSON(http.StatusOK, resp)
	}

	// 1. 경매 종료되었는지 확인
	if !IsAuctionEnd(auction, req.AucId) {
		log.Error("Auction is not over yet")
		resp.SetReturn(resultcode.Result_Auc_Auction_NotOverYet)
	} else {
		// 2. 낙찰자가 맞는지 확인
		if successBid, err := model.GetDB().GetAucBidAttendee(req.AucId, ctx.WalletAddr()); err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			if successBid != nil && successBid.BidState == context_auc.Bid_state_success {
				bidResp := &context_auc.BidWinnerVerifyResponse{
					Bid:       *successBid,
					TokenType: successBid.TokenType,
					Payment:   successBid.BidAmount - auction.BidDeposit, // 입찰 금액에서 입찰 보증금을 빼고 지불할 금액을 전달한다.
				}
				resp.Success()
				resp.Value = bidResp
			} else {
				resp.SetReturn(resultcode.Result_Auc_Bid_NotWinner)
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
