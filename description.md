Table of Contents
- [Project Structure Tree](#project-structure-tree)
- [Caching](#caching)
- [HTTP Server](#http-server)
- [Consistent Hashing](#consistent-hashing)
- [Distributed Nodes](#distributed-nodes)
- [Single Flight](#single-flight)
- [Protocol Buffers](#protocol-buffers)

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
In a distributed system, data will be distributed into different nodes. To save space and time, we should try to map the same key to the same node, the simplest method of hashing is modulo operations. For example, there're 3 nodes, key is mapped based on the formula `hash(key) % 3`.

However, one big issue coming up: one change in a node could result in [Cache Avalanche](#single-flight) because the change of cache mapping.

**Consistent Hashing**  
Imagine a hash ring consists of 2^32 nodes:

- First, perform a hash calculation on the storage nodes, that is, to perform a hash mapping of the storage nodes, such as hashing based on the node's IP address.
- Then perform a hash mapping(clockwise) on the ring when storing or accessing the data.

However, yes another however, there's still a problem: nodes will be unevenly placed on the ring, causing many keys point to one node, which might also lead to [Cache Avalanche](#single-flight).

**Virtual Nodes**  
Virtual nodes are mapped onto the hash ring, and these virtual nodes are then mapped to the actual nodes. 
```
vitual node(A-1) 
vitual node(A-2) -> real node(A) -> hash ring
vitual node(A-3) 
```
Now even with changes in nodes, mutiple virtual nodes together will take in changes and increase stability.

### crc32.ChecksumIEEE
```go
func ChecksumIEEE(data []byte) uint32
```
The CRC-32 checksum is a type of hash function that generates a 32-bit (4-byte) hash value,

### [strconv.Itoa()](https://pkg.go.dev/strconv#Itoa)
int to string
```go
func Itoa(i int) string
```

### [strconv.Atoi()](https://pkg.go.dev/strconv#Itoa)
string to int
```go
func Atoi(s string) (int, error)
```

### sort.Ints()
Ints sorts a slice of ints in increasing order.
### [sort.Search()](https://pkg.go.dev/sort#example-Search)
Search uses **binary search** to find and return the smallest index i in [0, n) at which f(i) is true, assuming that on the range [0, n), f(i) == true implies f(i+1) == true. 
```go
func Search(n int, f func(int) bool) int

idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
})
```

## Distributed Nodes
```
                  1
key --> cached? -----> return cache(key)
                |  0                        	1
                |-----> get from remote nodes? -----> interacts with remote nodes --> return cache(key)
                            |  0
                            |-----> callback func gets value and add to cache --> return cache(key)
```

### [url.QueryEscape](https://pkg.go.dev/net/url#QueryEscape)
QueryEscape escapes the string so it can be safely placed inside a URL query.
```go
func main() {
	query := url.QueryEscape("my/cool+blog&about,stuff")
	fmt.Println(query)

}
// my%2Fcool%2Bblog%26about%2Cstuff
```
## Single Flight
### Cache Avalanche: 
cache server crashes or setting the same expiration time for cached keys, causing a sudden increase in database request volume and pressure, leading to an avalanche.

### Cache Penetration: 
When a **key that exists expires**, and simultaneously, a large number of requests occur, these requests will **penetrate through to the database**, causing a sudden increase in database request volume and pressure.

### Cache Piercing: 
**Querying for data that does not exist.** Since it does not exist, it will not be written to the cache, so every request will go to the database. If there is a large amount of traffic in a short time, it can penetrate through to the database, leading to a crash.

｜-------------------------------------------------｜  
To make sure that same keys call HTTP once such that no cache penetration, we use singleflight to protect database.

**Mutex and WaitGroup are important for the implenmentation.**
```go
	// if the key is already in-flight, wait for it
	if c, ok := g.m[key]; ok {
		g.mu.Unlock() // unlock before wg.Wait()
		c.wg.Wait()   // wait for the call to complete
		return c.val, c.err
	}

	// the key is not in-flight; make the fn call
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c

	// unlock and get key from remote
	g.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()
```
- if we can get the key from map, we unlock it and wait for it to finish fetching
- if the key has to be fetched from remote, we add one on wait, unlock and fetch, then lock it

```go
	// remove the key from the map
	// prevents memory leak and updates the map
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
```

### sync.WaitGroup
sync.WaitGroup拥有一个内部计数器，当计数器等于0时，Wait()方法会立即返回
```go
func (wg *WaitGroup) Add(delta int) // Add添加n个并发协程
func (wg *WaitGroup) Done()  		// Done完成一个并发协程
func (wg *WaitGroup) Wait()  		// Wait等待其它并发协程结束

package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    wg := &sync.WaitGroup{}
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(i int) { // Pass the loop variable as an argument to the goroutine
            defer wg.Done()
            time.Sleep(1 * time.Second)
            fmt.Printf("hello world ~ %d\n", i) // Optionally print the value of i
        }(i) // Pass the loop variable i to the goroutine
    }
    // Wait for all goroutines to finish
    wg.Wait()
    fmt.Println("WaitGroup all process done ~")
}
```
## [Protocol Buffers](https://protobuf.dev/)
**Protocol Buffers are language-neutral, platform-neutral extensible mechanisms for serializing structured data.**

|Protobuf|JSON/XML|
|---|---|
|binary|text-based|
|smaller|larger|
|separation of context and data||

> To install and start using Protocol Buffers, see [ProtoUsage.md](ProtoUsage.md)

