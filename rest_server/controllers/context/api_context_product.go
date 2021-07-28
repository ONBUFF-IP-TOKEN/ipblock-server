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

//product_prict struct
type ProductPrice struct {
	TokenType string  `json:"token_type"`
	Price     float64 `json:"price"`
}

// product info
type ProductInfo struct {
	Id                int64          `json:"product_id"`
	Title             string         `json:"product_title"`
	Thumbnail         string         `json:"product_thumbnail_url"`
	Prices            []ProductPrice `json:"product_prices"`
	ProductType       string         `json:"product_type"`
	CreateTs          int64          `json:"create_ts"`
	Creator           string         `json:"creator"`
	Desc              string         `json:"description"`
	Content           string         `json:"content"`
	QuantityTotal     int64          `json:"quantity_total"`
	QuantityRemaining int64          `json:"quantity_remaining"`
	State             int64          `json:"state"`
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

// product delete
type UnregisterProduct struct {
	ProductId int64 `query:"product_id"`
}

func NewUnregisterProduct() *UnregisterProduct {
	return new(UnregisterProduct)
}

func (o *UnregisterProduct) CheckValidate() *base.BaseResponse {
	if o.ProductId == 0 {
		return base.MakeBaseResponse(resultcode.Result_Product_RequireVaildId)
	}
	return nil
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

// product list request
type ProductList struct {
	PageInfo
}

func NewProductList() *ProductList {
	return new(ProductList)
}

func (o *ProductList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type ProductListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	Products []ProductInfo    `json:"products"`
}

////////////////////////////////////////////////
