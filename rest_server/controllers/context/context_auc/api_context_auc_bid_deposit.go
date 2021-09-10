package context_auc

import (
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type Deposit_state int64

const (
	Deposit_state_fail       = 0 // 보증금 확인 안됨
	Deposit_state_checking   = 1 // 보증금 확인중
	Deposit_state_complete   = 2 // 보증금 확인 완료
	Deposit_state_not_refund = 3 // 보증금 반환 금지
)

type BidDeposit struct {
	Id                    int64          `json:"id"`
	AucId                 int64          `query:"auc_id" json:"auc_id"`
	ProductId             int64          `json:"product_id"`
	BidAttendeeWalletAddr string         `query:"bid_attendee_wallet_address" json:"bid_attendee_wallet_address"`
	DepositAmount         float64        `json:"deposit_amount"`
	DepositTxHash         string         `json:"deposit_txhash"`
	DepositState          Deposit_state  `json:"deposit_state"`
	DepositTs             int64          `json:"deposit_ts"`
	TokenType             string         `json:"token_type"`
	TermsOfService        TermsOfService `json:"terms_of_service"`
}

// 입찰 보증금 정보 전송
type BidDepositSubmit struct {
	BidDeposit
}

func NewBidDeposit() *BidDepositSubmit {
	return new(BidDepositSubmit)
}

func (o *BidDepositSubmit) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
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
