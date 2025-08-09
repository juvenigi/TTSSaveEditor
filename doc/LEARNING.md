# Things to learn about golang

- how nulls are properly handled
    - some marshaller implementations are sneaky and create nullable primitives, because they are playing with evil
      reflection
- reliable concurrency
    - channel-based communications for goroutines that may panic: do go devs expect you to add timeouts everywhere? -> I
      suppose you do have deadlock detection
    - callback-based communication
    - wait groups / std sync primitives
    - synchronised shared state

# Things I've learned

## Output stream needs to be flushed

There is almost always hidden buffering, if not in the application code, then at the OS's level.  
Don't ask me why `Close` doesn't flush for you, but you will have a mysterious problem if you aren't diligent of that
matter.

## I overengineer things too eagerly

```
package tabletop_save

func (sm *SaveManager) lsDirForNames() ([]string, error) {
	dir, err := os.ReadDir(sm.properties.GameDir)
	if err != nil {
		return nil, err
	}
	
	var saves = make([]string, len(dir))
	for _, file := range dir {
		name := file.Name()
		if strings.HasSuffix(name, ".json") && strings.Compare(name, "SaveFileInfos.json") == 0 {
			saves = append(saves, filepath.Join(sm.properties.GameDir, name))
		}
	}

	return saves, nil
}
```

## Go's concept of concurrency finally clicked with me

- channels are best made/owned by producers so that they can safely close them, or accept the fact that your producer
  can panic when sending on a closed channel (a somewhat brutal yet effective short-circuiting technique)
- make sure the consumer lives on a separate thread as the producer and the consumer does not block (or add buffering)
- barriers / wait groups are your friend, but something tells me that any non-toy example needs a `context.Context`
    - The use of `context.Context` is similar to `InterruptedException`, except that context is a form of function
      coloring.
    - Yet this coloring is less scary, because barely anything changes if a function simply passes context without doing
      anything with it.
        - adding context into a function is literally go's version of getter setter POJO boilerplate
- you can use callbacks / continuation-style passing
    - or even spinlocks, if you're an engineer of your own destruction.

## As it turns out, Wails handles panics

But this means that I need to listen to its events, otherwise I'll miss out on the knowledge that something went wrong
This also puts some complexity on presenting a coherent state -- I need to invalidate all transient state after a panic.

## Redundancy of CancellableReader for fs ops

Filesystem based io.Reader and io.Writer do not come with ctx variant, because you don't actually need it when dealing
with filesystems. If you want a cancellable write operation, just close the file descriptor in a separate goroutine, and
the write op will error. Best of all, you don't need to clean up an interrupted write, because _most_ operating systems
will delete the incompletely written file for you (which is why when a browser like chrome saves a file, chunks are
committed to disk, this way you have a chance to retry, or at least retain a partial download without the os insta
deleting the file)

However, if you want a cancel operation to be more responsive, make sure the buffer used during copy operation is small
or close the reader first, then the destination.

## Concurrency cancellation

**there are multiple _flavors_ of it:

- calling an async api that returns you a promise, cancellation is possible via a different method
    - Resource revocation
- calling a method that supports a cancellation signal
- simply not waiting for the return of an async
    - resource scope exit
- having a stop flag in a loop (some kind of shared memory approach)
- (brutish) cry to mama: call the os to kill a process/thread

### TODO AI slop (reduce informational noise)

**Flavors of concurrency cancellation**

* **Calling an async API that returns a handle, canceled via another method**

    * Resource revocation (closing file/socket, dropping DB connection)
    * Calling an explicit `.cancel()` or `Abort()` on the returned object
    * Linked cancellation: canceling one handle propagates to others
    * Timeout-driven cancellation tied to the handle
* **Calling a method that supports a cancellation signal**

    * Passing a cancellation token or context object (`context.Context` in Go, `AbortSignal` in JS)
    * Linked tokens/contexts for cascading cancellation
    * Timeout or deadline baked into the signal
    * Structured concurrency scope exit triggers automatic cancellation
* **Simply not waiting for the return of an async (fire-and-forget)**

    * Resource scope exit: exiting the scope implicitly cancels work (structured concurrency, RAII)
    * Losing references to the task/future so it becomes unreachable
    * Parent context cancellation propagates automatically to orphaned work
* **Having a stop flag in a loop (shared memory approach)**

    * Volatile/atomic boolean flag periodically checked by the worker
    * Shared concurrent data structure that indicates shutdown
    * Channel/queue closure as an implicit stop flag (Go, CSP)
    * Interrupt-based flag setting (`Thread.interrupt()`, signal handlers)
* **Cancellation via control flow interruption**

    * Throwing a specific cancellation exception (`OperationCanceledException`, `CancelledError`)
    * Early return from generator/iterator with cleanup (`generator.return()`)
    * Dependency invalidation: upstream task cancellation aborts dependent tasks
