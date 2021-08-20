package model

import (
	"fmt"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

type BidListCache struct {
	PageInfo *context_auc.PageInfoResponse `json:"page_info"`
	BidList  *[]context_auc.Bid            `json:"list"`
}

func (o *DB) SetBidListCache(aucId int64, reqPageInfo *context_auc.PageInfo, pageInfo *context_auc.PageInfoResponse, data *[]context_auc.Bid) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	productListCache := &BidListCache{
		PageInfo: pageInfo,
		BidList:  data,
	}

	cKey := genCacheKeyByBid("BID-LIST:" + strconv.FormatInt(aucId, 10))
	field := fmt.Sprintf("%v-%v", reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetProductListCache ", field)
	return o.Cache.HSet(cKey, field, productListCache)
}

func (o *DB) GetBidListCache(aucId int64, pageInfo *context_auc.PageInfo) (*context_auc.PageInfoResponse, *[]context_auc.Bid, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	bidListCache := &BidListCache{}
	cKey := genCacheKeyByBid("BID-LIST:" + strconv.FormatInt(aucId, 10))
	field := fmt.Sprintf("%v-%v", pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, bidListCache)

	return bidListCache.PageInfo, bidListCache.BidList, err
}

func (o *DB) DeleteBidList(aucId int64) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyByBid("BID-LIST:" + strconv.FormatInt(aucId, 10))
	return o.Cache.Del(cKey)
}

func genCacheKeyByBid(id string) string {
	return config.GetInstance().DBPrefix + ":AUCTION:" + id
}
