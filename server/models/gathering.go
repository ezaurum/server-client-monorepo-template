package models

// Gathering 모임 정보, 모임 타입은 모임, 세미나, 워크샵 등이 있음
type Gathering struct {
	Model64int
	Name string `json:"name"`
	Type string `json:"type"`
}

// GatheringUser 모임에 참여하는 사용자 정보, 모임에 종속
type GatheringUser struct {
	Model64int
	GatheringID int64 `json:"gatheringID"`
}

// GatheringUserAccount 사용자가 따로 계정이 있을 수 있다.
type GatheringUserAccount struct {
	Model64int
	GatheringUserID int64 `json:"gatheringUserID"`
	AccountID       int64 `json:"accountID"`
}

// GatheringUserCredential 모임에 참여하는 사용자의 로그인 정보
type GatheringUserCredential struct {
	Model64int
	//todo 비번로그인?
}

// GatheringSession 모임 세션 정보 - 실제로 출석을 체 - 매일 하는 출석부라면 이게 날짜별로 생성될 거고, 한 번 하는 출석부라면 이게 그 날짜에만 생성될 것. 출석이 각 수업이나 기타 조건이 있다면 거기마다 생성크하는 정보
type GatheringSession struct {
	Model64int
	GatheringID int64 `json:"gatheringID"`
}

//todo
// 사용자의 QR코드 => otp

//todo
// 초대 링크
// 로그인 링크 - 여러 사용자가 써야하므로... policy에 연결?
