# GoLightCache: A High-Performance Caching System
> ✨Inspired by [GeeCache](https://geektutu.com/post/geecache.html)✨

A standalone **Golang distributed caching system**, utilizing **gRPC and etcd** for efficient cache management and network communication.

- Implemented the **LRU and LFU** algorithm to optimize cache storage by automatically discarding the least accessed items, enhancing system performance and resource utilization.
- Engineered **gRPC-based** distributed cache to enable multiple nodes to work together, improving system scalability and fault tolerance.
- Applied **Consistent Hashing** for node selection within the distributed system, ensuring effective load balancing and improving scalability and fault tolerance.
- Integrated **Protocol Buffers (protobuf)** for node communication, optimizing binary data exchange to reduce latency and bandwidth usage, resulting in faster response times and improved overall efficiency.
- Utilized **etcd** for service registration and discovery, enabling nodes to automatically discover each other and work together, improving system scalability and fault tolerance.

✨Improvements on original design✨:
- Add **LFU** algorithm to cache
- Add **TTL** and **lazy delete** to cache
- Add **gRPC** to communicate between nodes
- Add **etcd** to register and discover nodes

# GoLightCache Workflow
![workflow](./public/golightcache.png)

# Structure Tree
```bash
│  go.mod
│  go.sum
│  main.go	
│  README.md	
│  run.sh	
│
└─geecache
    │  byteview.go	// cache abstraction layer
    │  cache.go	    // cocurrent safe cache
    │  geecache.go	负责与外部交互，控制缓存存储和获取的主流程
    │  geecache_test.go 			
    │  peers.go	// abstract PeerPicker
    │  grpc.go	// Server/Client for gRPC
    │
    ├─consistenthash
    │      consistenthash.go	
    │      consistenthash_test.go	
    │
    ├─geecachepb
    │      geecachepb.pb.go
    │      geecachepb.proto	
    │      geecachepb_grpc.pb.go
    │
    ├─lfu
    │      lfu.go	
    │      lfu_test.go
    │
    ├─lru
    │      lru.go	
    │      lru_test.go
    │
    ├─registry	
    │      discover.go	
    │      register.go	
    │
    └─singleflight
            singleflight.go	防止缓存击穿
            singleflight_test.go
```
## Install protoc
please see [ProtoUsage.md](./ProtoUsage.md)