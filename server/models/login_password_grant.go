package models

// LoginPasswordGrant 로그인 비밀번호 정보
type LoginPasswordGrant struct {
	Model64int
	PasswordString         string               `json:"passwordString"`
	LoginHandleNameGrant   LoginHandleNameGrant `json:"-"`
	LoginHandleNameGrantID int64                `json:"loginHandleNameGrantID"`
}
