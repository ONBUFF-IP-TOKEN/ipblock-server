package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

const auction_list_key = "AUCTION-LIST"

type AuctionListCache struct {
	PageInfo    *context_auc.PageInfoResponse `json:"page_info"`
	AuctionList *[]context_auc.AucAuction     `json:"list"`
}

func genCacheKeyByAucAuction(id string) string {
	return config.GetInstance().DBPrefix + ":AUCTION:" + id
}

func (o *DB) DeleteAuctionList() error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyByAucProduct(auction_list_key)
	return o.Cache.Del(cKey)
}

func (o *DB) SetAuctionListCache(reqPageInfo *context_auc.PageInfo, pageInfo *context_auc.PageInfoResponse, data *[]context_auc.AucAuction) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	auctionListCache := &AuctionListCache{
		PageInfo:    pageInfo,
		AuctionList: data,
	}

	cKey := genCacheKeyByAucProduct(auction_list_key)
	field := fmt.Sprintf("%v-%v", reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetAuctionListCache ", field)
	return o.Cache.HSet(cKey, field, auctionListCache)
}

func (o *DB) GetAuctionListCache(pageInfo *context_auc.PageInfo) (*context_auc.PageInfoResponse, *[]context_auc.AucAuction, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	auctionListCache := &AuctionListCache{}
	cKey := genCacheKeyByAucProduct(auction_list_key)
	field := fmt.Sprintf("%v-%v", pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, auctionListCache)

	return auctionListCache.PageInfo, auctionListCache.AuctionList, err
}
