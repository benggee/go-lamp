## lamp脚手架工具

lamp 工具可以通过命令行生成一个即时可用的项目脚手架

### lamp 安装
``` 
git clone https://github.com/seepre/go-lamp.git
cd go-lamp/cli/lamp
make 
```

### lamp 使用说明

```lamp new -d ./myworkspace -name example-project```

    -d  脚手架要创建到哪个目录
    -name 脚手架名字

### 目录结构
```bigquery
├── go.mod
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── demo
│   │   │   └── demo.go
│   │   └── route.go
│   └── service
│       └── demo.go
└── main.go
```

### 运行
```bigquery
go run main.go
```