package models

// LoginHandleNameGrant 로그인 핸들(계정 id등) 정보
type LoginHandleNameGrant struct {
	Model64int
	HandleName string  `json:"handleName"`
	Account    Account `json:"-"`
	AccountID  int64   `json:"accountID"`
}
