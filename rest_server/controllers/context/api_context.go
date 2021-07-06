package context

import (
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/constant"
)

const (
	Wallet_type_metamask = "metamask"
)

// page info
type PageInfo struct {
	PageOffset int64 `query:"page_offset" validate:"required"`
	PageSize   int64 `query:"page_size" validate:"required"`
}

// page response
type PageInfoResponse struct {
	PageOffset int64 `json:"page_offset"`
	PageSize   int64 `json:"page_size"`
	TotalSize  int64 `json:"total_size"`
}

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

///////////// API ///////////////////////////////
// login 관련 정보
type LoginAuth struct {
	WalletAddr string `json:"wallet_address" validate:"required"`
	Message    string `json:"message" validate:"required"`
	Sign       string `json:"sign" validate:"required"`
}
type LoginParam struct {
	WalletType string    `json:"wallet_type" validate:"required"`
	WalletAuth LoginAuth `json:"wallet_auth" validate:"required"`
}

func NewLoginParam() *LoginParam {
	return new(LoginParam)
}

func (o *LoginParam) CheckValidate() *constant.OnbuffBaseResponse {
	if len(o.WalletType) == 0 && (Wallet_type_metamask != o.WalletType) {
		return constant.MakeOnbuffBaseResponse(constant.Result_Auth_InvalidWalletType)
	}
	if len(o.WalletAuth.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if len(o.WalletAuth.Message) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_Auth_RequireMessage)
	}
	if len(o.WalletAuth.Sign) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_Auth_RequireSign)
	}
	return nil
}

type LoginResponse struct {
	AuthToken  string `json:"auth_token" validate:"required"`
	ExpireDate int64  `json:"expire_date" validate:"required"`
}

/////////////////////////

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

func (o *RegisterItem) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if len(o.Title) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireTitle)
	}
	if len(o.TokenType) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireTokenType)
	}
	if len(o.Thumbnail) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireThumbnailUrl)
	}
	if o.TokenPrice <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidTokenPrice)
	}
	if o.ExpireDate <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidExpireDate)
	}
	if len(o.Creator) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireCreator)
	}
	if len(o.Description) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireDescription)
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

func (o *UnregisterItem) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if o.ItemId <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidItemId)
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

func (o *GetItemList) CheckValidate() *constant.OnbuffBaseResponse {
	if o.PageOffset < 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidPageSize)
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
	WalletAddr string `json:"wallet_address" validate:"required"`
	ItemId     int64  `json:"item_id" validate:"required"`
}

func NewPostPurchaseItem() *PostPurchaseItem {
	return new(PostPurchaseItem)
}

func (o *PostPurchaseItem) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if o.ItemId < 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidItemId)
	}
	return nil
}

type PostPurchaseItemResponse struct {
	ItemId int64  `json:"item_id"`
	TxHash string `json:"tx_hash"`
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

func (o *GetHistoryTransferItem) CheckValidate() *constant.OnbuffBaseResponse {
	if o.PageOffset < 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidPageSize)
	}
	if o.ItemId <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
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

func (o *GetHistoryTransferMe) CheckValidate(ctx *IPBlockServerContext) *constant.OnbuffBaseResponse {
	if o.WalletAddr != ctx.WalletAddr() {
		return constant.MakeOnbuffBaseResponse(constant.Result_InvalidWalletAddress)
	}
	if o.PageOffset < 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireValidPageSize)
	}
	if len(o.WalletAddr) == 0 {
		return constant.MakeOnbuffBaseResponse(constant.Result_RequireWalletAddress)
	}

	return nil
}

type GetHistoryTransferMeResponse struct {
	PageInfo PageInfoResponse      `json:"page_info"`
	Historys []ItemTransferHistory `json:"historys"`
}

/////////////////////////
