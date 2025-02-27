# Models

数据库使用的两种方式

## 方式一: 使用 `ExecTx` 和 `DBTx`

关于事务的使用，在 service 层开始事务，将待执行的 repo 方法封装在 fn 参数中，  
传递给 gorm 实例的 Transaction() 方法待执行。

同时在 Transaction 内部，触发 `fn()` 函数，也就是聚合的N个 repo 操作，  
需要注意的是，此时 `contextTxKey` 将携带事务 tx 的 ctx 作为参数传递给了  
`fn` 函数，因此下游的两个 repo 可以获取到 service 层的事务会话。

```go
// geeeter_svc.go
// CreateGreeter creates a Greeter, and returns the new Greeter.
func (s *GreeterService) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
   s.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
   var (
      greater *Greeter
      err     error
   )
   // 开始使用事务
   err = s.db.ExecTx(ctx, func(ctx context.Context) error {
      // 更新所有 hello 为 hello + "updated"，且插入新的 hello
      greater, err = s.repo.Save(ctx, g)
      _, err = s.repo.Update(ctx, g)
      return err
   })
   if err != nil {
      return nil, err
   }

   return greater, nil
}
```

在 repo 层执行数据库操作的时候，尝试通过 DBTx() 方法，从 ctx 中获取到上游传递  
下来的事务会话，如果有则使用，如果没有，则使用 repo 层自己持有的 DB，进行数据访问操作。

举例如下：

```go
// geeeter_repo.go
func (r *greeterRepo) Save(ctx context.Context, g *model.Greeter) (*model.Greeter, error) {
    result := r.db.DBTx(ctx).Create(g)
    return g, result.Error
}
​
func (r *greeterRepo) Update(ctx context.Context, g *model.Greeter) (*model.Greeter, error) {
    result := r.db.DBTx(ctx).Model(&model.Greeter{}).Where("hello = ?", g.Hello).Update("hello", g.Hello+"updated")
    if result.RowsAffected == 0 {
        return nil, fmt.Errorf("greeter %s not found", g.Hello)
    }
    return nil, fmt.Errorf("custom error")
}
```

## 方式二: 使用 `gorm/gen` 事务

gorm/gen 生成的代码中，没有使用事务，需要手动添加事务支持。 

```go
import (
	"context"

	"github.com/go-eagle/eagle-layout/internal/dal"
	"github.com/go-eagle/eagle-layout/internal/dal/db/dao"
)

func doSomething(ctx context.Context, id int64, data model.UserInfoModel) error {
    ...
	query := dao.Use(dal.DB)

    // 开始事务
    // 在事务内部，使用 gen 生成的 Query 进行多表操作
	query.Transaction(func(tx *dao.Query) error {
		_, err = tx.UserInfoModel.WithContext(ctx).Where(tx.UserInfoModel.ID.Eq(id)).Updates(data)
		if err != nil {
			return err
		}
		// 其他db操作
		// ...
		return nil
	})
    ...
}
```

## References

- https://juejin.cn/post/7399984522094149659?share_token=B3E5040F-3BC1-481A-A700-AFCF37F124BC
