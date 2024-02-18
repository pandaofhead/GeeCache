# GeeCache: A High-Performance Caching System

A standalone and **HTTP-based** distributed caching system, utilizing Go for efficient cache management and network communication.

- Implemented the **Least Recently Used (LRU)** algorithm to optimize cache storage by automatically discarding the least accessed items, enhancing system performance and resource utilization.
- Engineered a robust **lock mechanism in Go** to safeguard against cache penetration, significantly increasing system reliability and stability under high-load conditions.
- Applied consistent hashing for node selection within the distributed system, ensuring effective load balancing and improving scalability and fault tolerance.
- Integrated **Protocol Buffers (protobuf)** for node communication, optimizing binary data exchange to reduce latency and bandwidth usage, resulting in faster response times and improved overall efficiency.

# Structure Tree
```
gee-cache
	|--go.mod
	|--main.go
	|--geecache/  
		|--lru/
			|--lru.go // lru 缓存淘汰策略
			|--lru_test.go  
		|--byteview.go // 缓存值的抽象与封装
		|--cache.go    // 并发控制
		|--geecache.go	// 负责与外部交互，控制缓存存储和获取的主流程
		|--geecache_test.go 
		|--http.go     // 提供被其他节点访问的能力(基于http)
```
## Caching

- FIFO
First In First Out: dosen't take frequency into account.

- LFU
Least Frequently Used: only take frequency into account.
- LRU
Least Recently Used: if some data has been used, then move it to the end of list, and the head of list will be the least recently used data, hence detele it.

## LRU
![LRU](/public/lru.jpg)

**sync.Mutex:**  
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

## HTTP in Go
> `http.ListenAndServer` takes two arguments, the first is address of service, the second is `Handler`.

Usage of standard HTTP module:
```go
package main

import (
	"log"
	"net/http"
)

type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!"))
}

func main() {
	var s server
	http.ListenAndServe("localhost:9999", &s)
}
```


