package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
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
