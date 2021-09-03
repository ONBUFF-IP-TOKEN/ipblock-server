package context_auc

const (
	TRUE  = "1"
	FALSE = "0"
)

type Localization struct {
	Ko string `json:"ko"`
	En string `json:"en"`
}

type Urls struct {
	Url string `json:"url"`
}

type ProductPrice struct {
	TokenType string  `json:"token_type"`
	Price     float64 `json:"price"`
}

type PageInfo struct {
	PageOffset int64 `query:"page_offset" validate:"required"`
	PageSize   int64 `query:"page_size" validate:"required"`
}

// page response
type PageInfoResponse struct {
	PageOffset int64 `json:"page_offset"`
	PageSize   int64 `json:"page_size"`
	TotalSize  int64 `json:"total_size"`
}

func CheckPrice(value float64) float64 {
	return float64(int(value*1000)) / 1000 // 소숫점 3자리 미만은 버린다.
}

func CheckDepositPrice(value float64) float64 {
	return CheckPrice(value / 10)
}
