package request

type CustomerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomerBalanceRequest struct {
	Action  UpdateNumberAction `json:"action"`
	Balance float64            `json:"balance"`
}
