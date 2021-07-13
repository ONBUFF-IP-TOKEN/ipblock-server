package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type ProductInfo struct {
	Id                int64   `json:"product_id"`
	Title             string  `json:"product_title"`
	Thumbnail         string  `json:"product_thumbnail_url"`
	Price             float64 `json:"product_price"`
	ProductType       string  `json:"product_type"`
	TokenType         string  `json:"token_type"`
	CreateTs          int64   `json:"create_ts"`
	Creator           string  `json:"creator"`
	Desc              string  `json:"description"`
	Content           string  `json:"content"`
	QuantityTotal     int64   `json:"quantity_total"`
	QuantityRemaining int64   `json:"quantity_remaining"`
	State             int64   `json:"state"`
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
