package context_auc

import (
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type Bid_state int64

const (
	Bid_state_ready   = 0 // 입찰 준비
	Bid_state_submit  = 1 // 입찰 완료
	Bid_state_success = 2 // 낙찰 성공
	Bid_state_fail    = 3 // 입찰 실패
)

type Deposit_state int64

const (
	Deposit_state_fail     = 0 // 보증금 확인 안됨
	Deposit_state_checking = 1 // 보증금 확인중
	Deposit_state_complete = 2 // 보증금 확인 완료
)

type Bid_winner_state int64

const (
	Bid_winner_state_none            = 0 // 낙찰 결재 안됨
	Bid_winner_state_submit_checking = 1 // 낙찰 결재 확인중
	Bid_winner_state_submit_complete = 2 // 낙찰 결재 확인 완료
	Bid_winner_state_giveup          = 3 // 낙찰 포기
)

type TermsOfService struct {
	DepositAgree string `json:"deposit_agree"`
	PrivacyAgree string `json:"privacy_agree"`
}

type Bid struct {
	Id                    int64          `json:"id"`
	AucId                 int64          `query:"auc_id" json:"auc_id"`
	ProductId             int64          `json:"product_id"`
	BidState              int64          `json:"bid_state"`
	BidTs                 int64          `json:"bid_ts"`
	BidAttendeeWalletAddr string         `query:"bid_attendee_wallet_address" json:"bid_attendee_wallet_address"`
	BidAmount             float64        `json:"bid_amount"`
	BidWinnerTxHash       string         `json:"bid_winner_txhash"`
	BidWinnerState        int64          `json:"bid_winner_state"`
	DepositAmount         float64        `json:"deposit_amount"`
	DepositTxHash         string         `json:"deposit_txhash"`
	DepositState          int64          `json:"deposit_state"`
	TokenType             string         `json:"token_type"`
	TermsOfService        TermsOfService `json:"terms_of_service"`
}

// 입찰 보증금 확인
type BidDepositVerify struct {
	Bid
}

func NewBidDepositVerify() *BidDepositVerify {
	return new(BidDepositVerify)
}

func (o *BidDepositVerify) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if len(o.BidAttendeeWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireWalletAddress)
	}
	if !strings.EqualFold(o.BidAttendeeWalletAddr, ctx.WalletAddr()) {
		log.Error("api:", o.BidAttendeeWalletAddr, " auth:", ctx.WalletAddr())
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_InvalidWalletAddress)
	}

	return nil
}

////////////////////////////////////////////////

// 입찰 보증금 정보 전송
type BidDeposit struct {
	Bid
}

func NewBidDeposit() *BidDeposit {
	return new(BidDeposit)
}

func (o *BidDeposit) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if len(o.BidAttendeeWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireWalletAddress)
	}
	if len(o.DepositTxHash) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireDepositTxHash)
	}
	if !strings.EqualFold(o.BidAttendeeWalletAddr, ctx.WalletAddr()) {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_InvalidWalletAddress)
	}
	if !strings.EqualFold(o.TermsOfService.DepositAgree, "true") {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireDepoistAgree)
	}
	if !strings.EqualFold(o.TermsOfService.PrivacyAgree, "true") {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequirePrivacyAgree)
	}

	return nil
}

////////////////////////////////////////////////

// bid 입찰 신청
type BidSubmit struct {
	Bid
}

func NewBidSubmit() *BidSubmit {
	return new(BidSubmit)
}

func (o *BidSubmit) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if len(o.BidAttendeeWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireWalletAddress)
	}
	if o.BidAmount <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAmount)
	}
	if !strings.EqualFold(o.BidAttendeeWalletAddr, ctx.WalletAddr()) {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_InvalidWalletAddress)
	}

	return nil
}

////////////////////////////////////////////////

// 입찰자 리스트
type BidAttendeeList struct {
	PageInfo
	AucId int64 `query:"auc_id"`
}

func NewBidAttendeeList() *BidAttendeeList {
	return new(BidAttendeeList)
}

func (o *BidAttendeeList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}

	return nil
}

type BidListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	Bids     []Bid            `json:"bids"`
}

////////////////////////////////////////////////

// 낙찰 받기
type BidWinner struct {
	Bid
}

func NewBidSuccess() *BidWinner {
	return new(BidWinner)
}

func (o *BidWinner) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if len(o.BidAttendeeWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireWalletAddress)
	}
	if len(o.BidWinnerTxHash) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireDepositTxHash)
	}
	if !strings.EqualFold(o.BidAttendeeWalletAddr, ctx.WalletAddr()) {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_InvalidWalletAddress)
	}

	return nil
}

////////////////////////////////////////////////

// 낙찰 포기
type BidWinnerGiveup struct {
	Bid
}

func NewBidSuccessGiveup() *BidWinnerGiveup {
	return new(BidWinnerGiveup)
}

func (o *BidWinnerGiveup) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if len(o.BidAttendeeWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireWalletAddress)
	}
	if !strings.EqualFold(o.BidAttendeeWalletAddr, ctx.WalletAddr()) {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_InvalidWalletAddress)
	}

	return nil
}

////////////////////////////////////////////////

// 입찰 삭제
type BidRemove struct {
	Bid
}

func NewBidRemove() *BidRemove {
	return new(BidRemove)
}

func (o *BidRemove) CheckValidate() *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if len(o.BidAttendeeWalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireWalletAddress)
	}

	return nil
}

////////////////////////////////////////////////

// 입찰 보증금 반환 리스트
type BidDepositRefundList struct {
	PageInfo
	AucId int64 `query:"auc_id"`
}

func NewBidDepositRefundList() *BidDepositRefundList {
	return new(BidDepositRefundList)
}

func (o *BidDepositRefundList) CheckValidate() *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}

	return nil
}

type BidDepositRefundListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	Bids     []Bid            `json:"bids"`
}

////////////////////////////////////////////////

// 낙찰 확인
type BidWinnerVerify struct {
	AucId int64 `query:"auc_id"`
}

func NewBidWinnerVerify() *BidWinnerVerify {
	return new(BidWinnerVerify)
}

func (o *BidWinnerVerify) CheckValidate() *base.BaseResponse {
	if o.AucId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	return nil
}

type BidWinnerVerifyResponse struct {
	Bid       Bid     `json:"bid"`
	TokenType string  `json:"token_type"`
	Payment   float64 `json:"payment"`
}

////////////////////////////////////////////////
