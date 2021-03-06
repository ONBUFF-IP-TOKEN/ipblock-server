package commonapi_auc

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/basenet"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/cdn/azure"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/schedule"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/token"
)

func PostAucProductRegister(product *context_auc.ProductInfo, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//0. 물품 가격 소숫점 3자리이하 버림
	product.Prices[0].Price = context_auc.CheckPrice(product.Prices[0].Price)

	//1. auc_product table에 저장
	product.CreateTs = datetime.GetTS2MilliSec()
	if id, err := model.GetDB().InsertAucProduct(product); err != nil {
		log.Error("InsertProduct :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		product.Id = id

		//2. nft 토큰 생성
		data := &basenet.CommandData{
			CommandType: token.TokenCmd_CreatNftByAut,
			Data:        product,
			Callback:    make(chan interface{}),
		}
		*resp = commonapi.GetTokenProc(data)
		if resp.Return == base.Return_Success {
			resp.Value = product
		} else {
			// nft 생성 실패 하면 입력했던 db 삭제
			model.GetDB().DeleteAucProduct(id)
		}

	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// product 등록하면서 자동으로 경매까지 등록되도록
func PostAucProductRegisterAuction(product *context_auc.AllRegister, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//0. 물품 가격 소숫점 3자리이하 버림
	product.ProductInfo.Prices[0].Price = context_auc.CheckPrice(product.ProductInfo.Prices[0].Price)
	//1. auc_product table에 저장
	product.ProductInfo.CreateTs = datetime.GetTS2MilliSec()
	if id, err := model.GetDB().InsertAucProduct(&product.ProductInfo); err != nil {
		log.Error("InsertProduct :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		product.ProductInfo.Id = id

		//2. nft 토큰 생성
		data := &basenet.CommandData{
			CommandType: token.TokenCmd_CreatNftByAut,
			Data:        &product.ProductInfo,
			Callback:    make(chan interface{}),
		}
		*resp = commonapi.GetTokenProc(data)
		if resp.Return != base.Return_Success {
			// nft 생성 실패 하면 입력했던 db 삭제
			model.GetDB().DeleteAucProduct(id)
			return ctx.EchoContext.JSON(http.StatusOK, resp)
		}

		//3. 경매 등록
		product.AucAuctionRegister.BidStartAmount = product.ProductInfo.Prices[0].Price
		product.AucAuctionRegister.BidCurAmount = 0
		product.AucAuctionRegister.BidDeposit = context_auc.CheckDepositPrice(product.AucAuctionRegister.BidStartAmount)
		// 만약 보증금이 0이 되버리면 0.001로 강제로 고정시킨다.
		if product.AucAuctionRegister.BidDeposit == 0 {
			product.AucAuctionRegister.BidDeposit = 0.001
		}
		if product.AucAuctionRegister.AucStartTs == 0 {
			product.AucAuctionRegister.AucStartTs = datetime.GetTS2MilliSec()
		}
		if product.AucAuctionRegister.AucEndTs == 0 {
			product.AucAuctionRegister.AucEndTs = product.AucAuctionRegister.AucStartTs + 2592000000
		}
		product.AucAuctionRegister.ProductId = id
		product.AucAuctionRegister.TokenType = product.ProductInfo.Prices[0].TokenType
		product.AucAuctionRegister.Price = product.ProductInfo.Prices[0].Price

		//4. auc_auctions table에 저장
		product.AucAuctionRegister.CreateTs = datetime.GetTS2MilliSec()
		if id, err := model.GetDB().InsertAucAuction(&product.AucAuctionRegister); err != nil {
			log.Error("InsertProduct :", err)
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			product.AucAuctionRegister.Id = id
			// 스케줄러 리셋
			schedule.GetScheduler().ResetAuctionScheduler()
		}

		resp.Value = product
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostAucProductUpdate(product *context_auc.ProductInfo, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. 물품 가격 소숫점 3자리이하 버림
	product.Prices[0].Price = context_auc.CheckPrice(product.Prices[0].Price)

	//2. 해당 product 이 존재하는지 check
	if getProduct, err := model.GetDB().GetAucProductById(product.Id); err != nil {
		log.Error("GetAucProductById error : ", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if getProduct == nil {
			log.Error("GetAucProductById invalid product id ")
			resp.SetReturn(resultcode.Result_Auc_Auction_RequireProductId)
		} else {
			//3. auc_product table에 업데이트
			if _, err := model.GetDB().UpdateAucProduct(product); err != nil {
				log.Error("UpdateAucProduct error : ", err)
				resp.SetReturn(resultcode.Result_DBError)
			} else {
				resp.Value = product

				go func() {
					// nft 링크용 json 파일 cdn 업로드
					data := &token.NftUri_AuctionProduct{
						SNo:      product.SNo,
						Media:    product.Media,
						CardInfo: product.CardInfo,
					}

					nftData, _ := json.Marshal(data)
					productId := strconv.FormatInt(product.Id, 10)
					azure.GetAzure().UploadNftInfoBuffer(nftData, productId+"/"+productId+"_0.json")
				}()
			}
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func DeleteAucProductRemove(product *context_auc.RemoveProduct, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_products table 에서 삭제
	if ret, err := model.GetDB().DeleteAucProduct(product.Id); err != nil {
		log.Error("DeleteAucProductRemove :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if !ret {
			resp.SetReturn(resultcode.Result_DBNotExistProduct)
		}
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetAucProductList(productList *context_auc.ProductList, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)

	//redis exist check
	if pageInfo, products, err := model.GetDB().GetProductListCache(&productList.PageInfo); err == nil {
		resp.Success()
		resp.Value = context_auc.ProductListResponse{
			PageInfo: *pageInfo,
			Products: *products,
		}
	} else {
		// cache 에 없다면 db에서 직접 로드
		products, totalCount, err := model.GetDB().GetAucProductList(productList)
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.Success()
			pageInfo := context_auc.PageInfoResponse{
				PageOffset: productList.PageOffset,
				PageSize:   int64(len(products)),
				TotalSize:  totalCount,
			}
			resp.Value = context_auc.ProductListResponse{
				PageInfo: pageInfo,
				Products: products,
			}
			model.GetDB().SetProductListCache(&productList.PageInfo, &pageInfo, &products)
		}
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
