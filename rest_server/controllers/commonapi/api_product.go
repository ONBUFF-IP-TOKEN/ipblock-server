package commonapi

import (
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
)

func PostRegisterProduct(product *context.ProductInfo, ctx *context.IPBlockServerContext) error {
	//1. product table에 저장
	model.GetDB().InsertProduct(product)
	//2. nft 생성 요청 하고 product_nft table에 전체 개수만큼 저장
	//3. 전체 개수 만큼 nft 생성되면 product_nft의 token_id를 업데이트 해주고
	//   product table의 state를 2(판매중)로 변경

	return nil
}
