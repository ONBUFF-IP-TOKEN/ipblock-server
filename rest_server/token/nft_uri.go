package token

type NftUriShose struct {
	Idx           int64  `json:"idx" validate:"required"`
	Name          string `json:"name" validate:"required"`
	SerialNo      string `json:"serial_no" validate:"required"`
	Info          string `json:"info"`
	Certification string `json:"certification" validate:"required"`
}

type NftUriInfo struct {
	UriType  string      `json:"type" validate:"required"`
	CreateTs int64       `json:"create_ts" validate:"required"`
	Data     interface{} `json:"data" validate:"required"`
}
