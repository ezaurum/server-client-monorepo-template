package models

type Account struct {
	Model64int
	Name string `json:"name"`
}

// AccountRole 계정 권한
type AccountRole struct {
	Model64int
	AccountID int64 `json:"accountID"`
	// RoleID 롤은 미리 정의된 값으로 사용,
	//todo 사용자가 직접 생성할 수 있게 할 것인가?
	// 직접 생성하게 할 경우 Role 모델을 만들어야 함, 만들어진 Role은 만든 사용자만 사용 가능해야 함. 커스텀 Role을 하고 따로 처리해야 함 or 최고 관리자 0을 만들어서 뭔가 처리해야 함.
	RoleID uint `json:"roleID"`
}

type Role struct {
	Model64int
	Name string `json:"name"`
}

// AccountGroupMapping 계정 그룹 매핑 - 계정과 계정 그룹을 매핑
type AccountGroupMapping struct {
	Model64int
	AccountID       int64 `json:"accountID"`
	AccountGroupID  int64 `json:"accountGroupID"`
	AccountPriority int   `json:"accountPriority"`
}

// AccountGroup 계정 그룹, 여러 계정을 묶어서 처리할 때 사용
type AccountGroup struct {
	Model64int
}
