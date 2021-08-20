package context_auc

import "github.com/ONBUFF-IP-TOKEN/baseapp/base"

type SystemRedisRemove struct {
	AuctionList string `query:"auction_list"`
	ProductList string `query:"product_list"`
	BidList     string `query:"bid_list"`
}

func NewSystemRedisRemove() *SystemRedisRemove {
	return new(SystemRedisRemove)
}

func (o *SystemRedisRemove) CheckValidate() *base.BaseResponse {

	return nil
}
