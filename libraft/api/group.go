package api

import (
	"context"
	"time"
)

type IRaftGroup interface {
	Propose(ctx context.Context, command []byte) error
	ProposeAndWait(ctx context.Context, command []byte, timeout time.Duration) error

	AddMember(ctx context.Context, member RaftMember) error
	RemoveMember(ctx context.Context, memberId uint64) error
	TransferLeader(ctx context.Context, leader RaftMember) error

	Leader() RaftMember
	Members() []RaftMember

	Committed() (uint64, uint64)
	Applied() (uint64, uint64)
	Snapshotted() (uint64, uint64)
}
