# CacheServer

## 简介

分布式缓存服务器

## 依赖库安装

### 安装snappy

```bash
sudo apt-get install libsnappy-dev
```

### 安装zlib

``` bash
sudo apt-get install zlib1g-dev
```

### 安装bzip2

``` bash
sudo apt-get install libbz2-dev
```

### 安装lz4

``` bash
sudo apt-get install liblz4-dev
```

### 安装zstandard

``` bash
sudo apt-get install libzstd-dev
```

### clone子模块

```bash
git submodule update --init
```

### 编译rocksdb

``` bash
cd rocksdb && make static_lib
```

## 参数

```bash
--type 数据储存方式，包括inmemory、rocksdb、badger
--ttl 缓存失效时间单位s
--node 节点地址
--cluster 集群地址
```

### 单节点启动

```
$./server --node 1.1.1.1
```

### 集群启动

启动节点1

```
$./server --node 1.1.1.1
```

启动节点2

```
$./server --node 1.1.1.2 --cluster 1.1.1.1
```

