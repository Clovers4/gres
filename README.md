# 系统设计

## 数据结构与类型(data structure)
1. 字符串 string 对应 string
2. 双端链表 list 对应 clist
3. 字典 hash 对应 cmap
4. 集合 set 对应 cset
5. 有序集合键 zset 对应 ryszard/goskiplist

## 功能
* 持久化
* 日志

额外:
* pubsub
* multi

## 初始化配置来源
1. default settings
2. flag 命令行
3. config file 配置文件


## 系统目标
* 性能强大 （读写时不会锁整个db，持久化的时候也能渐进式地进行），利用goroutine协程以超越redis单线程性能瓶颈
* 有基本的kv存取，如，get/set。
* golang原生，丰富go生态，方便阅读研究，二次开发修改，最好可独立出内嵌型db，轻量，以再利用。
* 原理尽量简单，方便维护和理解。直接存储在内存中（类redis），因此需要实现lru逐出。内存hashmap还有一些好处，相比于b+tree,lsm等需要wal的/部分数据在磁盘上的，内存的反应一定更快。也显得更加轻量


提高方向（二选一？）：
1. 支持更多数据结构（仿照redis，相较之下仅有性能提升，且持久化实现难度提高），如，list，set，zset，hash。(涉及lpush/lpop,sadd,zadd,hset等相应功能)

2. 强一致性kv，并提供分布式功能，如，分布式锁trylock，分布式唯一id。（与“性能强大”可能互斥，etcd貌似用mvcc等技术实现高并发，但自己实现的话技术难度和复杂度会大大提高，并且由于使用了raft，基本上就是要实现成etcd这样，创新点又少了——>因为这些功能实际上可以用etcd 在proxy层或client层做，除非是要为了这些功能做特殊优化，难度大）
（raft库本身带有很多）




