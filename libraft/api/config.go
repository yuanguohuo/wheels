package api

type GroupConfig struct {
	//config items used by libraft
	TickIntervalMs         int64
	SnapshotCount          uint64
	SnapshotCatchUpEntries uint64
	LogWriteSync           bool

	//config items passed to etcd raft
	ElectionTick              int
	HeartbeatTick             int
	MaxSizePerMsg             uint64
	MaxCommittedSizePerReady  uint64
	MaxUncommittedEntriesSize uint64
	MaxInflightMsgs           int
	CheckQuorum               bool
	PreVote                   bool
	DisableProposalForwarding bool
}

type MessengerTypeEnum byte
type StorageType_Enum byte

const (
	MessengerType_gRPC   MessengerTypeEnum = 1
	MessengerType_Thrift MessengerTypeEnum = 2
)
const (
	StorageType_LevelDB StorageType_Enum = 1
)

type GroupManagerConfig struct {
	LocalIp       string
	LocalPort     int32
	MessengerType MessengerTypeEnum
	StorageType   StorageType_Enum
}
