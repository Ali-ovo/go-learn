package main

import (
	"flag"
	"os"
	"os/exec"
	"path"
	"shop/tools/protocGenGin/generator"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func registerProtoFile(srcDir string, filename string) error {
	// 首先，将.proto文件转换为文件描述符集
	tmpFile := path.Join(srcDir, filename+".pb")
	cmd := exec.Command("protoc", // cmd 命令
		"--include_source_info",
		"--descriptor_set_out="+tmpFile,
		"--proto_path="+srcDir,
		path.Join(srcDir, filename))

	cmd.Stdout = os.Stdout // 设置 标准输出流
	cmd.Stderr = os.Stderr // 设置 标准输入流
	err := cmd.Run()       // 运行 cmd 命令
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpFile) // 关闭文件

	// 现在将该临时文件作为文件描述符集protobuf加载
	protoFile, err := os.ReadFile(tmpFile)
	if err != nil {
		panic(err)
	}

	pbSet := new(descriptorpb.FileDescriptorSet)
	if err := proto.Unmarshal(protoFile, pbSet); err != nil {
		panic(err)
	}

	// 我们知道protoc是通过一个.proto文件调用的
	pb := pbSet.GetFile()[0]

	// Initialize the file descriptor object.
	fd, err := protodesc.NewFile(pb, protoregistry.GlobalFiles)
	if err != nil {
		panic(err)
	}

	// and finally register it.
	return protoregistry.GlobalFiles.RegisterFile(fd) // 将文件的元数据信息注册到全局文件注册表中
}

func registerProto() {
	files := []string{
		"google/api/http.proto",
		"google/api/annotations.proto",
		"google/api/client.proto",
		"google/api/field_behavior.proto",
		"google/api/resource.proto",
	}
	for _, f := range files {
		err := registerProtoFile("./", f)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// protoc 会在一些默认路径查找所需要的 公共 proto 文件
	// 如果没查找到, 只能手动注入 以下函数就是手动注入代码
	// 也可以将 那些 公共 proto 文件放在 当前目录下 也可以 成功
	//registerProto()

	flag.Parse()
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generator.GenerateFile(gen, f)
		}
		return nil
	})
}
