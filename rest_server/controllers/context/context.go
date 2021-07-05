package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
)

const (
// // ParamFieldID id 입력 파라미터 Key
// ParamFieldID = iota + base.ParamFieldBaseMax + 1
// // ParamFieldTask task 입력 파라미터 Key
// ParamFieldTask
// // ParamFieldDone done 입력 파라미터 Key
// ParamFieldDone
// // ParamFieldDoneType done_type 입력 파라미터 Key
// ParamFieldDoneType
// // ParamFieldMemberNo member_no 입력 파라미터 Key
// ParamFieldMemberNo
// // ParamFieldUrl 요청 url
// ParamFieldUrl
// // ParamFieldRefresh 무조건 새롭게 요청
// ParamFieldRefresh
// // ParamFieldMessageId 연동된 메세지 id
// ParamFieldCustomId
// // ParamFieldRequestTimeout 크롤링 timeout 대기 시간
// ParamFieldRequestTimeout
// // ParamFieldRequestHeaders 크롤링 api 호출 header 정보
// ParamFieldRequestHeaders
// // ParamFieldAsyncResponseUrl 비동기 응답 시 응답 url
// ParamFieldAsyncResponseUrl
)

const (
// DoneTypeAll = iota
// DoneTypeDone
// DoneTypeUndone
// DoneTypeMax
)

// IPBlockServerContext API의 Request Context
type IPBlockServerContext struct {
	*base.BaseContext
}

// NewIPBlockServerContext 새로운 IPBlockserver Context 생성
func NewIPBlockServerContext(baseCtx *base.BaseContext) interface{} {
	if baseCtx == nil {
		return nil
	}

	ctx := new(IPBlockServerContext)
	ctx.BaseContext = baseCtx

	return ctx
}

// AppendRequestParameter BaseContext 이미 정의되어 있는 ReqeustParameters 배열에 등록
func AppendRequestParameter() {
	// base.SetParamMeta(ParamFieldID, base.NewParamMeta(base.ParamTypePath, "id"))
	// base.SetParamMeta(ParamFieldTask, base.NewParamMeta(base.ParamTypeFormValue, "task"))
	// base.SetParamMeta(ParamFieldDone, base.NewParamMeta(base.ParamTypeFormValue, "done"))
	// base.SetParamMeta(ParamFieldMemberNo, base.NewParamMeta(base.ParamTypeFormValue, "member_no"))
	// base.SetParamMeta(ParamFieldUrl, base.NewParamMeta(base.ParamTypeFormValue, "url"))
	// base.SetParamMeta(ParamFieldRefresh, base.NewParamMeta(base.ParamTypeFormValue, "refresh"))
	// base.SetParamMeta(ParamFieldCustomId, base.NewParamMeta(base.ParamTypeFormValue, "custom_id"))
	// base.SetParamMeta(ParamFieldRequestTimeout, base.NewParamMeta(base.ParamTypeFormValue, "request_timeout"))
	// base.SetParamMeta(ParamFieldRequestHeaders, base.NewParamMeta(base.ParamTypeFormValue, "request_headers"))
	// base.SetParamMeta(ParamFieldAsyncResponseUrl, base.NewParamMeta(base.ParamTypeFormValue, "async_response_url"))
}

// // ID id 입력파라미터 값 반환
// func (o *IPBlockServerContext) ID() string {
// 	return o.GetParam(ParamFieldID)
// }

// // Task task 입력파라미터 값 반환
// func (o *IPBlockServerContext) Task() string {
// 	return o.GetParam(ParamFieldTask)
// }

// // Done done 입력파라미터 값 반환
// func (o *IPBlockServerContext) Done() bool {
// 	doneStr := o.GetParam(ParamFieldDone)
// 	b, err := strconv.ParseBool(doneStr)
// 	if err != nil {
// 		log.Error(err)
// 		return false
// 	}
// 	return b
// }

// func (o *IPBlockServerContext) MemberNo() string {
// 	return o.GetParam(ParamFieldMemberNo)
// }

// func (o *IPBlockServerContext) RequestUrl() string {
// 	return o.GetParam(ParamFieldUrl)
// }

// func (o *IPBlockServerContext) Refresh() string {
// 	return o.GetParam(ParamFieldRefresh)
// }

// func (o *IPBlockServerContext) CustomId() string {
// 	return o.GetParam(ParamFieldCustomId)
// }

// func (o *IPBlockServerContext) RequestTimeout() string {
// 	return o.GetParam(ParamFieldRequestTimeout)
// }

// func (o *IPBlockServerContext) RequestHeaders() string {
// 	return o.GetParam(ParamFieldRequestHeaders)
// }

// func (o *IPBlockServerContext) AsyncResponseUrl() string {
// 	return o.GetParam(ParamFieldAsyncResponseUrl)
// }
