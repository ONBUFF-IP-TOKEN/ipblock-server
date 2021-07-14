package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

const (
	Product_state_registering = 0
	Product_state_ready       = 1
	Product_state_saleing     = 2
	Product_state_soldout     = 3
)

const (
	Product_nft_state_pending = 0
	Product_nft_state_mint    = 1
)

// product info
type ProductInfo struct {
	Id                int64   `json:"product_id,omitempty"`
	Title             string  `json:"product_title,omitempty"`
	Thumbnail         string  `json:"product_thumbnail_url,omitempty"`
	Price             float64 `json:"product_price,omitempty"`
	ProductType       string  `json:"product_type,omitempty"`
	TokenType         string  `json:"token_type,omitempty"`
	CreateTs          int64   `json:"create_ts,omitempty"`
	Creator           string  `json:"creator,omitempty"`
	Desc              string  `json:"description,omitempty"`
	Content           string  `json:"content,omitempty"`
	QuantityTotal     int64   `json:"quantity_total,omitempty"`
	QuantityRemaining int64   `json:"quantity_remaining,omitempty"`
	State             int64   `json:"state,omitempty"`
}

func NewProductInfo() *ProductInfo {
	return new(ProductInfo)
}

func (o *ProductInfo) CheckValidate() *base.BaseResponse {
	if len(o.Title) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequiredTitle)
	}
	if len(o.Thumbnail) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequiredThumbnailUrl)
	}
	if len(o.ProductType) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequiredProductType)
	}
	if len(o.ProductType) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequiredProductType)
	}
	if len(o.Creator) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequiredTokenType)
	}
	if len(o.Desc) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequiredDesc)
	}
	if o.QuantityTotal == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequireQuantityTotal)
	}
	return nil
}

func (o *ProductInfo) SetStateRegistering() {
	o.State = Product_state_registering
}

func (o *ProductInfo) SetStateReady() {
	o.State = Product_state_ready
}

func (o *ProductInfo) SetStateSaleing() {
	o.State = Product_state_saleing
}

func (o *ProductInfo) SetStateSoldOut() {
	o.State = Product_state_soldout
}

////////////////////////////////////////////////

// product update state
type ProductUpdateState struct {
	ProductId int64 `json:"product_id"`
	State     int64 `json:"state"`
}

func NewProductUpdateState() *ProductUpdateState {
	return new(ProductUpdateState)
}

func (o *ProductUpdateState) CheckValidate() *base.BaseResponse {
	if o.ProductId == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequireVaildId)
	}
	if o.State < Product_state_registering ||
		o.State > Product_state_soldout {
		return base.MakeBaseResponse(resultcode.Result_Product_RequiredThumbnailUrl)
	}
	return nil
}

////////////////////////////////////////////////
