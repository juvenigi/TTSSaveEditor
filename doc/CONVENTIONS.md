# Terminology

**AllowedInput** - possible

# Naming conventions

**NewFoo** - constructor, i.e. it guarantees a valid creation of a struct that does not have a nil-safe init value

- implies that the initial state is supplied as args, construction at most will only validate state coherence

**InitFoo** - initializer, more than just a constructor

- may contain size-effects for the purpose of retrieving the original state
- e.g. `GameCacheFinder` looks up the file directory during initialization, a constructor would have just accepted
  a map as input.

# Setter, Getter

- **SetFoo**, **GetFoo** : avoid primitive setters and getter like the plague
    - you are not in Java land, and things here are generally accepted as mutable. Be careful about what you pass!
    - if no sophisticated type conversion or cleanup happens, why not just use fields?
    - exposing an external 'setter' api is still best done via "Option structs", it's a much nicer pattern than the
      cringe
      of OOP brainrot

- **MutFoo** : the first argument is the object being mutated, the foo

- **FooRecur** : `FooRecur` calls itself, I recommend come with a prelude function **Foo**
    - however, if you hide away recursion, then you need to take the responsibility of making sure it's bounded for all
      possible input. Trust me, defensive programming is good programming. Redundant programming is not defensive, it's
      offensive as it's an insult to the reviewer who has to parse the code.