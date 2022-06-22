### TODO: implement simple memory cache

- data structure-> key - byte data
  for example, you have,
	- key: 1232, some data
	- key: 1282, some data
- API should include following:
	- add (equivalent to upsert)
	- search
	- evict
	- cache (save current state to file)
- Eviction Rule
	- LRU
- TTL Rule
	- PERMA, or unixTimeMilli
- Some restriction
	- caching size should be less or equal to 1GB
	- connection should be protected by ID/PW, provided in ENV
- Some decisions
    - use rpc to connect this

## Gcache
Gcache is simple key-value store cache built with Go.


## Main Mechanism

### LRU
- Policy
    - Least Recently Used Node would be evicted.
    - Max node count will be provieded by ENV
- Data Structure
    - Doubly LinkedList
### TTL
- Policy
    - 0: permanent node. It willnot be removed by TTL Policy
    - timestamp: each node except perma will have ttl value
    - if hit, ttl will be extended
- Active TTL
    - Pros: Optimize Space, since it actively removes expired one.
    - Cons: Periodically locks other operation since it should clearout expired node.
- Passive TTL
    - Pros: No locking. If hitted object is expired, it just deletes it
    - Cons: Space will be wasted

## I/O
- RPC
    - Put(key, data, ttl?)
    - Get(key)
    - Delete(key)
    - Clear()
- Synchronization
    - Accessing Btree and LL should be treated as critical section
    - each RPC would create goroutine
    - each RPC should be queued and processed with order