Table of Contents
- [Project Structure Tree](#project-structure-tree)
- [Caching](#caching)
- [HTTP in Go](#http-in-go)
- [Consistent Hashing](#consistent-hashing)
- [Distributed Nodes](#distributed-nodes)
- [Single Flight](#single-flight)
- [Functions and Packages](#functions-and-packages)

# Project Structure Tree
```
gee-cache
	|--README.md
	|--public
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
		|--go.mod // dependency mamagement
```
## Caching

- **FIFO**  
First In First Out: dosen't take frequency into account.

- **LFU**  
Least Frequently Used: only take frequency into account.

- **LRU**  
Least Recently Used: if some data has been used, then move it to the end of list, and the head of list will be the least recently used data, hence delete it.

## LRU
![LRU](/public/lru.jpg)
1. We first struct Cache() using a map: cache and a double-linked list:ll to store cache, also struct entry represents cache entry in list, interface Value could take Len() method, then New() cache.

2. Implement Get()
- if element in cache, move it to the front and return its value
- else just return
3. Implement Add()
- if element in cache, move it to the front and update its value in cache
- update new nbytes
- else put it at the front and add in cache
- update nbytes
- remove the oldest cache if exceeds maxBytes
4. Implement RemoveOldest()
- move the least visited node(at the back) in list and in cache map.
- update nbytes length

In addition to main component of LRU, we need to add mutex to support goroutine

1. struct ByteView to compute bytes(return length, slice) in cache.
2. now we add and get cache
- struct cache(mu, lru, cacheBytes)
- add cache using mu.Lock() and defer mu.Unock()
- get cache using mu.Lock() and defer mu.Unock()
3. Group is the main data struture that controls cache stream and interacts with users

### [map](https://go.dev/blog/maps)
```go
m = make(map[string]int)

i, ok := m["route"] // i is value, ok is whether key exists

_, ok := m["route"] // test for a key without retrieving the value

```

### [sync.Mutex](https://pkg.go.dev/sync#Mutex) 
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
### [sync.RWMutex](https://pkg.go.dev/sync#RWMutex)
	A RWMutex is a reader/writer mutual exclusion lock.	
1. Lock()
2. UnLock()
3. Rlock()
4. RUnlock()

## HTTP Server
Distributed caching needs HTTP to communicate between nodes.


### [http](https://pkg.go.dev/net/http#hdr-HTTP_2)
`http.ListenAndServer` takes two arguments, the first is address of service, the second is `Handler`.

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
### [strings.SplitN(s, sep string, n int)](https://pkg.go.dev/strings#SplitN)
slices s into substrings separated by sep and returns a slice of the substrings between those separators.
```go
func main() {
	fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 2))
	z := strings.SplitN("a,b,c", ",", 0)
	fmt.Printf("%q (nil = %v)\n", z, z == nil)
}
// ["a" "b,c"]
// [] (nil = true)
```
### [variadic parameter](https://www.digitalocean.com/community/tutorials/how-to-use-variadic-functions-in-go)
A variadic function(`...`) is a function that accepts zero, one, or more values as a single argument. The most common Println is a variadic function.
```go
func Println(a ...interface{}) (n int, err error)

func Printf(format string, a ...any) (n int, err error)
```
## [Consistent Hashing](https://www.xiaolincoding.com/os/8_network_system/hash.html)
**Hashing**  
To map the same key to the same node, the simplest method of hash algorithms is modulo operations. For example, in a distributed system with 3 nodes, data is mapped based on the formula hash(key) % 3.

However, there's a fetal problem: one change in node could result in Cache Avalanche.

**Consistent Hashing**  
Consistent hashing involves two steps:

- The first step is to perform a hash calculation on the storage nodes, that is, to perform a hash mapping of the storage nodes, such as hashing based on the node's IP address.
- The second step is to perform a hash mapping of the data when storing or accessing the data.

**Virtual Nodes**  
Instead of mapping the real nodes onto the hash ring, virtual nodes are mapped onto the hash ring, and these virtual nodes are then mapped to the actual nodes. 
```
			-> vitual node 
real node   -> vitual node -> hash ring
			-> vitual node 
```

## Distributed Nodes


## Single Flight

- **Cache Avalanche:** cache server crashes or setting the same expiration time for cached keys, causing a sudden increase in database request volume and pressure, leading to an avalanche.

- **Cache Penetration:** When a **key that exists expires**, and simultaneously, a large number of requests occur, these requests will **penetrate through to the database**, causing a sudden increase in database request volume and pressure.

- **Cache Piercing:** **Querying for data that does not exist.** Since it does not exist, it will not be written to the cache, so every request will go to the database. If there is a large amount of traffic in a short time, it can penetrate through to the database, leading to a crash.

