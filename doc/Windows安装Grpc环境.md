# Windows安装Grpc环境向导

## 前言

[Proto3 Language Guide](https://protobuf.dev/programming-guides/proto3/)

## 搭建环境

### 准备环境

- Golang
- 配置$GOPATH环境变量

### 环境搭建目标

- protoc-23.2-win64.zip
- google.golang.org/grpc v1.55.0

## 正文

### 安装protoc编辑器

[protoc编辑器下载链接](https://objects.githubusercontent.com/github-production-release-asset-2e65be/23357588/77943df3-bf70-431a-b6b6-ed4fc72f815c?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20230603%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230603T035031Z&X-Amz-Expires=300&X-Amz-Signature=f9bc48cf9498a044199ee492425967a691842bbc565774fb6ca847e992cecdf0&X-Amz-SignedHeaders=host&actor_id=35629966&key_id=0&repo_id=23357588&response-content-disposition=attachment%3B%20filename%3Dprotoc-23.2-win64.zip&response-content-type=application%2Foctet-stream)

下载解压完成之后，将解压后的bin目录添加到系统环境变量的Path路径中，然后运行`protoc --version`是否配置成功，如下图所示，则安装成功。
```cmd
C:\Users\stone-98>protoc --version
libprotoc 23.2
```

### 安装protoc-gen-go

- 获取protoc-gen-go

```bash
service get -u google.golang.org/protobuf/cmd/protoc-gen-service
```

- 打包到GoPath路径

```bash
service install google.golang.org/protobuf/cmd/protoc-gen-service
```

- 在Gopath路径中将打包的bin目录拷贝到第一步解压的路径下，然后运行`1`测试是否配置成功，如下图所示，则配置成功。

```bash
C:\Users\stone-98>protoc-gen-service --version
protoc-gen-service v1.30.0
```

### 安装protoc-gen-go-grpc

- 获取protoc-gen-go

```bash
service get -u google.golang.org/grpc/cmd/protoc-gen-service-grpc
```

- 打包到GoPath路径

```bash
service install google.golang.org/grpc/cmd/protoc-gen-service-grpc
```

- 在Gopath路径中将打包的bin目录拷贝到第一步解压的路径下，然后运行`1`测试是否配置成功，如下图所示，则配置成功。

```bash
C:\Users\stone-98>protoc-gen-service-grpc --version
protoc-gen-service-grpc 1.3.0
```

## 测试

### 创建一个 proto 文件

```protobuf
syntax = "proto3";

option go_package = "proto/helloworld";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

### 生成proto文件

- 运行生成命令
```cmd
protoc --go_out=./grpc ./grpc/proto/*.proto
protoc --go-grpc_out=./grpc ./grpc/proto/*.proto
```
命令定义如下:
```cmd
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/*.proto
```
- SRC_DIR(应用程序源代码所在的目录——如果不提供值，则使用当前目录)，
- DST_DIR(生成的代码要去的目录;通常与$SRC_DIR相同)，以及.proto的路径。

执行完命令之后，在命令里指定的文件夹路径下将会生成对应的 helloworld.pb.go 文件