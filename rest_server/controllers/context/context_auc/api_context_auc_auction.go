package context_auc

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type AucAuction struct {
	Id             int64        `json:"auc_id"`
	BidStartAmount float64      `json:"bid_start_amount"`
	BidCurAmount   float64      `json:"bid_cur_amount"`
	BidUnit        float64      `json:"bid_unit"`
	BidDeposit     float64      `json:"bid_deposit"`
	AucStartTs     int64        `json:"auc_start_ts"`
	AucEndTs       int64        `json:"auc_end_ts"`
	AucState       int64        `json:"auc_state"`
	AucRound       int64        `json:"auc_round"`
	CreateTs       int64        `json:"create_ts"`
	ActiveState    int64        `json:"active_state"`
	ProductId      int64        `json:"product_id"`
	Recommand      int64        `json:"recommand"`
	ProductInfo    *ProductInfo `json:"product,omitempty"`
}

// auc 정보 생성
type AucAuctionRegister struct {
	AucAuction
}

func NewAucAuctionRegister() *AucAuctionRegister {
	return new(AucAuctionRegister)
}

func (o *AucAuctionRegister) CheckValidate() *base.BaseResponse {
	if o.BidStartAmount == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireBidStartAmount)
	}
	if o.AucStartTs == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireStartTs)
	}
	if o.AucEndTs == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireEndTs)
	}
	if o.AucRound != 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireRound)
	}
	if o.ActiveState < 0 && o.ActiveState > 1 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireActiveState)
	}
	if o.ProductId == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireProductId)
	}
	return nil
}

////////////////////////////////////////////////

// 경매 정보 업데이트
type AucAuctionUpdate struct {
	AucAuction
}

func NewAucAuctionUpdate() *AucAuctionUpdate {
	return new(AucAuctionUpdate)
}

func (o *AucAuctionUpdate) CheckValidate() *base.BaseResponse {
	if o.Id == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireAucId)
	}
	if o.BidStartAmount == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireBidStartAmount)
	}
	if o.AucStartTs == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireStartTs)
	}
	if o.AucEndTs == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireEndTs)
	}
	if o.AucRound != 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireRound)
	}
	if o.ActiveState < 0 && o.ActiveState > 1 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireActiveState)
	}
	if o.ProductId == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireProductId)
	}
	return nil
}

////////////////////////////////////////////////

// auction list request
type AuctionList struct {
	PageInfo
	ActiveState int64 `query:"active_state"`
}

func NewAuctionList() *AuctionList {
	return new(AuctionList)
}

func (o *AuctionList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type AuctionListResponse struct {
	PageInfo    PageInfoResponse `json:"page_info"`
	AucAuctions []AucAuction     `json:"auctions"`
}

////////////////////////////////////////////////

// auction list by auc_state request
type AuctionListByAucState struct {
	PageInfo
	AucState    int64 `query:"auc_state"`
	ActiveState int64 // Auction_active_state_all, Auction_active_state_inactive, Auction_active_state_active
}

func NewAuctionListByAucState() *AuctionListByAucState {
	return new(AuctionListByAucState)
}

func (o *AuctionListByAucState) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type AuctionListByAucStateResponse struct {
	PageInfo    PageInfoResponse `json:"page_info"`
	AucAuctions []AucAuction     `json:"auctions"`
}

////////////////////////////////////////////////

// auction list by recommand request
type AuctionListByRecommand struct {
	PageInfo
	ActiveState int64 // Auction_active_state_all, Auction_active_state_inactive, Auction_active_state_active
}

func NewAuctionListRecommand() *AuctionListByRecommand {
	return new(AuctionListByRecommand)
}

func (o *AuctionListByRecommand) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type AuctionListRecommandResponse struct {
	PageInfo    PageInfoResponse `json:"page_info"`
	AucAuctions []AucAuction     `json:"auctions"`
}

////////////////////////////////////////////////

// 경매 삭제
type RemoveAuction struct {
	Id int64 `query:"auc_id"`
}

func NewRemoveAuction() *RemoveAuction {
	return new(RemoveAuction)
}

func (o *RemoveAuction) CheckValidate() *base.BaseResponse {
	if o.Id < 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireAucId)
	}
	return nil
}

////////////////////////////////////////////////

// 단일 경매 정보 요청
type GetAuction struct {
	Id int64 `query:"auc_id"`
}

func NewGetAuction() *GetAuction {
	return new(GetAuction)
}

func (o *GetAuction) CheckValidate() *base.BaseResponse {
	if o.Id < 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Auction_RequireAucId)
	}
	return nil
}
