package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
)

type NoticeListCache struct {
	PageInfo   *context_auc.PageInfoResponse `json:"page_info"`
	NoticeList *[]context_auc.Notice         `json:"list"`
}

func (o *DB) DeleteNoticeList() error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	cKey := genCacheKeyNotice("NOTICE-LIST")
	return o.Cache.Del(cKey)
}

func (o *DB) SetNoticeListCache(reqPageInfo *context_auc.PageInfo, pageInfo *context_auc.PageInfoResponse, data *[]context_auc.Notice) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	productListCache := &NoticeListCache{
		PageInfo:   pageInfo,
		NoticeList: data,
	}

	cKey := genCacheKeyNotice("NOTICE-LIST")
	field := fmt.Sprintf("%v-%v", reqPageInfo.PageSize, reqPageInfo.PageOffset)
	log.Info("SetProductListCache ", field)
	return o.Cache.HSet(cKey, field, productListCache)
}

func (o *DB) GetNoticeListCache(pageInfo *context_auc.PageInfo) (*context_auc.PageInfoResponse, *[]context_auc.Notice, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	noticeListCache := &NoticeListCache{}
	cKey := genCacheKeyNotice("NOTICE-LIST")
	field := fmt.Sprintf("%v-%v", pageInfo.PageSize, pageInfo.PageOffset)

	err := o.Cache.HGet(cKey, field, noticeListCache)

	return noticeListCache.PageInfo, noticeListCache.NoticeList, err
}

func genCacheKeyNotice(id string) string {
	return config.GetInstance().DBPrefix + ":AUCTION:" + id
}
