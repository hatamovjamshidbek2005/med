package schemas

type IDRequest struct {
	ID string `json:"id" binding:"required" example:"UUID"`
}
type IDResponse struct {
	ID string `json:"id"`
}
type ResponseError struct {
	Error  string `json:"data"`
	Status int    `json:"status"`
}
type ResponseSuccess struct {
	Data   any `json:"data"`
	Status int `json:"status"`
}

type GetListRequest struct {
	Search string `json:"search"`
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
}

type GetListRequestOfUserPayload struct {
	ID     string `json:"-" `
	Search string `json:"search"`
	Limit  int64  `json:"limit"`
	Page   int64  `json:"page"`
}
type GetSearchRequest struct {
	Search string `json:"search"`
}
