package models

import "time"

type BelongToGSU struct {
	BelongToGathering
	BelongToGatheringSession
	BelongToGatheringUser
}

// CheckInLog 출석 기록
type CheckInLog struct {
	Model64int
	BelongToGSU
	CheckInTime time.Time `json:"checkInTime"`
}

// CheckOutLog 퇴장 기록
type CheckOutLog struct {
	Model64int
	BelongToGSU
	CheckOutTime time.Time `json:"checkOutTime"`
}

// AttendanceAuditLog 출석 기록 수정 로그
type AttendanceAuditLog struct {
	Has64intID
	Type string `json:"type"`
	// 이건 FK를 안 만든다. 여러 테이블에서 참조할 수 있어야 하기 때문에
	// Polymorphic Associations 사용
	TargetID int64 `json:"targetID"`
	// Target
	TargetType     string    `json:"targetType"`
	OldValue       string    `json:"oldValue"`
	NewValue       string    `json:"newValue"`
	AuditTimestamp time.Time `json:"auditTimestamp"`
	// todo 이걸 오거나이저로 쓰게 할까?
	OrganizerID uint `json:"organizerID"`
}

// AttendanceLog 출석 정보, 시간 정보
// 보정값은 없도록 함. 시간 보정은 로그를 남기고 출석, 퇴장 기록을 수정하는 방식으로 처리
type AttendanceLog struct {
	Model64int
	BelongToGSU
	CheckInLogID int64 `json:"checkInLogID"`
	// CheckOutLogID 퇴장 기록이 없을 수 있음, null을 허용하는 대신 0을 넣어서 처리, 0은 placeholder 값을 넣어서 문제 없도록 처리
	CheckOutLogID int64         `json:"checkOutLogID"`
	CheckInTime   *time.Time    `json:"checkInTime"`
	CheckOutTime  *time.Time    `json:"checkOutTime"`
	Score         float32       `json:"score"`
	Duration      time.Duration `json:"duration"`
}

// GatheringSessionAttendanceSetting 모임 세션 출석 설정, 출석 방식을 설정
type GatheringSessionAttendanceSetting struct {
	Model64int
	GatheringID        uint   `json:"gatheringID"`
	GatheringSessionID uint   `json:"gatheringSessionID"`
	Type               string `json:"type"`
}

// AttendanceAggregate 출석 정보 합계, 여러 세션의 출석 정보를 합산해야 할 때 사용
type AttendanceAggregate struct {
	Model64int
	GatheringID uint   `json:"gatheringID"`
	Name        string `json:"name"`
}

// AggregateAttendanceLog 출석 정보 합계, 여러 세션의 출석 정보를 합산해야 할 때 사용
type AggregateAttendanceLog struct {
	Model64int
	GatheringID           uint          `json:"gatheringID"`
	AttendanceAggregateID uint          `json:"attendanceAggregateID"`
	CheckInTime           *time.Time    `json:"checkInTime"`
	CheckOutTime          *time.Time    `json:"checkOutTime"`
	Score                 float32       `json:"score"`
	Duration              time.Duration `json:"duration"`
}

// AggregateAttendanceLogAttendanceLog 출석 정보 합계, 여러 세션의 출석 정보를 합산해야 할 때 사용
type AggregateAttendanceLogAttendanceLog struct {
	Model64int
	AttendanceAggregateID    uint `json:"attendanceAggregateID"`
	AggregateAttendanceLogID uint `json:"aggregateAttendanceLogID"`
	AttendanceLogID          uint `json:"attendanceLogID"`
}
