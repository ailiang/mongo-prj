package data

const (
	TablePlayerInfoKey = "PlayerInfoKey"
	TableReputationKey = "PlayerReputationKey"
)

type PlayerInfo struct {
	PlayerId   int64
	PlayerName string
}

type ReputationData struct {
	RepuMap map[int32]int32
}
