# Repository

`repository` 层负责为上层（通常是服务层或业务逻辑层）提供统一的数据访问接口，  
它封装了数据的获取、存储和更新等操作，对上层屏蔽了底层数据存储的具体实现细节

调用链路如下

```bash
            service
              |
              |
            repository
              | (calls)
              |
               ------
              |      |      |
             dal
             /   |   \
            /    |    \
           /     |     \
        db   rpc   cache
```