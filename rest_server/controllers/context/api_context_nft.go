package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

const (
	Nft_state_pending = 0
	Nft_state_mint    = 1
)

const (
	Nft_order_state_sale_ready    = 0 // nft 판매 대기중
	Nft_order_state_saleing       = 1 // nft 판매 진행중
	Nft_order_state_sale_complete = 2 // nft 판매 완료
)

type NftInfo struct {
	NftId           int64  `json:"nft_id"`
	ProductId       int64  `json:"product_id"`
	CreateTs        int64  `json:"create_ts"`
	CreateHash      string `json:"create_hash"`
	TokenId         int64  `json:"token_id"`
	QuantityIndex   int64  `json:"quantity_index"`
	OwnerWalletAddr string `json:"owner_wallet_address"`
	NftUri          string `json:"nft_uri"`
	State           int64  `json:"state"`
	OrderState      int64  `json:"order_state"`
}

// nft list 요청
type NftList struct {
	PageInfo
	ProductId int64 `query:"product_id"`
}

func NewNftList() *NftList {
	return new(NftList)
}

func (o *NftList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type NftListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	Nfts     []NftInfo        `json:"nfts"`
}
