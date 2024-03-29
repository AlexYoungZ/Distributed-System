package kvraft

const (
	OK             = "OK"
	ErrNoKey       = "ErrNoKey"
	ErrWrongLeader = "ErrWrongLeader"
	ErrTimeOut     = "ErrTimeOut"
)

type msgId int64
type Err string

// PutAppendArgs
// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	Op    string // "Put" or "Append"
	// TODO You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	MsgId    msgId
	ClientId int64
}

type PutAppendReply struct {
	Err Err
}

type GetArgs struct {
	Key string
	// TODO You'll have to add definitions here.
	MsgId    msgId
	ClientId int64
}

type GetReply struct {
	Err   Err
	Value string
}
