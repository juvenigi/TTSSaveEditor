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