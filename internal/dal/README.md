# DAL 数据访问层

数据访问层主要实现具体的数据访问逻辑，比如与数据库交互、cache、RPC操作等。

repository 层会调用 `dal` 层的方法来完成实际的数据操作，供上层 `service` 调用

> 注：禁止 `service` 层直接操作 `dal` 进行数据查询

## 目录结构

- `dal/db` 层按照表名划分文件，该层的逻辑尽量保持原子性，只是单纯的增删改查  
- `dal/db/model` 中只放 `struct` 定义和成员函数，其它函数不应该放在 `model` 中  
- `dal/db/method` 里放的是自定义的查询函数，是基于sql注解生成的  
- `dal/db/query` 是基于 `gorm/gen` 生成的原子性的查询方法，比如 `Find`, `First` 等

## 数据获取

- 通过 `dal/db/query` 或 `dal/db/method` 直接从数据库获取
- `dal/rpc` RPC访问接口
- `dal/cache` 获取redis中的数据
