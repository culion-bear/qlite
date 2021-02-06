![](./qlite.png)
# QLite [[中文文档](./README.md)]
> from YouDao translate

QLite is a NewSQL database system based on micro-services. Different from the traditional integrated database, this system divides the built-in multiple data structures (STL) into multiple service modules. Each module is an independent node, and each node is connected with its main gateway, thus forming a distributed storage structure.

The main storage mode of QLite is KV storage. The main gateway has a built-in Hash structure, like relational data reports, which can store all the data structures that have been loaded (connected) to the main gateway.
> A Hash structure can store the Hash itself. This results in a QLite database that can be looked up and stored as a tree like a folder.

The system is equipped with lightweight QLite-STL development framework, you can write their own data structure and QLite connection, with high scalability, but currently the framework only supports GO language, see the specific [go-qlite](https://www.github.com/culion-bear/go-qlite)。

## Background

QLite before release to the lot has experienced two versions, respectively is TCP version and integration edition, the effect is not very ideal, this version is distributed version, using HTTP as the interface communication protocol, the information processing to the application layer, completely lost in the premise of certain performance improves the extensibility and reusability, more convenient from personnel of course of data processing.

## Change Log

### V 2.1.1

- Resolved the BUG where the service module reconnect failed but no error was reported

- Optimized service connection operation at gateway startup, giving consumers the option to reconnect if the connection fails rather than simply cancel the startup

### V 2.0.2
- Optimized graceful service shutdown (CTRL +C), solved the BUG that data could not be stored locally in time when service shutdown

- Optimized version information

## Feature

- Plug-in integration: A data structure is a module that can be registered to a gateway at will

- Provide GO language version of STL development framework, can be customized development

- Use the Hash structure as the gateway built-in structure, equivalent to folders in the operating system

- KV storage scheme is selected as the storage mode. With the support of Hash structure, the relationship between data is no longer chaotic.

## Architecture

![](./architecture.png)

## Install

[linux-amd64-latest](https://github.com/culion-bear/qlite/releases/download/v2.0.2/qlite-linux-amd64)

[linux-arm64-latest](https://github.com/culion-bear/qlite/releases/download/v2.0.2/qlite-linux-arm64)

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

## How to Start

[API Document](./doc/api.md)

## Other

This project is currently in the testing stage, there must be some bugs, and there must be a lot of optimization aspects, I hope interested brothers can improve this project together. If you have better ideas to give me feedback, I will check them in time.