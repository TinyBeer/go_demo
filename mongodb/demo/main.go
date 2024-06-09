package main

import (
	"context"
	"encoding/json"
	"fmt"

	_ "compress/zlib"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
D：一个BSON文档。这种类型应该在顺序重要的情况下使用，比如MongoDB命令。
M：一张无序的map。它和D是一样的，只是它不保持顺序。
A：一个BSON数组。
E：D里面的一个元素。
*/

// Replace the placeholder with your Atlas connection string
const uri = "mongodb://192.168.56.101:27017"

func main() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	opts := options.Client().ApplyURI(uri)
	// authorization
	opts.SetAuth(options.Credential{
		// AuthMechanism:           "",
		// AuthMechanismProperties: map[string]string{},
		AuthSource: "admin",
		Username:   "zero",
		Password:   "123456",
		// PasswordSet: false,
	})
	// connection pool
	opts.SetMaxPoolSize(20)
	// compressors setting
	opts.SetCompressors([]string{"zlib"})
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Check connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	// fmt.Println("Ping result =", result)
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	coll := client.Database("test").Collection("user")
	name := "zhangsan"
	err = coll.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the name %s\n", name)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	// res, err := coll.InsertOne(context.TODO(), bson.D{{"name", "wangwu"}, {"age", 18}, {"habby", []string{"sleep", "football", "jogging"}}})
	// if mongo.IsDuplicateKeyError(err) {
	// 	fmt.Println("duplicated key")
	// 	return
	// }

	// jsData, _ := json.MarshalIndent(res, "", "  ")
	// fmt.Println(string(jsData))

	_, err = coll.DeleteOne(context.Background(), bson.D{{"name", "lisi"}})
	if err != nil {
		fmt.Println("something wrong", err)
		return
	}

	res, err := coll.UpdateOne(
		context.TODO(),
		bson.D{{Key: "name", Value: "wangwu"}},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "age", Value: 2}}},
			{Key: "$currentDate", Value: bson.D{
				{Key: "lastModified", Value: true},
			}},
		},
	)
	if err != nil {
		fmt.Println("something wrong", err)
		return
	}
	jsd, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(jsd))
}
