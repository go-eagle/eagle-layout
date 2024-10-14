package main

import (
	"github.com/go-eagle/eagle-layout/internal/dal/method"
	"github.com/go-eagle/eagle-layout/internal/dal/model"
	"gorm.io/gen"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../../internal/dal/query",
		ModelPkgPath: "../../internal/dal/model",                    // 默认情况下会跟随OutPath参数，在同目录下生成model目录
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(model.DB) // reuse your gorm db

	// 可直接指定结构体
	// Generate basic type-safe DAO API for struct `model.User` following conventions
	// g.ApplyBasic(model.UserModel{})

	// 指定要同步的表名
	// g.ApplyBasic(g.GenerateModel("user_info"))

	// 指定别名 {table_name} -> {model_name}
	g.ApplyBasic(g.GenerateModelAs("user_info", "UserInfoModel"))

	// s
	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User`
	g.ApplyInterface(func(method.Querier) {})

	// Generate the code
	g.Execute()
}
