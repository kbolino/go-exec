# Execution framework for Go

[![GoDoc](https://godoc.org/github.com/kbolino/go-exec?status.svg)](https://godoc.org/github.com/kbolino/go-exec)

This is a simple local execution framework for Go, enabling the
interchangeable use of synchronous and asynchronous execution strategies.

## Strategies

The core interface is `Strategy` which has a simple signature:
```
type Strategy interface {
	Do(task func())
}
```

There are 4 implementations of the `Strategy` interface:

* `Bounded` uses a semaphore to limit the maximum number of goroutines
* `Direct` runs every task directly on the caller's goroutine
* `Pool` uses a fixed-size pool of goroutines to run queued tasks
* `Unbounded` starts a new goroutine for each task

Except for `Direct`, every implementation provides a `Wait` method which
blocks until all running tasks are done (although you have to call `Stop`
first on a `Pool`).

`Bounded` and `Pool` also provide additional methods for running tasks:

* `Try` won't block at all
* `TryUntil` won't block past a timeout duration

## Testing

There are unit tests with good coverage on all of the functions and methods,
but they are timing sensitive. For accurate results, the tests should not be
run on machines under heavy load.

## Acknowledgements

The semaphore implementation backing `Bounded` is taken from
[Go Language Patterns][go-ptrns-sem].

[go-ptrns-sem]: http://www.golangpatterns.info/concurrency/semaphores
