# worker
- 从etcd获取job并同步到内存中
- 实现调度模块，基于cron expression调度N个job
- 实现执行模块，并发的执行多个job
- 对job的分布式锁，防止集群并发
- 把执行日志保存到mongodb
