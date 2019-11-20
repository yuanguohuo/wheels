package api

import "io"

type IUserFsm interface {
	OnLeaderStart(term, index uint64) error
	OnLeaderStop(term, index uint64) error

	Apply(term, index uint64, command []byte) error
	OnMemberJoined(term, index uint64, member RaftMember) error
	OnMemberLeft(term, index uint64, member RaftMember) error
	OnTruncateLog(term, index uint64) error

	CreateSnapshot(term, index uint64) error
	RecoverSnapshot(term, index uint64) error

	LoadSnapshot(term, index uint64, writer io.Writer) error
	StoreSnapshot(term, index uint64, reader io.Reader) error

	//when leader replicates an command to a follower, it can optionally add an attachment
	//to the command; on the other hand, when a follower received an command, it detaches
	//the attachment from the command;
	//if you don't need this, just return the original command and nil;
	Attach(term, index uint64, command []byte) ([]byte, error)
	Detach(term, index uint64, command []byte) ([]byte, error)

	OnStop(term, index uint64) error
	OnCrash(term, index uint64, err error) error
}
