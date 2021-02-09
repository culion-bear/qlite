![](./qlite.png)
# QLite [[English Document](./README_EN.md)]

![](https://img.shields.io/badge/go-V1.14.3-brightgreen.svg)
![](https://img.shields.io/badge/release-V2.2.1-blue.svg)

QLite 是基于微服务的 NewSQL 型数据库系统，与传统的一体化数据库不同，该系统将本该内置的多种数据结构（STL）拆分成多个服务模块，每个模块都是独立的一个节点，每个节点都与其主网关进行连接，从而形成分布式存储结构。

QLite 主要存储方式为KV存储，主网关内置Hash结构，如同关系型数据报表，可存储已加载（连接）至主网关的所有数据结构。
> Hash结构可以存储Hash本身。这就导致QLite 数据库可以像文件夹一样进行树状查找及存储操作。

该系统配有轻量级 QLite - STL 开发框架，可以将自己编写的数据结构与QLite进行连接，具有高扩展性，但目前该框架只支持Go语言，具体见 [go-qlite](https://www.github.com/culion-bear/go-qlite)。

## Background

QLite 在发布至GitHub前已经经历了两个版本，分别是TCP版和集成版，效果皆不是很理想，该版本为分布式版本，采用HTTP作为接口交互协议，将信息处理完全交给了应用层，在损失了一定的性能的前提下提高了扩展性和复用性，更方便从业人员进行数据处理。

## Change Log

### V 2.2.1
- 重大更新！
- 现在添加服务节点时需要服务密码了，让服务变得更安全！至于密码是多少，得问服务节点部署者:)
- 更新了服务节点开发框架，现在多个网关能连接同一个服务节点而不产生混乱了
- 优化了join请求
- 优化了启动时添加节点的操作

### V 2.1.2
- 优化了服务连接模块
- 当网关关闭时会自动清空其他服务模块的缓存数据，防止再次开启时垃圾数据过多的情况

### V 2.1.1
- 解决了服务模块重连失败但未报错的BUG
- 优化了网关启动时服务连接操作，将连接失败时是否重连的选择权交给了使用者，而不是直接取消启动

### V 2.0.2
- 优化了优雅关闭服务（Ctrl+C）,解决关闭服务时数据不能及时存储至本地的BUG
- 优化了版本信息

## Feature

- 插件化集成：一种数据结构为一个模块，可随意注册至网关
- 提供Go语言版本的STL开发框架，可进行自定义开发
- 将Hash结构作为网关内置结构，等同于操作系统中的文件夹
- 存储模式选用KV型存储方案，在Hash结构的支持下，使数据间的关系不再混乱。

## Architecture

![](./architecture.png)

## Install

[linux-amd64-latest](https://github.com/culion-bear/qlite/releases/download/v2.2.1/qlite-linux-amd64)

[linux-arm64-latest](https://github.com/culion-bear/qlite/releases/download/v2.2.1/qlite-linux-arm64)

[windows-2.0.1-BETA](https://github.com/culion-bear/qlite/releases/download/v2.0.1-beta/qlite-windows.exe)

## Usage

```shell script
windows
qlite.exe -f [yaml path]
```

```shell script
linux
chmod 777 qlite
./qlite -f [yaml path]
```

[[click to download yaml file](./qlite.yaml)]

## Related Efforts

- STL
    - [qlite-string](https://github.com/culion-bear/qlite-stl-string)

- Other
    - [go-qlite](https://github.com/culion-bear/go-qlite)

## 快速开始

[API文档](./doc/api.md)

## 未来计划

- 产品
    - VLite(QLite可视化解决方案，正在测试中)
    - BLite(QLite中心化分布式集群解决方案，正在商议中)
    - QLite-Cli(QLite命令行客户端)
    - GLite(Go语言版QLite-ORM)
    - XLite(其他语言的ORM)
- 优化/未完成的功能
    - [√] 服务节点的多模块划分功能，主网关将携带密匙访问服务，这样多网关就能访问同一个服务节点而不造成数据混乱
    - [√] 加密的服务节点，在主网关join时初始化密码，在每次访问时携带，及Token
    - 具有守护进程的QLite
    - 子进程启动，关闭，重启父进程
- 其他
    - Docker一键部署

## 作者有话说

这个项目目前处于测试阶段，BUG肯定是有的，优化的方面也肯定有很多，希望感兴趣的兄弟们可以一起完善这个项目。如果您有更好的点子可以反馈给我，我会及时查阅的。