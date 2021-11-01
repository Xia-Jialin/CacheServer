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
