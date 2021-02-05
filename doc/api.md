# Api文档

该文档将系统的带您了解该数据库主要的（有可能是所有的）接口，让您更快的掌握如何进行交互。

## Table of Contents

- [前言](./api.md#前言)
- [用户操作](./api.md#用户操作)
    - [登录](./api.md#登录)
    - [ping](./api.md#ping)
- [数据库操作](./api.md#数据库操作)
    - [获取数据库信息](./api.md#获取数据库信息)
    - [获取数据结构(模板)列表](./api.md#获取数据结构(模板)列表)
    - [获取接口列表](./api.md#获取接口列表)
    - [添加一个服务模块](./api.md#添加一个服务模块)
    - [选择数据库](./api.md#选择数据库)
    - [清空数据库](./api.md#清空数据库)
    - [获取数据库个数](./api.md#获取数据库个数)
- [Hash操作](./api.md#Hash操作)
    - [添加数据](./api.md#添加数据)
    - [浏览数据](./api.md#浏览数据)
    - [获取元素个数](./api.md#获取元素个数)
    - [删除多个元素](./api.md#删除多个元素)
    - [获取元素类别](./api.md#获取元素类别)
    - [判断元素是否存在](./api.md#判断元素是否存在)
    - [设定时间](./api.md#设定时间)
    - [获取时间](./api.md#获取时间)
    - [重命名](./api.md#重命名)
- [其他](./api.md#其他)
    - [模板类API调用](./api.md#模板类API调用)
    - [算法类API调用](./api.md#算法类API调用)
    
    

## 前言
该系统的网络协议为HTTP协议，将信息处理完全交给了应用层，在损失了一定的性能的前提下提高了扩展性和复用性，更方便从业人员进行数据处理。

假设主机IP为127.0.0.1,端口号为5792

则客户端需访问该数据库时，应输入url为http://127.0.0.1:5792

但若开启跨域，则url应为http://127.0.0.1:5792/cors

另注意，除了登录（login）以及ping指令外，其余指令的url起始为http://127.0.0.1:5792/token或http://127.0.0.1:5792/cors/token并为请求携带token

token密钥携带方法为：HEAD[Authorization]="Bearer {token}"

如需POST访问，则文本协议均为JSON（**自定义服务模块除外**）

若服务端出现错误，则返回
- 响应体
    - code:int //非200值
    - msg:string //错误信息

## 用户操作
用户操作无需携带token，该操作是服务端对客户端进行验证处理

### 登录
- [POST]
- /login
- 请求体
    - password:string //登录密码
- 响应体
    - code:200
    - token:string//将返回的token放在请求头，已验证自身信息
### ping
- [POST]
- /ping
- 请求体
    - password:string //登录密码
- 响应体
    - code:200

## 数据库操作
url += /service

### 获取数据库信息
- [GET]
- /info
- 响应体
    - code:200
    - version:string//版本信息
### 获取数据结构(模板)列表
- [GET]
- /list/stl
- 响应体
    - code:200
    - list:[ ]
        - name:string //数据结构名
        - is_orderly:bool //是否为有序集合
### 获取接口列表
- [GET]
- /list/api/{stl:数据结构名}
- 响应体
    - code:200
    - list:[ ]
        - name:string //接口名
        - info:string //接口备注
### 添加一个服务模块
- [POST]
- /join
- 请求体
    - url:string //服务模块地址
- 响应体
    - code:200
### 选择数据库
- [GET]
- /database/{number:int}
- 响应体
    - code:200
    - token:string//将返回的token放在请求头，已验证自身信息
### 清空数据库
- [GET]
- /flush
- 响应体
    - code:200
- 作者注：删库跑路说的就是这个指令！慎用！
### 获取数据库个数
- [GET]
- /number
- 响应体
    - code:200
    - num:int//数据库个数

## Hash操作
url += /hash

由于Hash类型像操作系统中的文件夹，可以在Hash中放置新的Hash，以至于可以出现树状查询及存储操作，若要实现该功能，仅需在url末尾加入一段路径（path）

假如在数据库0中有类别为hash的键值k1,k1中存储类别为hash的键值k2,k2存储一个键值为k3的元素，若要获取到k3的类别，仅需POST该url：http://127.0.0.1:5792/token/hash/type/k1/k2,请求体中key为k3，详情请访问[type](./api.md#获取元素类别)

### 添加数据
- [POST]
- /set(添加类别为Hash的空键值)
- 请求体
    - key:string //键值
    - time:int //持续时间，若小于等于0则为永久
- 响应体
    - code:200
---
- [POST]
- /set_x(若key存在则覆盖)
- 请求体
    - key:string //键值
    - time:int //持续时间，若小于等于0则为永久
- 响应体
    - code:200
### 浏览数据
- [GET]
- /select
- 响应体
    - code:200
    - list:[]string //获取key的列表
---
- [GET]
- /select_x
- 响应体
    - code:200
    - list:[] //获取key的列表
        - key:string
        - type:string //类别
        - time:int
### 获取元素个数
- [GET]
- /size
- 响应体
    - code:200
    - num:int //列表中元素的个数
### 删除多个元素
- [POST]
- /delete
- 请求体
    - keys:[]string //key列表
- 响应体
    - code:200
    - num:int //返回成功删除的个数
### 获取元素类别
- [POST]
- /type
- 请求体
    - key:string
- 响应体
    - code:200
    - type:string //获取元素类别
### 判断元素是否存在
- [POST]
- /exists
- 请求体
    - key:string
- 响应体
    - code:200
    - exists:bool
### 设定时间
- [POST]
- /pex
- 请求体
    - key:string //键值
    - time:int //持续时间，若小于等于0则为永久
- 响应体
    - code:200
---
- [POST]
- /pex_to(至何时结束)
- 请求体
    - key:string //键值
    - time:int //秒级时间戳
- 响应体
    - code:200
### 获取时间
- [POST]
- /time(查询剩余时间)
- 请求体
    - key:string //键值
- 响应体
    - code:200
    - time:int
---
- [POST]
- /time_to(查询至何时结束)
- 请求体
    - key:string //键值
- 响应体
    - code:200
    - time:int //秒级时间戳
### 重命名
- [POST]
- /rename
- 请求体
    - key:string //键值
    - new_key:string //要改成的键值名
- 响应体
    - code:200
---
- [POST]
- /rename_x(若新键值名存在，则覆盖)
- 请求体
    - key:string //键值
    - new_key:string //要改成的键值名
- 响应体
    - code:200

## 其他
客户端通过以下方式跟其他数据结构进行对接，其数据结构内置接口如何使用，请联系其开发者

模板类API通常是对其单一数据进行增删改查，算法类API则是对多组数据进行算法调度，但对本身数据并无影响

### 模板类API调用
- [POST]
- /stl/{service}/{api}/{path}(path视情况而定，同hash，可省略)
- 请求体
    - key:string //键值
    - time:int //时间，可省略
    - opt:object //视第三方数据结构而定
- 响应体
    - ？ //视第三方数据结构而定
### 算法类API调用
- [POST]
- /alg/{service}/{api}
- 请求体
    - keys:[]
        - key:string
        - path:string
    - opt:object //视第三方数据结构而定
- 响应体
    - ？ //视第三方数据结构而定
