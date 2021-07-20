package commonapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
)

func PostRegisterProduct(product *context.ProductInfo, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. product table에 저장
	product.CreateTs = datetime.GetTS2MilliSec()
	product.SetStateRegistering()
	if product.QuantityRemaining == 0 {
		product.QuantityRemaining = product.QuantityTotal
	}
	if id, err := model.GetDB().InsertProduct(product); err != nil {
		log.Error("InsertProduct :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		product.Id = id

		//2. product_nft table에 전체 개수만큼 저장
		data := &basenet.CommandData{
			CommandType: token.TokenCmd_CreateNft,
			Data:        product,
			Callback:    make(chan interface{}),
		}
		*resp = GetTokenProc(data)
		resp.Value = product
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func DeleteUnregisterProduct(product *context.UnregisterProduct, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	ret, err := model.GetDB().DeleteProduct(product)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	}
	if !ret && err == nil {
		//삭제할 product이 없는 경우
		resp.SetReturn(resultcode.Result_DBNotExistItem)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostUpdateProduct(product *context.ProductInfo, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if _, err := model.GetDB().UpdateProduct(product); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostUpdateProductState(product *context.ProductUpdateState, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if _, err := model.GetDB().UpdateProductState(product); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetProductList(productList *context.ProductList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	products, totalCount, err := model.GetDB().GetProductList(productList)
	if err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		resp.Success()
		pageInfo := context.PageInfoResponse{
			PageOffset: productList.PageOffset,
			PageSize:   int64(len(products)),
			TotalSize:  totalCount,
		}
		resp.Value = context.ProductListResponse{
			PageInfo: pageInfo,
			Products: products,
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostProductOrder(order *context.OrderProduct, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)

	if productInfo, err := model.GetDB().GetProductInfo(order.ProductId); err != nil {
		log.Error("PostProductOrder::GetProductInfo errr : ", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		//해당 product의 상태가 판매 중인지 확인
		if productInfo.State == context.Product_state_saleing {
			//해당 product의 잔여 수량 확인
			if productInfo.QuantityRemaining > 0 {
				// 잔여 수량을 1개 줄이고 token thread에게 넘긴다.
				if _, err := model.GetDB().UpdateProductRemain(false, order.ProductId); err != nil {
					log.Error("PostProductOrder::UpdateProductRemain errr : ", err)
					resp.SetReturn(resultcode.Result_DBError)
				} else {
					// 판매 가능 상품 index 추출
					quantityIndex, tokenId, err := GetAvailableNft(productInfo.Id)
					if err != nil {
						resp.SetReturn(resultcode.Result_Product_LackOfQuantity)
						model.GetDB().UpdateProductRemain(true, order.ProductId) //실패시 수량 증가
					} else {
						// 주문서를 작성한다.
						orderInfo := &context.OrderInfo{
							Date:               datetime.GetTS2MilliSec(),
							PurchaseTxHash:     order.PurchaseTxHash,
							State:              context.Order_state_txhash_checking,
							ProductId:          productInfo.Id,
							Price:              productInfo.Price,
							QuantityIndex:      quantityIndex,
							QuantityTotal:      productInfo.QuantityTotal,
							CustomerWalletAddr: order.WalletAddr,
							CustomerEmail:      order.CustomerEmail,
							TokenId:            tokenId,
						}
						order.TokenId = tokenId              //보낼 토큰 정보 담아서 전달
						model.GetDB().InsertOrder(orderInfo) // order table 추가
						data := &basenet.CommandData{
							CommandType: token.TokenCmd_OrderProduct,
							Data:        order,
							Callback:    nil, //콜백은 필요 없다.
						}
						GetTokenProc(data)

						resp.Success()
					}
				}
			} else {
				resp.SetReturn(resultcode.Result_Product_LackOfQuantity)
			}
		} else {
			resp.SetReturn(resultcode.Result_Product_NotOnSale)
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetTokenProc(data *basenet.CommandData) base.BaseResponse {
	if ch, exist := context.GetChanInstance().Get(context.TokenChannel); exist {
		ch.(chan *basenet.CommandData) <- data
	}

	if data.Callback == nil {
		return base.BaseResponse{}
	}

	ticker := time.NewTicker(90 * time.Second)

	resp := base.BaseResponse{}
	select {
	case callback := <-data.Callback:
		ticker.Stop()
		msg, ok := callback.(*base.BaseResponse)
		if ok {
			resp = *msg
		}
	case <-ticker.C:
		ticker.Stop()
		resp = base.BaseResponseInternalServerError()
	}

	return resp
}

func GetAvailableNft(productId int64) (int64, int64, error) {
	nfts, err := model.GetDB().GetNftListByProductId(productId)
	if err != nil {
		log.Error("GetAvailableNft error : ", err)
		return 0, 0, err
	}

	for _, nft := range nfts {
		if nft.OrderState == context.Nft_order_state_sale_ready {
			// 판매 상태를 판매중으로 바꾸고 상품 정보 리턴
			model.GetDB().UpdateProductNftOrderState(nft.TokenId, context.Nft_order_state_saleing)
			return nft.QuantityIndex, nft.TokenId, nil
		}
	}

	return 0, 0, errors.New("not exist nft")
}
