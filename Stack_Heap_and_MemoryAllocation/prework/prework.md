## The Stack {-}

Read 3.7 (Procedures)[^errata] of _CS:APP_; the goal is to understand how function calls, parameters, and local variables interact with "the stack".

As (or after) you read, you might also find it helpful to build and run the visualization[^ncurses] in the included exercise files.

[^errata]:
	The last sentence in the second-to-last paragraph of 3.7.1 should say "Procedure P can pass up to six integral values (i.e., pointers and integers) _in registers_"; see [Errata for CS:APP3e](https://csapp.cs.cmu.edu/3e/errata.html).

[^ncurses]:
	You may need to install the `ncurses` library before you can run `make viz`.

## The Heap {-}

Read the following sections of _CS:APP_:

* 9.9 (Dynamic Memory Allocation), up to and including 9.9.5 (Implementation Issues)
* 9.11 (Common Memory-Related Bugs in C Programs)

## Memory Allocation in Go {-}

Watch [Understanding Allocations: the Stack and the Heap](https://www.youtube.com/watch?v=ZMZpH4yT7M0), a talk on memory allocation in Go. As you watch, consider differences between C and Go; for example, what happens in each language if a function returns a pointer to a local variable?

## Exploration {-}

Try some short experiments / benchmarks to explore some these questions for your own computer:

* What's the largest array that you an allocate:
	* As a local variable (on the stack)?
	* As a global variable?
	* Dynamically (on the heap)?
* How fast can you allocate memory (e.g. how many `malloc` / `free` calls can you handle per second)?
	* How does it depend on the amount of memory being allocated?
	* What if you call `malloc` from multiple threads?
* What's the typical range of memory addresses for the following?
	* Stack memory?
	* Heap memory?
	* Global variables?
	* String constants?
* Do the address ranges agree with the result of:
	* `vmmap <pid>` (OSX)?
	* `cat /proc/<pid>/maps` (Linux)?