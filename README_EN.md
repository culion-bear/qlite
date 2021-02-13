![](./qlite.png)
# QLite [[中文文档](./README.md)]
> from YouDao translate

![](https://img.shields.io/badge/go-V1.14.3-brightgreen.svg)
![](https://img.shields.io/badge/release-v2.2.5-blue.svg)

QLite is a NewSQL database system based on micro-services. Different from the traditional integrated database, this system divides the built-in multiple data structures (STL) into multiple service modules. Each module is an independent node, and each node is connected with its main gateway, thus forming a distributed storage structure.

The main storage mode of QLite is KV storage. The main gateway has a built-in Hash structure, like relational data reports, which can store all the data structures that have been loaded (connected) to the main gateway.
> A Hash structure can store the Hash itself. This results in a QLite database that can be looked up and stored as a tree like a folder.

The system is equipped with lightweight QLite-STL development framework, you can write their own data structure and QLite connection, with high scalability, but currently the framework only supports GO language, see the specific [go-qlite](https://www.github.com/culion-bear/go-qlite)。

## Background

QLite before release to the lot has experienced two versions, respectively is TCP version and integration edition, the effect is not very ideal, this version is distributed version, using HTTP as the interface communication protocol, the information processing to the application layer, completely lost in the premise of certain performance improves the extensibility and reusability, more convenient from personnel of course of data processing.

## Change Log

### V 2.2.5

- Fixed Bug

### V 2.2.4

- Last update before New Year

- Added daemons to start programs

- Added an operation to send a signal from a child process to the main process

- Added QLite initialization

### V 2.2.3

- Fixed Bug

### V 2.2.2

- The author is trying to update it

- A very serious BUG was found, when a service node crashes and restarts, when the gateway reconnects to the service node, all the data will be lost unless the gateway is restarted. Of course, the author found it and fixed it, but it still needs to be optimized.

### V 2.2.1
- Major updates!

- Service password is now required when adding service nodes, making services more secure! Ask the service node deployer what the password is :)

- Updated service node development framework so that multiple gateways can now connect to the same service node without confusion

- Optimized join requests

- Optimized the operation of adding nodes at startup

### V 2.1.2
- Optimized service connection module

- When the gateway is closed, it will automatically empty the cache data of other service modules to prevent excessive garbage data when it is opened again

## Feature

- Plug-in integration: A data structure is a module that can be registered to a gateway at will

- Provide GO language version of STL development framework, can be customized development

- Use the Hash structure as the gateway built-in structure, equivalent to folders in the operating system

- KV storage scheme is selected as the storage mode. With the support of Hash structure, the relationship between data is no longer chaotic.

## Architecture

![](./architecture.png)

## Install

[linux-amd64-latest](https://github.com/culion-bear/qlite/releases/download/v2.2.5/qlite-linux-amd64)

[linux-arm64-latest](https://github.com/culion-bear/qlite/releases/download/v2.2.5/qlite-linux-arm64)

## Usage

```shell script
[get version]
./qlite -v

[get help]
./qlite -h

[start qlite]
./qlite -f [yaml path]

[init qlite](run as root)
./qlite -i

[run as daemon]
./qlite -d

[stop qlite]
./qlite -s stop

[reload qlite]
./qlite -s reload
```

[[click to download yaml file](https://github.com/culion-bear/qlite/releases/download/v2.2.5/qlite.yaml)]

## Related Efforts

- STL
    - [qlite-string](https://github.com/culion-bear/qlite-stl-string)

- Other
    - [go-qlite](https://github.com/culion-bear/go-qlite)

## How to Start

[API Document](./doc/api.md)

## Future

- product

    - VLite(QLite visualization solution, under test)
    - Blite (QLite Centralized Distributed Cluster Solution, under negotiation)
    - QLite-Cli (QLite command line client)
    - Glite (Go QLite-ORM)
    -Xlite (ORM in other languages)
- Optimized/Unfinished features
    - [√]  Multi-module partitioning of service nodes. The main gateway will carry the key to access the service, so that multiple gateways can access the same service node without causing data clutter
    - [√]  Encrypted service node that initializes the password on the main gateway join, carries it on each access, and the Token
    - [√] QLite with daemon
    - [√] Child process start, shut down, restart parent process
    - [√] Data loss of built-in nodes when the service node crashes
- other
    - Docker one-click deployment

## Other

This project is currently in the testing stage, there must be some bugs, and there must be a lot of optimization aspects, I hope interested brothers can improve this project together. If you have better ideas to give me feedback, I will check them in time.