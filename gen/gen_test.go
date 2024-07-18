package gen

type CommonReq struct {
	TraceId  string `json:"traceId" bson:"traceId"`
	Username string `json:"username" bson:"username"`
	Currency string `json:"currency" bson:"currency"`
	Token    string `json:"token" bson:"token"`
}

type walletAdjustmentReq struct {
	CommonReq             `bson:",inline"`
	TransactionId         string  `json:"transactionId"`
	ExternalTransactionId string  `json:"externalTransactionId"`
	RoundId               string  `json:"roundId"`
	Amount                float64 `json:"amount"`
	GameCode              string  `json:"gameCode"`
	Timestamp             int64   `json:"timestamp"`
}
