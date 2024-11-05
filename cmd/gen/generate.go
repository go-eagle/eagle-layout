package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"gorm.io/gen"

	"github.com/go-eagle/eagle-layout/internal/dal"
	"github.com/go-eagle/eagle-layout/internal/dal/db/method"
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	v "github.com/go-eagle/eagle/pkg/version"
)

var (
	cfgDir  = pflag.StringP("config dir", "c", "config", "config path.")
	env     = pflag.StringP("env name", "e", "dev", "env var name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	if *version {
		ver := v.Get()
		marshaled, err := json.MarshalIndent(&ver, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	// init config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	dal.Init()

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/dal/db/dao",                       // 生成的查询类代码的输出路径，默认./query, 这里改为dao
		ModelPkgPath: "./internal/dal/db/model",                     // 默认情况下会跟随OutPath参数，在同目录下生成model目录
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		WithUnitTest: true,                                          // 为生成的查询类代码生成单元测试
	})

	g.UseDB(dal.DB) // reuse your gorm db

	// 可直接指定结构体
	// Generate basic type-safe DAO API for struct `model.User` following conventions
	// g.ApplyBasic(model.UserModel{})

	// 指定要同步的表名
	// g.ApplyBasic(g.GenerateModel("user_info"))

	// 指定别名 {table_name} -> {model_name}
	// g.ApplyBasic(g.GenerateModelAs("user_info", "UserInfoModel"))

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.UserInfo`
	g.ApplyInterface(func(method.Querier) {},
		g.GenerateModelAs("user_info", "UserInfoModel",
			gen.FieldType("gender", "int32"),
			gen.FieldType("status", "int32"),
		),
	)

	// Generate the code
	g.Execute()
}
