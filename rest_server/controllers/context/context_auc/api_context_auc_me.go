package context_auc

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type MeBid struct {
	Bid         Bid         `json:"bid"`
	ProductInfo ProductInfo `json:"product"`
}

// 내 입찰,낙찰,nft 리스트
type MeBidList struct {
	PageInfo
}

func NewMeBidList() *MeBidList {
	return new(MeBidList)
}

func (o *MeBidList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type MeBidListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	MeBids   []MeBid          `json:"bids"`
}

////////////////////////////////////////////////
