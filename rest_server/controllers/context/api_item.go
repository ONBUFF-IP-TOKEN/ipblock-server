package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

// item info
type ItemInfo struct {
	ItemId          int64   `json:"item_id" validate:"required"`
	WalletAddr      string  `json:"wallet_address" validate:"required"`
	Title           string  `json:"title" validate:"required"`
	TokenType       string  `json:"token_type" validate:"required"`
	Thumbnail       string  `json:"thumbnail_url" validate:"required"`
	TokenPrice      float64 `json:"token_price" validate:"required"`
	ExpireDate      int64   `json:"expire_date" validate:"required"`
	RegisterDate    int64   `json:"register_date" validate:"required"`
	Creator         string  `json:"creator" validate:"required"`
	Description     string  `json:"description" validate:"required"`
	OwnerWalletAddr string  `json:"owner_wallet_address" validate:"required"`
	Owner           string  `json:"owner" validate:"required"`
	TokenId         int64   `json:"token_id"`
	CreateHash      string  `json:"create_hash"`
	Content         string  `json:"content,omitempty"`
}

// item history
type ItemTransferHistory struct {
	Idx       int64  `json:"idx"`
	ItemId    int64  `json:"item_id"`
	FromAddr  string `json:"from_wallet_address"`
	ToAddr    string `json:"to_wallet_address"`
	TokenId   int64  `json:"token_id"`
	State     string `json:"state"`
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
}

func NewItemInfo() *ItemInfo {
	return new(ItemInfo)
}

// register item
type RegisterItem struct {
	WalletAddr  string  `json:"wallet_address" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	TokenType   string  `json:"token_type" validate:"required"`
	Thumbnail   string  `json:"thumbnail_url" validate:"required"`
	TokenPrice  float64 `json:"token_price" validate:"required"`
	ExpireDate  int64   `json:"expire_date" validate:"required"`
	Creator     string  `json:"creator" validate:"required"`
	Description string  `json:"description" validate:"required"`
	CreateHash  string
}

func NewRegisterItem() *RegisterItem {
	return new(RegisterItem)
}

func (o *RegisterItem) CheckValidate(ctx *IPBlockServerContext) *base.BaseResponse {
	if len(o.WalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireWalletAddress)
	}
	if o.WalletAddr != ctx.WalletAddr() {
		return base.MakeBaseResponse(resultcode.Result_InvalidWalletAddress)
	}
	if len(o.Title) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireTitle)
	}
	if len(o.TokenType) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireTokenType)
	}
	if len(o.Thumbnail) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireThumbnailUrl)
	}
	if o.TokenPrice <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidTokenPrice)
	}
	if o.ExpireDate <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidExpireDate)
	}
	if len(o.Creator) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireCreator)
	}
	if len(o.Description) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireDescription)
	}
	return nil
}

type RegisterItemResponse struct {
	ItemId int64  `json:"item_id" validate:"required"`
	TxHash string `json:"tx_hash" validate:"required"`
}

/////////////////////////

// unregister item
type UnregisterItem struct {
	WalletAddr string `query:"wallet_address" validate:"required"`
	ItemId     int64  `query:"item_id" validate:"required"`
}

func NewUnregisterItem() *UnregisterItem {
	return new(UnregisterItem)
}

func (o *UnregisterItem) CheckValidate(ctx *IPBlockServerContext) *base.BaseResponse {
	if len(o.WalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireWalletAddress)
	}
	if o.WalletAddr != ctx.WalletAddr() {
		return base.MakeBaseResponse(resultcode.Result_InvalidWalletAddress)
	}
	if o.ItemId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidItemId)
	}
	return nil
}

type UnregisterItemResponse struct {
	ItemId int64  `json:"item_id" validate:"required"`
	TxHash string `json:"tx_hash" validate:"required"`
}

/////////////////////////

// request item list
type GetItemList struct {
	PageInfo
}

func NewGetItemList() *GetItemList {
	return new(GetItemList)
}

func (o *GetItemList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type GetItemListResponse struct {
	PageInfo  PageInfoResponse `json:"page_info"`
	ItemInfos []ItemInfo       `json:"items"`
}

/////////////////////////

// item 구매 (purchase)
type PostPurchaseItem struct {
	WalletAddr     string `json:"wallet_address" validate:"required"`
	ItemId         int64  `json:"item_id" validate:"required"`
	PurchaseTxHash string `json:"purchase_tx_hash" validate:"required"`
}

func NewPostPurchaseItem() *PostPurchaseItem {
	return new(PostPurchaseItem)
}

func (o *PostPurchaseItem) CheckValidate(ctx *IPBlockServerContext) *base.BaseResponse {
	if len(o.WalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireWalletAddress)
	}
	if o.WalletAddr != ctx.WalletAddr() {
		return base.MakeBaseResponse(resultcode.Result_InvalidWalletAddress)
	}
	if o.ItemId < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidItemId)
	}
	if len(o.PurchaseTxHash) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequiredPurchaseTxHash)
	}
	return nil
}

type PostPurchaseItemResponse struct {
	ItemId int64  `json:"item_id"`
	TxHash string `json:"tx_hash,omitempty"`
}

/////////////////////////

// Item transfer History 조회
type GetHistoryTransferItem struct {
	PageInfo
	ItemId int64 `query:"item_id"`
}

func NewGetHistoryTransferItem() *GetHistoryTransferItem {
	return new(GetHistoryTransferItem)
}

func (o *GetHistoryTransferItem) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	if o.ItemId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireWalletAddress)
	}
	return nil
}

type GetHistoryTransferItemResponse struct {
	PageInfo PageInfoResponse      `json:"page_info"`
	Historys []ItemTransferHistory `json:"historys"`
}

/////////////////////////

// Me transfer History 조회
type GetHistoryTransferMe struct {
	PageInfo
	WalletAddr string `query:"wallet_address"`
}

func NewGetHistoryTransferMe() *GetHistoryTransferMe {
	return new(GetHistoryTransferMe)
}

func (o *GetHistoryTransferMe) CheckValidate(ctx *IPBlockServerContext) *base.BaseResponse {
	if o.WalletAddr != ctx.WalletAddr() {
		return base.MakeBaseResponse(resultcode.Result_InvalidWalletAddress)
	}
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	if len(o.WalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireWalletAddress)
	}

	return nil
}

type GetHistoryTransferMeResponse struct {
	PageInfo PageInfoResponse      `json:"page_info"`
	Historys []ItemTransferHistory `json:"historys"`
}

/////////////////////////
