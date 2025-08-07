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

This felt _inefficient_ because, God forbid I check for `SaveFileInfos.json` every element. Which is why in the current
impl., I collect a slice then copy a slice just to remove that single element. I leave the code as-is, as a testament
to my stupidity, even though a much more reasonable impl. is something like:

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

Filesystem-based io.Reader and io.Writer do not come with ctx variant, because you don't actually need it when dealing
with filesystems. If you want a cancellable write operation, just close the file descriptor in a separate goroutine, and
the write op will error. Best of all, you don't need to clean up an interrupted write, because _most_ operating systems
will delete the incompletely written file for you (which is why when a browser like chrome saves a file, chunks are
committed to disk, this way you have a chance to retry, or at least retain a partial download without the os
insta-deleting the file)

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

## Value vs reference (TODO: better title, restructuring, add more explanation because others don't have the same thoughts in your head as you do)

What conclusions did I draw from this?

1. the most common idiom in C-like languages is to make a zero/nil initialized struct valid, like my JsonPointer. Since
   it is a thin wrapper around strings.Builder, you can manipulate with it even when all references within it are set to
   nil
    - important side-note: Go treats nil as an "empty" and library authors usually make it safe to use. In C, NULL
      usually means "invalid until proven otherwise and I will segfault you if you forget to malloc/init it."
2. creating "constructor functions" that return a valid instance of the struct may sound appealing, until you realize
   that you will be returning it. Untether yourself from OOP shackles and embrace creating a "raw" chunk of memory
   first, then call an `myStruct.Init()` / `myStruct.Setup()` to set a valid initial state.
    - in all fairness, making the "constructor" a receiver method does _look_ more oop-ish in its syntax, but it breaks
      the assumption that the object is uninitialized before the constructor call.
    - this does make me a bit sad, because now I know that the struct being uninitialized is rather unavoidable and
      outright "banning" any use of `var unset MyStruct` is not always possible, because of the unavoidable copy that
      the "constructor method" would do, while returning `*MyStruct` always leaves a disgusting taste in my mouth. If I
      want things to live in the heap, I might as well write a higher-level language and enjoy its complete decadence
      and allocation debauchery.
3. returning a struct by value is fine when it is a dumb transfer object, e.g., a complex number, but the moment you
   introduce some sophisticated state and pointers, you have to be extra careful with handling it. Some things really
   don't like to be copied and defensive programming that has no reference-counting will backfire at you

I am skipping the stack vs heap explanation for now, but you may find a good one here:

https://medium.com/@pranoy1998k/understanding-escape-analysis-in-go-b2db76be58f0

Having programmed _a little bit_ in C, and living through my own fair share of segfaults caused by a null pointer
dereference, I naively thought that I've mastered it all. In truth, it was only the tip of the iceberg. I embarrassingly
have only now discovered why init-free structs are such a hot topic in the C world and why "constructors" are clunky.

I returned to Go knowing the initial footguns: first one being that math-inclined folks fumble when they attempt
treating return values as the function's _image_ or "output domain" because that is not what happens most of the time.
A more appropriate mental framework is to consider the return value as a declaration that the caller needs to provide
more memory (unless we are returning a reference to `this`, builder pattern-style), whereas a function that returns
`void` is doing side effects or mutates the references you pass into it.

Another blunder comes through to the corruption of youth caused by opinionated languages that may or may not even allow
mutability or even memory referencing of any kind, even though the latter is rarer and much more extreme. Most of the
time, "no references" is a blatant lie which is repeatedly told to us when we are children because every object passed
into a method Java is a reference of said object, it would've been too memory-inefficient and impractical otherwise.

If you are unfortunate to share the same fate such as me, remember this:

If your function returns a value, this value gets copied to the caller's stack. One could easily imagine a situation
where the compiler could simply try to be a bit smart about popping the stack by introducing some sort of stack-frame
aliasing/sharing to not require copying, but you have to understand this:

1. The callee's stack frame is ephemeral: once the function returns, its context vanishes. If the caller simply held on
   to that memory that wasn't his, that would be a forbidden dangling pointer and a memory corruption in the making,
   once an unsuspecting new stack-frame will happily take that memory and use it.
2. Abstractions start at the lowest level: Values may be passed back in registers or via caller-allocated buffers,
   something that you may forget when imagining how things look like low-level. Don't forget that C is the lowest
   high-level language after all, in a sense that you no longer see much of the CPU's machinations, meaning that the
   _stack_, _heap_, and _pointers_ are mere illusions (contractual obligations) that can be implemented differently on
   different target platforms. What sounds inefficient at a first glance might actually be necessary to ensure expected
   behavior. Another important mention is returns from dynamically-linked libraries and foreign-function interfaces
   where would be impossible to properly implement copy elision without each component knowing each other a bit too well
   than the minimum expectations from the two foreign parties' interaction. Therefore, when you are returning a value,
   the default is the safest and the most portable option which is to copy, unless the right conditions are met.
3. compilers are allowed to elide copies when possible (Return Value Optimization/move semantics), and much like
   optimizations are prone to backfire, neither are any compiler optimizations guaranteed: if you want to control
   memory ownership and avoid copies explicitly, the idiomatic way in C-like languages is **to pass a reference into
   the function** rather than rely on the compiler to optimize a return.

My blunder was defining JsonPointer struct like so:

```go
package internal

import "strings"

type JsonPointer struct {
	sb strings.Builder
}

// this one breaks
func (jp *JsonPointer) brokenClone() JsonPointer {
	var out JsonPointer
	s := jp.sb.String()
	out.sb.Grow(len(s))
	out.sb.WriteString(s)
	return out
}

// this one also breaks
func alsoBrokenclone(sb *strings.Builder) JsonPointer {
	var out JsonPointer
	out.sb = *sb
	return out
}

// this one works
func (jp *JsonPointer) CloneFrom(src *JsonPointer) {
	jp.sb.Reset()
	jp.sb.Grow(src.sb.Len())
	jp.sb.WriteString(src.sb.String())
}

// this works too, because you are making a heap allocation (but this is kinda sad / Java-brain approach to doing things)
func (jp *JsonPointer) clone() *JsonPointer {
	var out JsonPointer
	s := jp.sb.String()
	out.sb.Grow(len(s))
	out.sb.WriteString(s)
	return &out
}

// other receiver methods were omitted for brevity

```

This is where "returning a struct value" bit me, because the `strings.Builder` has a pointer to a buffer internally, and
even though the original version of `JsonPointer` and `strings.Builder` will cease to exist, Golang devs made a
deliberate decision to cleverly detect that the internal buffer of the `strings.Builder` actually came from another
instance.

This is actually a good protective measure that prevents misuse and unexpected behaviour where one builder corrupts
the other. To illustrate how that would look like:

```text
JsonPointer -> strings.Builder -> *buf
(copy occurs)
JsonPointer (copy) -> strings.Builder (copy) -> *buf (same underlying memory, danger!)
(strings.Builder copy doesn't know that the old one is dead, therefore it will still panic if *buf was non-nil prior to copying)
```

A zero-initialized `strings.Builder` can be copied safely, but this is the by-product of a `nil` buffer -- if
the original has the buffer initialized with _anything_, your code will break

```go
package demo

import "strings"
import "testing"

func TestStringsBuilder_nonnilBufferCopy(t *testing.T) {
	var b1 strings.Builder
	b1.WriteString("hello") // the buffer isn't nil anymore

	b2 := b1 // Shallow copy: b2 shares b1’s buffer!

	b2.WriteString(" world") // panics!

	t.Log("b1:", b1.String()) // would print "hello world" had the program not crashed
	t.Log("b2:", b2.String()) // would also print "hello world"
}

func TestStringsBuilder_nilBufferCopy(t *testing.T) {
	var b1 strings.Builder

	b2 := b1 // yes, it's a copy but b1's buffer is nil / uninitialized

	b2.WriteString("World")  // b2's buffer is initialized, i.e. a nil pointer is set to a non-nil pointer
	b1.WriteString("Hello?") // the moment when b1's buffer gets initialized

	t.Log("b1:", b1.String()) // returns "Hello?"
	t.Log("b2:", b2.String()) // => returns "World"
}

```

### meta-notes:

- inaccuracy about zero-init structs in C: they are not always valid as they might need a malloc or two
- I made too many jumps between languages
- I focused too much on golang, therefore I must not pretend to speak about pointers/values in general, or to add more
  context around certain topics.