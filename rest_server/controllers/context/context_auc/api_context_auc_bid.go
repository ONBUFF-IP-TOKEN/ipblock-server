package context_auc

import (
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

const (
	Bid_state_ready   = 0
	Bid_state_submit  = 1
	Bid_state_success = 2
	Bid_state_giveup  = 3
	Bid_state_fail    = 4

	Deposit_state_checking = 0
	Deposit_state_complete = 1
	Deposit_state_fail     = 2
)

type Bid struct {
	Id                    int64   `json:"id"`
	AucId                 int64   `query:"auc_id" json:"auc_id"`
	ProductId             int64   `json:"product_id"`
	BidState              int64   `json:"bid_state"`
	BidTs                 int64   `json:"bid_ts"`
	BidAttendeeWalletAddr string  `query:"bid_attendee_wallet_address" json:"bid_attendee_wallet_address"`
	BidAmount             float64 `json:"bid_amount"`
	DepositAmount         float64 `json:"deposit_amount"`
	DepositTxHash         string  `json:"deposit_txhash"`
	DepositState          int64   `json:"deposit_state"`
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

func (o *BidAttendeeList) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
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
