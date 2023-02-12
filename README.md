# 6.824: Distributed Systems

* * *

@ MIT
@ Instructor: [Robert Morris](https://pdos.csail.mit.edu/~rtm/)   
@ [Course Website](https://pdos.csail.mit.edu/6.824/schedule.html)

* * *

## Overview

<p align="justify">
6.824 is a core 12-unit graduate subject with lectures, readings, programming labs, an optional project, a mid-term exam, and a final exam. It will present abstractions and implementation techniques for engineering distributed systems. Major topics include fault tolerance, replication, and consistency. Much of the class consists of studying and discussing case studies of distributed systems.
</p>

* * *

### Distributed System Labs

| Lab Index                | Detailed Requirements                                                                        | Quick Link to My Solution                                                                                                             |
|--------------------------|----------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| [Lab 1](#lab1-mapreduce) | [MapReduce](http://nil.csail.mit.edu/6.824/2020/labs/lab-mr.html)                            | [mr](https://github.com/AlexYoungZ/Distributed-System-6.824/tree/master/src/mr)                                                       |
| [Lab 2](#lab2-raft)      | [Raft](http://nil.csail.mit.edu/6.824/2020/labs/lab-raft.html)                               | [raft](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_2)    |
| Lab 3                    | [Fault-tolerant Key/Value Service](http://nil.csail.mit.edu/6.824/2020/labs/lab-kvraft.html) | [kvraft](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_3)  |
| Lab 4                    | [Sharded Key/Value Service](http://nil.csail.mit.edu/6.824/2020/labs/lab-shard.html)         | [shardkv](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_4) |

* * *

## Lab1 MapReduce

<p align="justify">
In this lab you'll build a MapReduce system. You'll implement a worker process that calls application Map and Reduce functions and handles reading and writing files, and a master process that hands out tasks to workers and copes with failed workers. You'll be building something similar to the MapReduce paper.
</p>

<p align="justify">
Your job is to implement a distributed MapReduce, consisting of two programs, the master and the worker. There will be just one master process, and one or more worker processes executing in parallel. In a real system the workers would run on a bunch of different machines, but for this lab you'll run them all on a single machine. The workers will talk to the master via RPC. Each worker process will ask the master for a task, read the task's input from one or more files, execute the task, and write the task's output to one or more files. The master should notice if a worker hasn't completed its task in a reasonable amount of time (for this lab, use ten seconds), and give the same task to a different worker.</p>

* * *

## Lab2 Raft

### Introduction

<p align="justify">

This is the first in a series of labs in which you'll build a fault-tolerant key/value storage system. In this lab
you'll implement Raft, a replicated state machine protocol. In the next lab you'll build a key/value service on top of
Raft. Then you will “shard” your service over multiple replicated state machines for higher performance.

A replicated service achieves fault tolerance by storing complete copies of its state (i.e., data) on multiple replica
servers. Replication allows the service to continue operating even if some of its servers experience failures (crashes
or a broken or flaky network). The challenge is that failures may cause the replicas to hold differing copies of the
data.

Raft organizes client requests into a sequence, called the log, and ensures that all the replica servers see the same
log. Each replica executes client requests in log order, applying them to its local copy of the service's state. Since
all the live replicas see the same log contents, they all execute the same requests in the same order, and thus continue
to have identical service state. If a server fails but later recovers, Raft takes care of bringing its log up to date.
Raft will continue to operate as long as at least a majority of the servers are alive and can talk to each other. If
there is no such majority, Raft will make no progress, but will pick up where it left off as soon as a majority can
communicate again.

In this lab you'll implement Raft as a Go object type with associated methods, meant to be used as a module in a larger
service. A set of Raft instances talk to each other with RPC to maintain replicated logs. Your Raft interface will
support an indefinite sequence of numbered commands, also called log entries. The entries are numbered with index
numbers. The log entry with a given index will eventually be committed. At that point, your Raft should send the log
entry to the larger service for it to execute.

You should follow the design in the extended Raft paper, with particular attention to Figure 2. You'll implement most of
what's in the paper, including saving persistent state and reading it after a node fails and then restarts. You will not
implement cluster membership changes (Section 6). You'll implement log compaction / snapshotting (Section 7) in a later
lab.

You may find this guide useful, as well as this advice about locking and structure for concurrency. For a wider
perspective, have a look at Paxos, Chubby, Paxos Made Live, Spanner, Zookeeper, Harp, Viewstamped Replication, and
Bolosky et al.

This lab is due in three parts. You must submit each part on the corresponding due date.

</p>

* * *

### Part 2A

<p align="justify">
Implement Raft leader election and heartbeats (AppendEntries RPCs with no log entries). The goal for Part 2A is for a single leader to be elected, for the leader to remain the leader if there are no failures, and for a new leader to take over if the old leader fails or if packets to/from the old leader are lost. Run go test -run 2A to test your 2A code.
</p>

* * *

### Part 2B

<p align="justify">
Implement the leader and follower code to append new log entries, so that the go test -run 2B tests pass.</p>

* * *

### Part 2C

<p align="justify">
If a Raft-based server reboots it should resume service where it left off. This requires that Raft keep persistent state that survives a reboot. The paper's Figure 2 mentions which state should be persistent.

A real implementation would write Raft's persistent state to disk each time it changed, and would read the state from
disk when restarting after a reboot. Your implementation won't use the disk; instead, it will save and restore
persistent state from a Persister object (see persister.go). Whoever calls Raft.Make() supplies a Persister that
initially holds Raft's most recently persisted state (if any). Raft should initialize its state from that Persister, and
should use it to save its persistent state each time the state changes. Use the Persister's ReadRaftState() and
SaveRaftState() methods.
</p>

* * *



