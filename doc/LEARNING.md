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
  - The use of `context.Context` is similar to `InterruptedException`, except that context is a form of function coloring.
  - Yet this coloring is less scary, because barely anything changes if a function simply passes context without doing
    anything with it.
    - adding context into a function is literally go's version of getter setter POJO boilerplate
- you can use callbacks / continuation-style passing 
  - or even spinlocks, if you're an engineer of your own destruction.

## As it turns out, Wails handles panics

But this means that I need to listen to its events, otherwise I'll miss out on the knowledge that something went wrong
This also puts some complexity on presenting a coherent state -- I need to invalidate all transient state after a panic.
