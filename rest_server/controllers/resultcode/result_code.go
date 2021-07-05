package resultcode

const (
	// for application
	ResultLinkViewAlreadyExists    = "11001"
	ResultMemberNoIsRequired       = "11002"
	ResultRequestUrlRequired       = "11003"
	ResultRequestTimeout           = "11004"
	ResultRequestHeaderInvalid     = "11005"
	ResultRequestMaxLengthExceeded = "11006"
	ResultRequestNotAllowUrl       = "11007"
	ResultAsyncResonseUrlRequired  = "11008"
	ResultNotReady                 = "11009"
)

var IPBlockServerResultCodeMap = map[string]string{
	// for application
	ResultLinkViewAlreadyExists:    "LinkView already exists",
	ResultMemberNoIsRequired:       "member_no is required",
	ResultRequestUrlRequired:       "url is required",
	ResultRequestTimeout:           "Request Timeout",
	ResultRequestHeaderInvalid:     "Request headers is invalid",
	ResultRequestMaxLengthExceeded: "Max limit length excedded",
	ResultRequestNotAllowUrl:       "Not Allow url",
	ResultAsyncResonseUrlRequired:  "response url is required",
	ResultNotReady:                 "not ready",
}
