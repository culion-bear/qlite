![](./data/qlite.png)
# QLite 

![](https://img.shields.io/badge/go-V1.14.3-brightgreen.svg)
![](https://img.shields.io/badge/release-v4.0.1-blue.svg)
![](https://img.shields.io/badge/type-NoSQL-yellow.svg)

QLite是一款可插件化扩展的NoSQL型数据库系统，与其他的数据库不同，该系统可根据开发者需要来自定义插件扩展数据库的功能

## 背景

V2版本的QLite采用HTTP + GRPC作为开发方案，然而测试时发现QPS极其低下，故不得不考虑新的解决方案，在未发行的V3版本中采用TCP协议作为底层协议再进行二次封装，通过Go - plugins作为插件化解决方案，效果比较理想，但并未很好的处理细节问题，故进行重构并推出V4版本。

## 日志

### V 4.0.1
- 全新版本上线
- 更简洁的交互语句
- 更快的数据操作

## 亮点

- 插件化扩展：开发者可根据自需来自定义插件，进而扩展数据库的功能
- 高效又简洁的交互设计

## 作者有话说

这个项目目前处于测试阶段，BUG肯定是有的，优化的方面也肯定有很多，希望感兴趣的朋友们可以一起完善这个项目。如果您有更好的点子可以反馈给我，我会及时查阅的。
