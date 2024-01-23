package database

// Computed properties and user preferences per item per user
type ObjectiveView struct {
	Oid           ObjectiveId
	Vid           VersionId
	Uid           UserId
	Degree        NonNegativeNumber
	Depth         NonNegativeNumber
	ReadyToPickUp bool
	Completion    float64
	Fold          bool
}
