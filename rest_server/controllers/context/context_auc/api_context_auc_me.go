package context_auc

import (
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type MeBid struct {
	Bid         Bid         `json:"bid"`
	ProductInfo ProductInfo `json:"product"`
}

// 내 입찰,낙찰,nft 리스트
type MeBidList struct {
	PageInfo
	AucId      int64  `query:"auc_id"`
	WalletAddr string `query:"wallet_address"`
}

func NewMeBidList() *MeBidList {
	return new(MeBidList)
}

func (o *MeBidList) CheckValidate(ctx *context.IPBlockServerContext) *base.BaseResponse {
	if o.AucId < -1 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_RequireAucId)
	}
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	if len(ctx.WalletAddr()) > 0 {
		o.WalletAddr = ctx.WalletAddr()
	}
	return nil
}

type MeBidListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	MeBids   []MeBid          `json:"bids"`
}

////////////////////////////////////////////////

// 내 보유 토큰 정보
type MeTokenAmount struct {
	WalletAddr string  `query:"wallet_address" json:"wallet_address"`
	TokenType  string  `query:"token_type" json:"token_type"`
	Balance    float64 `json:"balance"`
}

func NewMeTokenAmount() *MeTokenAmount {
	return new(MeTokenAmount)
}

func (o *MeTokenAmount) CheckValidate(ctx *context.IPBlockServerContext, checkWallet bool) *base.BaseResponse {
	if len(o.TokenType) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireTokenType)
	}
	if checkWallet && !strings.EqualFold(ctx.WalletAddr(), o.WalletAddr) {
		log.Error("api:", o.WalletAddr, " auth:", ctx.WalletAddr())
		return base.MakeBaseResponse(resultcode.Result_Auc_Bid_InvalidWalletAddress)
	}
	return nil
}
