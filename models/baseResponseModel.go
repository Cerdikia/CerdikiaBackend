package models

type BaseResponseModel struct {
	// CodeResponse  int         `json:"Code"`
	// HeaderMessage string      `json:"HeaderMessage"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
	Total   int         `json:"total,omitempty"` // total ditampilkan jika ada
}
