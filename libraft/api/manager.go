package api

type IRaftGroupManager interface {
	Start() error
	Stop() error

	CreateRaftGroup(groupId uint64, memberId uint64, members []RaftMember, dataDir string, fsm IUserFsm, grpConf *GroupConfig) (IRaftGroup, error)
	LoadRaftGroup(groupId uint64, memberId uint64, dataDir string, fsm IUserFsm, grpConf *GroupConfig) (IRaftGroup, error)
	JoinRaftGroup(groupId uint64, memberId uint64, members []RaftMember, dataDir string, fsm IUserFsm, grpConf *GroupConfig) (IRaftGroup, error)
}
