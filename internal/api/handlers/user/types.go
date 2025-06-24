package user

type statusResponse struct {
	Status string `json:"status" example:"success"`
}

type internalErrResponse struct {
	Error string `json:"error" example:"internal server error"`
}

type invalidBodyErrResponse struct {
	Error string `json:"error" example:"invalid request"`
}
