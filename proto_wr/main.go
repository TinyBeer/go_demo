package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bufbuild/protocompile"
	"github.com/jhump/protoreflect/v2/protobuilder"
	"github.com/jhump/protoreflect/v2/protoprint"
)

func main() {
	errFilePath := "errors.proto"
	serverFilePath := "server.proto"
	targetFilePath := "bag.proto"

	// file := protobuilder.NewFile(fileName)
	// file.SetSyntax(protoreflect.Proto3)

	ctx := context.Background()
	compiler := &protocompile.Compiler{
		SourceInfoMode: protocompile.SourceInfoExtraComments,
		// 可以添加 import 搜索路径
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: []string{
				"/home/beer/workspace/sim/src/proto/def/inner/",
				"/home/beer/workspace/sim/src/proto/def/out/",
				"/home/beer/workspace/sim/src/services/bag/api_bag/",
			},
		}),
	}

	files, err := compiler.Compile(ctx, targetFilePath, serverFilePath, errFilePath)
	if err != nil {
		panic(err)
	}

	errFile := files.FindFileByPath(errFilePath)
	serverFile := files.FindFileByPath(serverFilePath)
	targetFile := files.FindFileByPath(targetFilePath)

	file, err := protobuilder.FromFile(targetFile)
	if err != nil {
		log.Fatal(err)
	}

	msg := file.GetMessage("HelloReq")
	if msg != nil {
		log.Fatal("message already exists !!!")
	}
	msg = protobuilder.NewMessage("HelloReq").
		AddField(protobuilder.NewField("role", protobuilder.FieldTypeImportedMessage(serverFile.Messages().ByName("RpcRoleInfo"))).SetNumber(1)).
		AddField(protobuilder.NewField("ErrCode", protobuilder.FieldTypeImportedEnum(errFile.Enums().ByName("ErrCode"))).SetNumber(2))
	file.AddMessage(msg)

	fd, err := file.Build()
	// fmt.Println(fd)
	fmt.Println(err)

	f, err := os.OpenFile("test.proto", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatal(err)
	}
	pr := protoprint.Printer{
		Compact:      true,
		SortElements: true,
	}
	err = pr.PrintProtoFile(fd, f)
	fmt.Println(err)
}

// // Parse element from proto file
// // 从proto文件中解析出各种元素
// func main() {
// 	ctx := context.Background()

// 	compiler := &protocompile.Compiler{
// 		// 可以添加 import 搜索路径
// 		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
// 			ImportPaths: []string{"/home/beer/workspace/sim/src/proto/def/out/"},
// 		}),
// 	}

// 	files, err := compiler.Compile(ctx, "lobby_mail.proto")
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, fd := range files {
// 		targetName := "MailData"
// 		if msg := fd.Messages().ByName(protoreflect.Name(targetName)); msg != nil {
// 			fmt.Printf("Message: %s\n", msg.Name())
// 			for fidx := range msg.Fields().Len() {
// 				field := msg.Fields().Get(fidx)
// 				fmt.Printf("  field %s = %d, type=%s\n", field.Name(), field.Number(), field.Kind())
// 			}
// 		} else {
// 			fmt.Printf("Message: %s not found\n", targetName)
// 		}
// 		fmt.Println()
// 		fmt.Printf("==== File: %s, package=%s ====\n", fd.Path(), fd.Package())

// 		// Messages
// 		for idx := range fd.Messages().Len() {
// 			msg := fd.Messages().Get(idx)
// 			fmt.Printf("Message: %s\n", msg.Name())
// 			for fidx := range msg.Fields().Len() {
// 				field := msg.Fields().Get(fidx)
// 				fmt.Printf("  field %s = %d, type=%s\n", field.Name(), field.Number(), field.Kind())
// 			}
// 		}

// 		// Services
// 		for sidx := range fd.Services().Len() {
// 			svc := fd.Services().Get(sidx)
// 			fmt.Printf("Service: %s\n", svc.Name())
// 			for mid := range svc.Methods().Len() {
// 				rpc := svc.Methods().Get(mid)
// 				fmt.Printf("  RPC %s: input=%s, output=%s\n", rpc.Name(), rpc.Input().Name(), rpc.Output().Name())
// 			}
// 		}

// 		// Enums
// 		for eid := range fd.Enums().Len() {
// 			enum := fd.Enums().Get(eid)
// 			fmt.Printf("Enum: %s\n", enum.Name())
// 			for vid := range enum.Values().Len() {
// 				ev := enum.Values().Get(vid)
// 				fmt.Printf("  %s = %d\n", ev.Name(), ev.Number())
// 			}
// 		}
// 	}
// }
