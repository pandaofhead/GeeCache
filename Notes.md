# Caching

## FIFO
First In First Out: dosen't take frequency into account.

## LFU
Least Frequently Used: only take frequency into account.
## LRU
Least Recently Used: if some data has been used, then move it to the end of list, and the head of list will be the least recently used data, hence detele it.

# LRU
![LRU](/public/lru.jpg)

# sync.Mutex: 
Mutexes only allow one goroutine to acquire the lock and access the shared resource, while other goroutines wait until the lock is released.
```go
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
    
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)

	}
	c.lru.Add(key, value)
}
```