# GeeCache: A High-Performance Caching System

A standalone and **HTTP-based** distributed caching system, utilizing Go for efficient cache management and network communication.

- Implemented the **Least Recently Used (LRU)** algorithm to optimize cache storage by automatically discarding the least accessed items, enhancing system performance and resource utilization.
- Engineered a robust **lock mechanism in Go** to safeguard against cache penetration, significantly increasing system reliability and stability under high-load conditions.
- Applied **consistent hashing** for node selection within the distributed system, ensuring effective load balancing and improving scalability and fault tolerance.
- Integrated **Protocol Buffers (protobuf)** for node communication, optimizing binary data exchange to reduce latency and bandwidth usage, resulting in faster response times and improved overall efficiency.

For more functions and project details, please see [description](/description.md).