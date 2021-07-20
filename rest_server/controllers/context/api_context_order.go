package context

import (
	"strings"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

const (
	Order_state_txhash_checking       = 0
	Order_state_txhash_complete       = 1
	Order_state_nft_transfer_start    = 2
	Order_state_nft_transfer_complete = 3
	Order_state_cancel                = 4
)

// product order 요청
type OrderProduct struct {
	WalletAddr     string `json:"wallet_address"`
	ProductId      int64  `json:"product_id"`
	PurchaseTxHash string `json:"purchase_tx_hash"`
	CustomerEmail  string `json:"customer_email"`

	QuantityIndex int64
	TokenId       int64
}

func NewOrderProduct() *OrderProduct {
	return new(OrderProduct)
}

func (o *OrderProduct) CheckValidate(ctx *IPBlockServerContext) *base.BaseResponse {
	if len(o.WalletAddr) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if !strings.EqualFold(strings.ToUpper(o.WalletAddr), strings.ToUpper(ctx.WalletAddr())) {
		return base.MakeBaseResponse(resultcode.Result_InvalidWalletAddress)
	}
	if o.ProductId <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	if len(o.PurchaseTxHash) == 0 {
		return base.MakeBaseResponse(resultcode.Result_RequiredPurchaseTxHash)
	}
	return nil
}

type OrderProductesponse struct {
	//PageInfo PageInfoResponse `json:"page_info"`
	//Products []ProductInfo    `json:"products"`
}

////////////////////////////////////////////////

// order 정보
type OrderInfo struct {
	OrderId            int64   `json:"order_id"`
	Date               int64   `json:"order_date"`
	PurchaseTxHash     string  `json:"purchase_tx_hash"`
	State              int64   `json:"state"`
	ProductId          int64   `json:"product_id"`
	Price              float64 `json:"product_price"`
	QuantityIndex      int64   `json:"quantity_index"`
	QuantityTotal      int64   `json:"quantity_total"`
	CustomerWalletAddr string  `json:"customer_wallet_address"`
	CustomerEmail      string  `json:"customer_email"`
	TokenId            int64   `json:"token_id"`
}

func NewOrderInfo() *OrderInfo {
	return new(OrderInfo)
}

////////////////////////////////////////////////
