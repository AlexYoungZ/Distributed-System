
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

### Parallel Programming Projects
Project Index | Detailed Requirements | Quick Link to My Solution
--------------- | --------------- | ---------------
Project 1 | [Reciprocal-Array-Sum using the Java Fork/Join Framework](https://www.coursera.org/learn/parallel-programming-in-java/supplement/aXTS6/mini-project-1-reciprocal-array-sum-using-the-java-fork-join-framework) | [miniproject_1](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_1)
Project 2 | [Analyzing Student Statistics Using Java Parallel Streams](https://www.coursera.org/learn/parallel-programming-in-java/supplement/UfGDN/mini-project-2-analyzing-student-statistics-using-java-parallel-streams) | [miniproject_2](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_2)
Project 3 | [Parallelizing Matrix-Matrix Multiply Using Loop Parallelism](https://www.coursera.org/learn/parallel-programming-in-java/supplement/PRuxj/mini-project-3-parallelizing-matrix-matrix-multiply-using-loop-parallelism) | [miniproject_3](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_3)
Project 4 | [Using Phasers to Optimize Data-Parallel Applications](https://www.coursera.org/learn/parallel-programming-in-java/supplement/j9xiZ/mini-project-4-using-phasers-to-optimize-data-parallel-applications) | [miniproject_4](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Parallel%20Programming/miniproject_4)

* * *

## Concurrent Programming

<p align="justify">
This course teaches learners (industry professionals and students) the fundamental concepts of concurrent programming in the context of Java 8.   Concurrent programming enables developers to efficiently and correctly mediate the use of shared resources in parallel programs.  By the end of this course, you will learn how to use basic concurrency constructs in Java such as threads, locks, critical sections, atomic variables, isolation, actors, optimistic concurrency and concurrent collections, as well as their theoretical foundations (e.g., progress guarantees, deadlock, livelock, starvation, linearizability).

The desired learning outcomes of this course are as follows:

•   Concurrency theory: progress guarantees, deadlock, livelock, starvation, linearizability  
•   Use of threads and structured/unstructured locks in Java  
•   Atomic variables and isolation  
•   Optimistic concurrency and concurrent collections in Java (e.g., concurrent queues, concurrent  hashmaps)  
•   Actor model in Java  
</p>

### Concurrent Programming Projects
Project Index | Detailed Requirements | Quick Link to My Solution
--------------- | --------------- | ---------------
Project 1 | [Locking and Synchronization](https://www.coursera.org/learn/concurrent-programming-in-java/supplement/RTnZE/mini-project-1-locking-and-synchronization) | [miniproject_1](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Concurrent%20Programming/miniproject_1)
Project 2 | [Global and Object-Based Isolation](https://www.coursera.org/learn/concurrent-programming-in-java/supplement/4ULtn/mini-project-2-global-and-object-based-isolation) | [miniproject_2](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Concurrent%20Programming/miniproject_2)
Project 3 | [Sieve of Eratosthenes Using Actor Parallelism](https://www.coursera.org/learn/concurrent-programming-in-java/supplement/v5hdr/mini-project-3-sieve-of-eratosthenes-using-actor-parallelism) | [miniproject_3](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Concurrent%20Programming/miniproject_3)
Project 4 | [Parallelization of Boruvka's Minimum Spanning Tree Algorithm](https://www.coursera.org/learn/concurrent-programming-in-java/supplement/aXPoC/mini-project-4-parallelization-of-boruvkas-minimum-spanning-tree-algorithm) | [miniproject_4](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Concurrent%20Programming/miniproject_4)

* * *

## Distributed Programming

<p align="justify">
This course teaches learners (industry professionals and students) the fundamental concepts of Distributed Programming in the context of Java 8.  Distributed programming enables developers to use multiple nodes in a data center to increase throughput and/or reduce latency of selected applications.  By the end of this course, you will learn how to use popular distributed programming frameworks for Java programs, including Hadoop, Spark, Sockets, Remote Method Invocation (RMI), Multicast Sockets, Kafka, Message Passing Interface (MPI), as well as different approaches to combine distribution with multithreading.

The desired learning outcomes of this course are as follows:

•   Distributed map-reduce programming in Java using the Hadoop and Spark frameworks  
•   Client-server programming using Java's Socket and Remote Method Invocation (RMI) interfaces  
•   Message-passing programming in Java using the Message Passing Interface (MPI)  
•   Approaches to combine distribution with multithreading, including processes and threads, distributed actors, and reactive programming  
</p>

### Distributed Programming Projects
Project Index | Detailed Requirements | Quick Link to My Solution
--------------- | --------------- | ---------------
Project 1 | [Page Rank with Spark](https://www.coursera.org/learn/distributed-programming-in-java/supplement/Jy0Hv/mini-project-1-page-rank-with-spark) | [miniproject_1](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Distributed%20Programming/miniproject_1)
Project 2 | [File Server](https://www.coursera.org/learn/distributed-programming-in-java/supplement/rIOck/mini-project-2-file-server) | [miniproject_2](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Distributed%20Programming/miniproject_2)
Project 3 | [Matrix Multiply in MPI](https://www.coursera.org/learn/distributed-programming-in-java/supplement/fqPLv/mini-project-3-matrix-multiply-in-mpi) | [miniproject_3](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Distributed%20Programming/miniproject_3)
Project 4 | [Multi-Threaded File Server](https://www.coursera.org/learn/distributed-programming-in-java/supplement/w6gEN/mini-project-4-multi-threaded-file-server) | [miniproject_4](https://github.com/AlexYoungZ/Parallel-Concurrent-Distributed-Programming/tree/master/Distributed%20Programming/miniproject_4)




