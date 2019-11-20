package libraft

import (
	"libraft/api"
	"libraft/impl"
	"libraft/utils"
)

//go:generate protoc --proto_path=impl/pbtypes --go_out=plugins=grpc:impl/pbtypes proposal.proto

var defaultMgrConfig = api.GroupManagerConfig{
	LocalIp:       "127.0.0.1",
	LocalPort:     7890,
	MessengerType: api.MessengerType_Thrift,
	StorageType:   api.StorageType_LevelDB,
}

var defaultGrpConfig = api.GroupConfig{
	TickIntervalMs:         100,
	SnapshotCount:          3000,
	SnapshotCatchUpEntries: 8000,
	LogWriteSync:           true,

	ElectionTick:              10,
	HeartbeatTick:             1,
	MaxSizePerMsg:             0,
	MaxCommittedSizePerReady:  0,
	MaxUncommittedEntriesSize: 0,
	MaxInflightMsgs:           16,
	CheckQuorum:               true,
	PreVote:                   true,
	DisableProposalForwarding: true,
}

func CreateRaftGroupManager2(logger utils.ILogger, magic string) (api.IRaftGroupManager, error) {
	return CreateRaftGroupManager4(logger, magic, defaultMgrConfig, defaultGrpConfig)
}

func CreateRaftGroupManager3(logger utils.ILogger, magic string, config api.GroupManagerConfig) (api.IRaftGroupManager, error) {
	return CreateRaftGroupManager4(logger, magic, config, defaultGrpConfig)
}

func CreateRaftGroupManager4(logger utils.ILogger, magic string, config api.GroupManagerConfig, defGrpConfig api.GroupConfig) (api.IRaftGroupManager, error) {
	return impl.CreateRaftGroupManager(logger, magic, config, defGrpConfig)
}
