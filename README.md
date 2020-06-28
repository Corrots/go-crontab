# go-crontab
Crontab service of golang

## etcd

### Install
```bash
# Linux
ETCD_VER=v3.3.20

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}

rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

/tmp/etcd-download-test/etcd --version
ETCDCTL_API=3 /tmp/etcd-download-test/etcdctl version

cd /tmp/etcd-download-test
nohup ./etcd --listen-client-urls 'http://0.0.0.0:2379' --advertise-client-urls 'http://0.0.0.0:2379' &

#
less nohup.out
```
### etcdctl
```bash
cd /tmp/etcd-download-test
# get
ETCDCTL_API=3 ./etcdctl get /cron/lock/ --prefix
# watch
ETCDCTL_API=3 ./etcdctl get /cron/lock/ --prefix
# delete
ETCDCTL_API=3 ./etcdctl del /cron/lock/ --prefix
```
## worker flow
1. 从etcd从取出task存入内存中
2. 实现调度模块，基于cronexpr调度多个task
3. 实现执行模块，并发的执行多个task
4. 通过分布式锁来限制task在集群中的并发
5. 将执行日志保存至mongodb