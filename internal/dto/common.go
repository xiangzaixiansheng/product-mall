package dto

type Response struct {
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
	Msg    string `json:"msg"`
	Error  string `json:"error,omitempty"`
}

type DataList struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
}

type PagedList struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type TokenData struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}

type TrackedErrorResponse struct {
	Response
	TrackID string `json:"track_id"`
}

func BuildListResponse(items any, total int64) Response {
	return Response{
		Status: 200,
		Data: DataList{
			Items: items,
			Total: total,
		},
		Msg: "ok",
	}
}

func BuildPagedResponse(items any, total int64, page, pageSize int) Response {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return Response{
		Status: 200,
		Data: PagedList{
			Items:      items,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
		Msg: "ok",
	}
}
