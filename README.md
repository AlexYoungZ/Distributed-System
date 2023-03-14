# 6.824: Distributed Systems

* * *

@ MIT
@ Instructor: [Robert Morris](https://pdos.csail.mit.edu/~rtm/)   
@ [Course Website](https://pdos.csail.mit.edu/6.824/schedule.html)

* * *

## Overview

<p align="justify">
6.824 is a core 12-unit graduate subject with lectures, readings, programming labs, an optional project, a midterm exam, and a final exam. It will present abstractions and implementation techniques for engineering distributed systems. Major topics include fault tolerance, replication, and consistency. Much of the class consists of studying and discussing case studies of distributed systems.
</p>

* * *

### Distributed System Labs

| Lab Index                                      | Detailed Requirements                                                                        | Quick Link to My Solution                                                                                                             |
|------------------------------------------------|----------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| [Lab 1](#lab1-mapreduce)                       | [MapReduce](http://nil.csail.mit.edu/6.824/2020/labs/lab-mr.html)                            | [mr](https://github.com/AlexYoungZ/Distributed-System-6.824/tree/master/src/mr)                                                       |
| [Lab 2](#lab2-raft)                            | [Raft](http://nil.csail.mit.edu/6.824/2020/labs/lab-raft.html)                               | [raft](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_2)    |
| [Lab 3](#lab3-fault-tolerant-keyvalue-service) | [Fault-tolerant Key/Value Service](http://nil.csail.mit.edu/6.824/2020/labs/lab-kvraft.html) | [kvraft](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_3)  |
| [Lab 4](#lab4-sharded-keyvalue-service)        | [Sharded Key/Value Service](http://nil.csail.mit.edu/6.824/2020/labs/lab-shard.html)         | [shardkv](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_4) |

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

## Lab3 Fault-tolerant Key/Value Service

### Introduction

<p align="justify">

In this lab you will build a fault-tolerant key/value storage service using your Raft library from lab 2. Your key/value
service will be a replicated state machine, consisting of several key/value servers that use Raft for replication. Your
key/value service should continue to process client requests as long as a majority of the servers are alive and can
communicate, in spite of other failures or network partitions.

The service supports three operations: Put(key, value), Append(key, arg), and Get(key). It maintains a simple database
of key/value pairs. Keys and values are strings. Put() replaces the value for a particular key in the database, Append(
key, arg) appends arg to key's value, and Get() fetches the current value for a key. A Get for a non-existant key should
return an empty string. An Append to a non-existant key should act like Put. Each client talks to the service through a
Clerk with Put/Append/Get methods. A Clerk manages RPC interactions with the servers.

Your service must provide strong consistency to application calls to the Clerk Get/Put/Append methods. Here's what we
mean by strong consistency. If called one at a time, the Get/Put/Append methods should act as if the system had only one
copy of its state, and each call should observe the modifications to the state implied by the preceding sequence of
calls. For concurrent calls, the return values and final state must be the same as if the operations had executed one at
a time in some order. Calls are concurrent if they overlap in time, for example if client X calls Clerk.Put(), then
client Y calls Clerk.Append(), and then client X's call returns. Furthermore, a call must observe the effects of all
calls that have completed before the call starts (so we are technically asking for linearizability).

Strong consistency is convenient for applications because it means that, informally, all clients see the same state and
they all see the latest state. Providing strong consistency is relatively easy for a single server. It is harder if the
service is replicated, since all servers must choose the same execution order for concurrent requests, and must avoid
replying to clients using state that isn't up to date.

This lab has two parts. In part A, you will implement the service without worrying that the Raft log can grow without
bound. In part B, you will implement snapshots (Section 7 in the paper), which will allow Raft to discard old log
entries. Please submit each part by the respective deadline.
</p>

* * *

### Part A: Key/value service without log compaction

<p align="justify">
Each of your key/value servers ("kvservers") will have an associated Raft peer. Clerks send Put(), Append(), and Get() RPCs to the kvserver whose associated Raft is the leader. The kvserver code submits the Put/Append/Get operation to Raft, so that the Raft log holds a sequence of Put/Append/Get operations. All of the kvservers execute operations from the Raft log in order, applying the operations to their key/value databases; the intent is for the servers to maintain identical replicas of the key/value database.

A Clerk sometimes doesn't know which kvserver is the Raft leader. If the Clerk sends an RPC to the wrong kvserver, or if
it cannot reach the kvserver, the Clerk should re-try by sending to a different kvserver. If the key/value service
commits the operation to its Raft log (and hence applies the operation to the key/value state machine), the leader
reports the result to the Clerk by responding to its RPC. If the operation failed to commit (for example, if the leader
was replaced), the server reports an error, and the Clerk retries with a different server.

Your kvservers should not directly communicate; they should only interact with each other through Raft. For all parts of
Lab 3, you must make sure that your Raft implementation continues to pass all of the Lab 2 tests.</p>
</p>

* * *

### Part B: Key/value service with log compaction

<p align="justify">
As things stand now with your code, a rebooting server replays the complete Raft log in order to restore its state. However, it's not practical for a long-running server to remember the complete Raft log forever. Instead, you'll modify Raft and kvserver to cooperate to save space: from time to time kvserver will persistently store a "snapshot" of its current state, and Raft will discard log entries that precede the snapshot. When a server restarts (or falls far behind the leader and must catch up), the server first installs a snapshot and then replays log entries from after the point at which the snapshot was created. Section 7 of the extended Raft paper outlines the scheme; you will have to design the details.

You must design an interface between your Raft library and your service that allows your Raft library to discard log
entries. You must revise your Raft code to operate while storing only the tail of the log. Raft should discard old log
entries in a way that allows the Go garbage collector to free and re-use the memory; this requires that there be no
reachable references (pointers) to the discarded log entries.

The tester passes maxraftstate to your StartKVServer(). maxraftstate indicates the maximum allowed size of your
persistent Raft state in bytes (including the log, but not including snapshots). You should compare maxraftstate to
persister.RaftStateSize(). Whenever your key/value server detects that the Raft state size is approaching this
threshold, it should save a snapshot, and tell the Raft library that it has snapshotted, so that Raft can discard old
log entries. If maxraftstate is -1, you do not have to snapshot. maxraftstate applies to the GOB-encoded bytes your Raft
passes to persister.SaveRaftState().
</p>

* * *

## Lab4 Sharded Key/Value Service

### Introduction

<p align="justify">

In this lab you'll build a key/value storage system that "shards," or partitions, the keys over a set of replica groups.
A shard is a subset of the key/value pairs; for example, all the keys starting with "a" might be one shard, all the keys
starting with "b" another, etc. The reason for sharding is performance. Each replica group handles puts and gets for
just a few of the shards, and the groups operate in parallel; thus total system throughput (puts and gets per unit time)
increases in proportion to the number of groups.

Your sharded key/value store will have two main components. First, a set of replica groups. Each replica group is
responsible for a subset of the shards. A replica consists of a handful of servers that use Raft to replicate the
group's shards. The second component is the "shard master". The shard master decides which replica group should serve
each shard; this information is called the configuration. The configuration changes over time. Clients consult the shard
master in order to find the replica group for a key, and replica groups consult the master in order to find out what
shards to serve. There is a single shard master for the whole system, implemented as a fault-tolerant service using
Raft.

A sharded storage system must be able to shift shards among replica groups. One reason is that some groups may become
more loaded than others, so that shards need to be moved to balance the load. Another reason is that replica groups may
join and leave the system: new replica groups may be added to increase capacity, or existing replica groups may be taken
offline for repair or retirement.

The main challenge in this lab will be handling reconfiguration -- changes in the assignment of shards to groups. Within
a single replica group, all group members must agree on when a reconfiguration occurs relative to client Put/Append/Get
requests. For example, a Put may arrive at about the same time as a reconfiguration that causes the replica group to
stop being responsible for the shard holding the Put's key. All replicas in the group must agree on whether the Put
occurred before or after the reconfiguration. If before, the Put should take effect and the new owner of the shard will
see its effect; if after, the Put won't take effect and client must re-try at the new owner. The recommended approach is
to have each replica group use Raft to log not just the sequence of Puts, Appends, and Gets but also the sequence of
reconfigurations. You will need to ensure that at most one replica group is serving requests for each shard at any one
time.

Reconfiguration also requires interaction among the replica groups. For example, in configuration 10 group G1 may be
responsible for shard S1. In configuration 11, group G2 may be responsible for shard S1. During the reconfiguration from
10 to 11, G1 and G2 must use RPC to move the contents of shard S1 (the key/value pairs) from G1 to G2.

This lab's general architecture (a configuration service and a set of replica groups) follows the same general pattern
as Flat Datacenter Storage, BigTable, Spanner, FAWN, Apache HBase, Rosebud, Spinnaker, and many others. These systems
differ in many details from this lab, though, and are also typically more sophisticated and capable. For example, the
lab doesn't evolve the sets of peers in each Raft group; its data and query models are very simple; and handoff of
shards is slow and doesn't allow concurrent client access.
</p>

* * *
