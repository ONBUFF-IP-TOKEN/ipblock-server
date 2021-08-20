package commonapi_auc

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/context/context_auc"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/ipblock-server/rest_server/model"
	"github.com/labstack/echo"
)

func PostNoticeRegister(notices *context_auc.NoticeRegister, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if err := model.GetDB().InsertNotice(notices); err != nil {
		log.Error("PostNoticeRegister error : ", err)
		resp.SetReturn(resultcode.Result_DBError)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func GetNotice(noticeList *context_auc.NoticeList, c echo.Context) error {
	resp := new(base.BaseResponse)

	//redis exist check
	if pageInfo, notices, err := model.GetDB().GetNoticeListCache(&noticeList.PageInfo); err == nil {
		resp.Success()
		resp.Value = context_auc.NoticeListResponse{
			PageInfo: *pageInfo,
			Notices:  *notices,
		}
	} else {
		// cache 에 없다면 db에서 직접 로드
		notices, totalCount, err := model.GetDB().GetNotice(noticeList)
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.Success()
			pageInfo := context_auc.PageInfoResponse{
				PageOffset: noticeList.PageOffset,
				PageSize:   int64(len(*notices)),
				TotalSize:  totalCount,
			}
			resp.Value = context_auc.NoticeListResponse{
				PageInfo: pageInfo,
				Notices:  *notices,
			}
			model.GetDB().SetNoticeListCache(&noticeList.PageInfo, &pageInfo, notices)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// 공지 사항 삭제
func DeleteNoticeRemove(notice *context_auc.NoticeRemove, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_products table 에서 삭제
	if ret, err := model.GetDB().DeleteNotice(notice.Id); err != nil {
		log.Error("DeleteNoticeRemove :", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if !ret {
			resp.SetReturn(resultcode.Result_DBNotExistProduct)
		}
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

// 공지 사항 업데이트
func PostNoticeUpdate(notices *context_auc.NoticeUpdate, ctx *context.IPBlockServerContext) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//1. auc_product table에 업데이트
	if err := model.GetDB().UpdateNotice(notices); err != nil {
		log.Error("UpdateAucProduct error : ", err)
		resp.SetReturn(resultcode.Result_DBError)
	}

	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
