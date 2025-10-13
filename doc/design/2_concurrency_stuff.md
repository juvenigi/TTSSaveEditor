# Troubles with concurrency

Currently, some methods will benefit from parallelism
- even though file io is a bottleneck, shouting at your OS gives you the benefit of doing other useful things in the meantime
- once files are handled in goroutines, computing the sha3-256 checksum should be good enough as well

However, goroutines are like `Runnable`s, not promises, and they panic silently (unless something _terrible_ occurs)
- cancelling
- awaiting / sync with panic consideration -> either make channels deadlock by design
