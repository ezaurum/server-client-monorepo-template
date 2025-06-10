package models

type BelongToGathering struct {
	GatheringID int64     `json:"gatheringID"`
	Gathering   Gathering `json:"gathering,omitempty"`
}

type BelongToGatheringSession struct {
	GatheringSessionID int64            `json:"gatheringSessionID"`
	GatheringSession   GatheringSession `json:"gatheringSession,omitempty"`
}

type BelongToGatheringUser struct {
	GatheringUserID int64         `json:"gatheringUserID"`
	GatheringUser   GatheringUser `json:"gatheringUser,omitempty"`
}
