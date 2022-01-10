module github.com/arktos-venture/golang-template

go 1.17

replace github.com/arktos-venture/grpc-apis => ../grpc-apis

require (
	github.com/arktos-venture/grpc-apis v0.0.0-20211213233223-74a70d9e73ae
	github.com/go-kratos/kratos/v2 v2.1.2
	github.com/hashicorp/go-hclog v1.0.0
	github.com/spf13/pflag v1.0.5
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.6.2 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
