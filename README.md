## 数据结构(data structure)
1. 字符串 string
2. 双端链表 container/list
3. 字典 map, concurrent/map
4. 跳跃表 ryszard/goskiplist
5. HyperLogLog

## 内存编码(encoding)
1. intset
2. ziplist

## 数据类型(type)
1. 字符串键 string
2. 列表键 list
3. 散列键 hash
4. 集合键 set
5. 有序集合键 zset
6. HyperLogLog

## 数据库实现
* redisDb
* notify
* rdb
* aof

额外:
* pubsub
* multi
* sort
* bitops

## 客户端 & 服务端
* ae 事件处理器实现（基于 Reactor 模式）
* networking 网络连接库
* 单机 Redis 服务器的实现

## 多机功能
* 复制 replication
* sentinel
* cluster

## 初始化
1. default settings
2. flag 命令行
3. config file 配置文件

## Other ?
* tls
* 日志 sirupsen/logrus or uber-go/zap
* ACL User id
*

## feature
* goroutine 并发
* new cmds, eg. LOCK, UNLOCK, SETJSON 

## 论文
第一章	绪论
1 引言
2 论文组织结构

第二章	系统相关技术概述
1 Golang简介
2 raft简介

第三章	系统的需求分析
3.1 系统需求概述
3.2 系统功能需求分析
	3.2.1 数据结构、类型，存取
	3.2.2 日志、授权等次要功能
	3.2.3 持久化
	3.2.4 分布式强一致性
3.3 系统非功能需求
	3.3.1 用户界面——命令行需求
	3.3.2 性能需求
	3.3.3 软硬件环境需求
	3.3.4 产品质量需求

第四章	系统的设计
4.1 系统总体设计 （cli，srv，benchmark）
4.2 功能模块
4.3 系统详细设计

第五章	系统实现
5.1 开发环境与系统配置
5.2 数据结构、类型，存取
5.3 日志、授权等次要功能
5.4 持久化
5.5 分布式强一致性

第六章 系统测试
6.1 功能测试
6.2 性能测试

第七章 总结与展望
7.1 项目总结
7.2 项目展望

参考文献



