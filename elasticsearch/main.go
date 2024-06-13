package main

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Movie struct {
	Title  string   `json:"title,omitempty"`
	Post   string   `json:"post,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Desc   string   `json:"desc,omitempty"`
	Source string   `json:"source,omitempty"`
}

func main() {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.56.101:9200",
		},
	}

	// 创建客户端连接
	typedClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
		return
	}

	err = count(typedClient)
	if err != nil {
		fmt.Println("failed to count, err:", err)
		return
	}

	// err = createIndex(typedClient)
	// if err != nil {
	// 	fmt.Println("failed to create index")
	// 	return
	// }

	// err = deleteIndex(typedClient)
	// if err != nil {
	// 	fmt.Println("failed to delete index")
	// 	return
	// }

	// err = getIndices(typedClient)
	// if err != nil {
	// 	fmt.Println("failed to get index, err:", err)
	// 	return
	// }

	// err = create(typedClient)
	// if err != nil {
	// 	fmt.Println("failed to create, err:", err)
	// 	return
	// }

	// err = search(typedClient)
	// if err != nil {
	// 	fmt.Println("failed to search, err:", err)
	// 	return
	// }

	// err = update(typedClient)
	// if err != nil {
	// 	fmt.Println("failed to update, err:", err)
	// 	return
	// }

	// err = delete(typedClient)
	// if err != nil {
	// 	fmt.Println("failed to delete, err:", err)
	// 	return
	// }
}

func create(client *elasticsearch.TypedClient) error {
	m := Movie{
		Title:  "title",
		Post:   "post",
		Tags:   []string{"tag1", "tag2"},
		Desc:   "desc",
		Source: "source",
	}
	resp, err := client.Index("movie").Document(m).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("result:", resp.Result)
	return nil
}

func search(client *elasticsearch.TypedClient) error {
	resp, err := client.Search().Index("movie").Query(&types.Query{
		Match: map[string]types.MatchQuery{"tags": {Query: "tag1"}},
	}).Do(context.Background())
	if err != nil {
		return err
	}
	for _, hit := range resp.Hits.Hits {
		fmt.Println(hit.Source_)
	}
	return nil
}

func update(client *elasticsearch.TypedClient) error {
	resp, err := client.UpdateByQuery("movie").
		Query(&types.Query{MatchPhrase: map[string]types.MatchPhraseQuery{"source": {Query: "souce"}}}).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func count(client *elasticsearch.TypedClient) error {
	resp, err := client.Count().Index("movie").
		Query(&types.Query{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"title": {Query: "fastx"},
			},
		}).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(resp.Count)
	return nil
}

func delete(client *elasticsearch.TypedClient) error {
	resp, err := client.DeleteByQuery("movie").
		Query(&types.Query{Match: map[string]types.MatchQuery{"tags": {Query: "tag2"}}}).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func createIndex(client *elasticsearch.TypedClient) error {
	resp, err := client.Indices.Create("movie").Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(resp.Index)
	return nil
}

func deleteIndex(client *elasticsearch.TypedClient) error {
	resp, err := client.Indices.Delete("movie").Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(resp.Acknowledged)
	return nil
}

func getIndices(client *elasticsearch.TypedClient) error {
	indices, err := client.Cat.Indices().Do(context.Background())
	if err != nil {
		return err
	}
	for _, index := range indices {
		fmt.Println(*index.Index)
	}
	return nil
}
