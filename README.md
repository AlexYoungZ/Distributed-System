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

## Lab 1: MapReduce

<p align="justify">
In this lab you'll build a MapReduce system. You'll implement a worker process that calls application Map and Reduce functions and handles reading and writing files, and a master process that hands out tasks to workers and copes with failed workers. You'll be building something similar to the MapReduce paper.
</p>

<p align="justify">
Your job is to implement a distributed MapReduce, consisting of two programs, the master and the worker. There will be just one master process, and one or more worker processes executing in parallel. In a real system the workers would run on a bunch of different machines, but for this lab you'll run them all on a single machine. The workers will talk to the master via RPC. Each worker process will ask the master for a task, read the task's input from one or more files, execute the task, and write the task's output to one or more files. The master should notice if a worker hasn't completed its task in a reasonable amount of time (for this lab, use ten seconds), and give the same task to a different worker.</p>

* * *

### Distributed System Labs

| Lab Index | Detailed Requirements                                                                        | Quick Link to My Solution                                                                                                             |
|-----------|----------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| Lab 1     | [MapReduce](http://nil.csail.mit.edu/6.824/2020/labs/lab-mr.html)                            | [mr](https://github.com/AlexYoungZ/Distributed-System-6.824/tree/master/src/mr)                                                       |
| Lab 2     | [Raft](http://nil.csail.mit.edu/6.824/2020/labs/lab-raft.html)                               | [raft](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_2)    |
| Lab 3     | [Fault-tolerant Key/Value Service](http://nil.csail.mit.edu/6.824/2020/labs/lab-kvraft.html) | [kvraft](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_3)  |
| Lab 4     | [Sharded Key/Value Service](http://nil.csail.mit.edu/6.824/2020/labs/lab-shard.html)         | [shardkv](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_4) |

* * *


