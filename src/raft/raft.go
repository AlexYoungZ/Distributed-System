package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"6.824/src/labgob"
	"6.824/src/labrpc"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// import "bytes"
// import "../labgob"

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	rand.Seed(time.Now().Unix())
}

const (
	ElectionTimeout  = time.Millisecond * 300 // election
	HeartBeatTimeout = time.Millisecond * 150 // leader send heartbeat
	ApplyInterval    = time.Millisecond * 100 // apply log
	RPCTimeout       = time.Millisecond * 100
	MaxLockTime      = time.Millisecond * 10 // debug
)

type Role int

const (
	Follower  Role = 0
	Candidate Role = 1
	Leader    Role = 2
)

// ApplyMsg
// as each Raft peer becomes aware that successive log entries are committed,
// the peer should send an ApplyMsg to the service (or tester) on the same server,
// via the applyCh passed to Make().
// set CommandValid to true to indicate that the ApplyMsg contains a newly committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g., snapshots) on the applyCh;
// at that point you can add fields to ApplyMsg,
// but set CommandValid to false for these other uses.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type LogEntry struct {
	Term    int
	Idx     int // only for debug log
	Command interface{}
}

// Raft
// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// TODO Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	role Role
	term int

	electionTimer       *time.Timer
	appendEntriesTimers []*time.Timer
	applyTimer          *time.Timer
	notifyApplyCh       chan struct{}
	stopCh              chan struct{}

	voteFor           int        // server id, -1 for null
	logEntries        []LogEntry // put lastSnapshot at index 0
	applyCh           chan ApplyMsg
	commitIndex       int
	lastSnapshotIndex int
	lastSnapshotTerm  int
	lastApplied       int   // log commit of the server
	nextIndex         []int // next index to send
	matchIndex        []int // check match

	DebugLog  bool      // print log
	lockStart time.Time // debug to find long live lock
	lockEnd   time.Time
	lockName  string
	gid       int
}

// GetState
// return currentTerm and whether this server believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	// TODO Your code here (2A).
	rf.lock("get state")
	defer rf.unlock("get state")
	return rf.term, rf.role == Leader
}

func (rf *Raft) lock(m string) {
	rf.mu.Lock()
	rf.lockStart = time.Now()
	rf.lockName = m
}

func (rf *Raft) unlock(m string) {
	rf.lockEnd = time.Now()
	rf.lockName = ""
	duration := rf.lockEnd.Sub(rf.lockStart)
	if rf.lockName != "" && duration > MaxLockTime {
		rf.log("lock too long:%s:%s:iskill:%v", m, duration, rf.killed())
	}
	rf.mu.Unlock()
}

func (rf *Raft) log(format string, a ...interface{}) {
	if rf.DebugLog == false {
		return
	}
	term, idx := rf.lastLogTermIndex()
	r := fmt.Sprintf(format, a...)
	s := fmt.Sprintf("gid:%d, me: %d, role:%v,term:%d, commitIdx: %v, snidx:%d, apply:%v, matchidx: %v, nextidx:%+v, lastlogterm:%d,idx:%d",
		rf.gid, rf.me, rf.role, rf.term, rf.commitIndex, rf.lastSnapshotIndex, rf.lastApplied, rf.matchIndex, rf.nextIndex, term, idx)
	log.Printf("%s:log:%s\n", s, r)
}

func (rf *Raft) lastLogTermIndex() (int, int) {
	term := rf.logEntries[len(rf.logEntries)-1].Term
	index := rf.lastSnapshotIndex + len(rf.logEntries) - 1
	return term, index
}

func (rf *Raft) getPersistData() []byte {
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(rf.term)
	e.Encode(rf.voteFor)
	e.Encode(rf.commitIndex)
	e.Encode(rf.lastSnapshotIndex)
	e.Encode(rf.lastSnapshotTerm)
	e.Encode(rf.logEntries)
	data := w.Bytes()
	return data
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
func (rf *Raft) persist() {
	// TODO Your code here (2C).
	data := rf.getPersistData()
	rf.persister.SaveRaftState(data)
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

// restore previously persisted state.
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// TODO Your code here (2C).
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}

	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)

	var term int
	var voteFor int
	var logs []LogEntry
	var commitIndex, lastSnapshotIndex, lastSnapshotTerm int

	if d.Decode(&term) != nil ||
		d.Decode(&voteFor) != nil ||
		d.Decode(&commitIndex) != nil ||
		d.Decode(&lastSnapshotIndex) != nil ||
		d.Decode(&lastSnapshotTerm) != nil ||
		d.Decode(&logs) != nil {
		log.Fatal("rf read persist err")
	} else {
		rf.term = term
		rf.voteFor = voteFor
		rf.commitIndex = commitIndex
		rf.lastSnapshotIndex = lastSnapshotIndex
		rf.lastSnapshotTerm = lastSnapshotTerm
		rf.logEntries = logs
	}
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

func (rf *Raft) changeRole(role Role) {
	rf.role = role
	switch role {
	case Follower:
	case Candidate:
		rf.term += 1
		rf.voteFor = rf.me
		rf.resetElectionTimer()
	case Leader:
		_, lastLogIndex := rf.lastLogTermIndex()
		rf.nextIndex = make([]int, len(rf.peers))
		for i := 0; i < len(rf.peers); i++ {
			rf.nextIndex[i] = lastLogIndex + 1
		}
		rf.matchIndex = make([]int, len(rf.peers))
		rf.matchIndex[rf.me] = lastLogIndex
		rf.resetElectionTimer()
	default:
		panic("unknown role")
	}
}

func randElectionTimeout() time.Duration {
	r := time.Duration(rand.Int63()) % ElectionTimeout
	return ElectionTimeout + r
}

// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// Start
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log.
// if this server isn't the leader, returns false.
// otherwise start the agreement and return immediately.
// there is no guarantee that this command will ever be committed to the Raft log,
// since the leader may fail or lose an election.
// even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	// TODO Your code here (2B).
	rf.lock("start")
	term := rf.term
	isLeader := rf.role == Leader
	_, lastIndex := rf.lastLogTermIndex()
	index := lastIndex + 1

	if isLeader {
		rf.logEntries = append(rf.logEntries, LogEntry{
			Term:    rf.term,
			Command: command,
			Idx:     index,
		})
		rf.matchIndex[rf.me] = index
		rf.persist()
	}
	rf.resetHeartBeatTimers()
	rf.unlock("start")
	return index, term, isLeader
}

func (rf *Raft) getLogByIndex(logIndex int) LogEntry {
	idx := logIndex - rf.lastSnapshotIndex
	return rf.logEntries[idx]
}

func (rf *Raft) getRealIdxByLogIndex(logIndex int) int {
	idx := logIndex - rf.lastSnapshotIndex
	if idx < 0 {
		return -1
	} else {
		return idx
	}
}

// Kill
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// TODO Your code here, if desired.
	close(rf.stopCh)
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

func (rf *Raft) startApplyLogs() {
	defer rf.applyTimer.Reset(ApplyInterval)

	rf.lock("applyLogs1")
	var msgs []ApplyMsg
	if rf.lastApplied < rf.lastSnapshotIndex {
		msgs = make([]ApplyMsg, 0, 1)
		msgs = append(msgs, ApplyMsg{
			CommandValid: false,
			Command:      "installSnapShot",
			CommandIndex: rf.lastSnapshotIndex,
		})

	} else if rf.commitIndex <= rf.lastApplied {
		// snapShot didn't update commit idx
		msgs = make([]ApplyMsg, 0)
	} else {
		rf.log("rfapply")
		msgs = make([]ApplyMsg, 0, rf.commitIndex-rf.lastApplied)
		for i := rf.lastApplied + 1; i <= rf.commitIndex; i++ {
			msgs = append(msgs, ApplyMsg{
				CommandValid: true,
				Command:      rf.logEntries[rf.getRealIdxByLogIndex(i)].Command,
				CommandIndex: i,
			})
		}
	}
	rf.unlock("applyLogs1")

	for _, msg := range msgs {
		rf.applyCh <- msg
		rf.lock("applyLogs2")
		rf.log("send applych idx:%d", msg.CommandIndex)
		rf.lastApplied = msg.CommandIndex
		rf.unlock("applyLogs2")
	}
}

// Make
// the service or tester wants to create a Raft server.
// the ports of all the Raft servers (including this one) are in peers[].
// this server's port is peers[me].
// all the servers' peers[] arrays have the same order.
// persister is a place for this server to save its persistent state,
// and also initially holds the most recent saved state, if any.
// applyCh is a channel on which the tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines for any long-running work.
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg, gid ...int) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.applyCh = applyCh

	rf.DebugLog = false
	// gid for test
	if len(gid) != 0 {
		rf.gid = gid[0]
	} else {
		rf.gid = -1
	}

	// TODO Your initialization code here (2A, 2B, 2C).
	// initialize from state persisted before a crash
	rf.stopCh = make(chan struct{})
	rf.term = 0
	rf.voteFor = -1
	rf.role = Follower
	rf.logEntries = make([]LogEntry, 1) // save lastSnapshot at idx ==0
	rf.readPersist(persister.ReadRaftState())

	rf.electionTimer = time.NewTimer(randElectionTimeout())
	rf.appendEntriesTimers = make([]*time.Timer, len(rf.peers))
	for i, _ := range rf.peers {
		rf.appendEntriesTimers[i] = time.NewTimer(HeartBeatTimeout)
	}
	rf.applyTimer = time.NewTimer(ApplyInterval)
	rf.notifyApplyCh = make(chan struct{}, 100)

	// apply log
	go func() {
		for {
			select {
			case <-rf.stopCh:
				return
			case <-rf.applyTimer.C:
				rf.notifyApplyCh <- struct{}{}
			case <-rf.notifyApplyCh:
				rf.startApplyLogs()
			}
		}
	}()

	// start election
	go func() {
		for {
			select {
			case <-rf.stopCh:
				return
			case <-rf.electionTimer.C:
				rf.startElection()
			}
		}
	}()

	// leader send log entry
	for i, _ := range peers {
		if i == rf.me {
			continue
		}
		go func(index int) {
			for {
				select {
				case <-rf.stopCh:
					return
				case <-rf.appendEntriesTimers[index].C:
					rf.appendEntriesToPeer(index)
				}
			}
		}(i)

	}
	// for debug
	//go func() {
	//	for !rf.killed() {
	//		time.Sleep(time.Second * 2)
	//		fmt.Println(fmt.Sprintf("rf who has lock:%s, time:%v", rf.lockName, time.Now().Sub(rf.lockStart)))
	//	}
	//
	//}()
	return rf
}
