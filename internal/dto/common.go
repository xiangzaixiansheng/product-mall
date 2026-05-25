package dto

type Response struct {
	Status int    `json:"status"`
	Data   any    `json:"data"`
	Msg    string `json:"msg"`
	Error  string `json:"error"`
}

type DataList struct {
	Item  any  `json:"item"`
	Total uint `json:"total"`
}

type TokenData struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}

type TrackedErrorResponse struct {
	Response
	TrackID string `json:"track_id"`
}

func BuildListResponse(items any, total uint) Response {
	return Response{
		Status: 200,
		Data: DataList{
			Item:  items,
			Total: total,
		},
		Msg: "ok",
	}
}
