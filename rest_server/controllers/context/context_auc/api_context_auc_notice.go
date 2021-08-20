package context_auc

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
)

type Notice struct {
	Id    int64        `json:"id"`
	Title Localization `json:"title"`
	Desc  Localization `json:"desc"`
	Urls  []string     `json:"urls"`
}

// 공지 사항 등록
type NoticeRegister struct {
	Notices []Notice `json:"notices"`
}

func NewNoticeRegister() *NoticeRegister {
	return new(NoticeRegister)
}

func (o *NoticeRegister) CheckValidate() *base.BaseResponse {
	if len(o.Notices) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Notice_RequireNotice)
	}

	return nil
}

////////////////////////////////////////////////

// 공지 사항 리스트 요청
type NoticeList struct {
	PageInfo
}

func NewNoticeList() *NoticeList {
	return new(NoticeList)
}

func (o *NoticeList) CheckValidate() *base.BaseResponse {
	if o.PageOffset < 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageOffset)
	}
	if o.PageSize <= 0 {
		return base.MakeBaseResponse(resultcode.Result_RequireValidPageSize)
	}
	return nil
}

type NoticeListResponse struct {
	PageInfo PageInfoResponse `json:"page_info"`
	Notices  []Notice         `json:"notices"`
}

////////////////////////////////////////////////

// 공지 사항 삭제
type NoticeRemove struct {
	Id int64 `query:"id"`
}

func NewNoticeRemove() *NoticeRemove {
	return new(NoticeRemove)
}

func (o *NoticeRemove) CheckValidate() *base.BaseResponse {
	return nil
}

////////////////////////////////////////////////

// 공지 사항 업데이트
type NoticeUpdate struct {
	Notices []Notice `json:"notices"`
}

func NewNoticeUpdate() *NoticeUpdate {
	return new(NoticeUpdate)
}

func (o *NoticeUpdate) CheckValidate() *base.BaseResponse {
	if len(o.Notices) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auc_Notice_RequireNotice)
	}

	return nil
}

////////////////////////////////////////////////
