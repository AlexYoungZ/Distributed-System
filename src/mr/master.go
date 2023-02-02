package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

const (
	TaskStatusReady   = 0
	TaskStatusQueue   = 1
	TaskStatusRunning = 2
	TaskStatusFinish  = 3
	TaskStatusErr     = 4
)

const (
	MaxTaskRunTime   = time.Second * 5
	ScheduleInterval = time.Millisecond * 500
)

type TaskStat struct {
	Status    int
	WorkerId  int
	StartTime time.Time
}

type Master struct {
	// TODO Your definitions here.
	files     []string
	nReduce   int
	taskPhase TaskPhase
	taskStats []TaskStat
	mu        sync.Mutex
	done      bool
	workerSeq int
	taskCh    chan Task
}

func (m *Master) getTask(taskSeq int) Task           {}
func (m *Master) schedule()                          {}
func (m *Master) initMapTask()                       {}
func (m *Master) initReduceTask()                    {}
func (m *Master) regTask(args *TaskArgs, task *Task) {}

// TODO Your code here -- RPC handlers for the worker to call.
func (m *Master) GetOneTask(args *TaskArgs, reply *TaskReply) error             {}
func (m *Master) ReportTask(args *ReportTaskArgs, reply *ReportTaskReply) error {}
func (m *Master) RegWorker(args *RegisterArgs, reply *RegisterReply) error      {}

// start a thread that listens for RPCs from worker.go
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// Done
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
func (m *Master) Done() bool {
	// TODO Your code here. Done
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.done
}

// should assign unique timer to each task
// simplify implementation
func (m *Master) tickSchedule() {
	for !m.Done() {
		go m.schedule()
		time.Sleep(ScheduleInterval)
	}
}

// MakeMaster
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// TODO Your code here.

	m.server()
	return &m
}

// Example
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}
