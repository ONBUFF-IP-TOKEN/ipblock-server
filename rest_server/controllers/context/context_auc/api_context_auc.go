package context_auc

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
