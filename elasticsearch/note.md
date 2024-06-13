ElasticSearch 是一个高度可扩展的开源实时搜索和分析引擎，它允许用户在近实时的时间内执行全文搜索、结构化搜索、聚合、过滤等功能。Elasticsearch 基于 Lucene 构建，提供了强大的全文搜索功能，并且具有广泛的应用领域，包括日志和实时分析、社交媒体、电子商务等。

# 环境搭建

- 创建 compose.yaml 文件

```yaml
version: '3'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.9.1
    environment:
      - ES_JAVA_OPTS=-Xms512m -Xmx512m
      - discovery.type=single-node
      - xpack.security.enabled=false
    networks:
      - elasticsearch
    ports:
      - 9200:9200
  kibana:
    image: docker.elastic.co/kibana/kibana:8.9.1
    networks:
      - elasticsearch
    ports:
      - 5601:5601
    environment:
      ELASTICSEARCH_HOSTS: '["http://elasticsearch:9200"]'
    depends_on:
      - elasticsearch
networks:
  elasticsearch:
    driver: bridge
```

- 使用 docker-compose 运行

```sh
docker-compose up -d
```

- 访问 http://127.0.0.1:5601
  可能需要等待一会，kibanna 准备需要一些时间。进入页面后，从左侧菜单栏进入 Manager-->DevTool 打开开发工具页面。后续我们就可以在页面左侧窗口中输入 curl 命令，点击 “▶” 符号发送请求后，在页面右侧窗口查看返回结果。

# API

## 查看健康状态

```
GET /_cat/health?v

epoch      timestamp cluster        status node.total node.data shards pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
1718103967 11:06:07  docker-cluster green           1         1      6   6    0    0        0             0                  -                100.0%

```

## 文档操作

### 创建文档

将 JSON 文档添加到指定的数据流或索引并使其可搜索。

- 如果索引不存在，则会创建默认配置的索引。索引相关内容想在后文详细介绍。
  我们已经通过索引一篇文档创建了一个新的索引 。这个索引采用的是默认的配置，新的字段通过动态映射的方式被添加到类型映射。
- 如果目标是索引并且文档已经存在，则请求更新文档并递增其版本。

```
PUT /<target>/_doc/<_id>

POST /<target>/_doc/

PUT /<target>/_create/<_id>

POST /<target>/_create/<_id>
```

```
POST /movie_index/_create/1
{
    "id":1,
    "title":"a movie",
    "post_url":"post url",
    "tags":["action","sci_fic"],
    "desc":"this is a movie",
    "source_url":"source url"
}

```

### 判断文档是否存在

```
HEAD /movie_index/_doc/1
```

如果存在，Elasticsearch 返回 200 - OK 的响应状态码，如果不存在则返回 404 - Not Found。

### 获取文档

```
GET /movie_index/_doc/1
```

返回整个文档的内容，包括元数据。

### 获取数据

```
GET /movie_index/_source/1
```

### 获取指定字段

```
GET /movie_index/_source/1?_source=title,source_url
```

### 更新文档

```
POST /<index>/_update/<_id>

POST /movie_index/_update/1
{
  "doc": {
    "title": "good movie"
  }
}
```

### 批量获取

```
GET /_mget
GET /<index>/_mget
```

### 删除文档

```
DELETE /movie_index/_doc/1
```

### 检索

```
GET /<target>/_search

GET /_search

POST /<target>/_search

POST /_search
```

- 查询 id=1 的文档。

```
GET /movie_index/_search
{
  "query": {
    "bool": {
      "filter":{
        "term":{"id": 1}
      }
    }
  }
}
```

- 查询 tags 包含 action 的文档。

```
GET /movie_index/_search
{
  "query": {
    "bool": {
      "filter":{
        "match_phrase":{
          "tags": "action"

        }
      }
    }
  }
}
```

- 查询 source 中包含 url 的文档。

```
GET /movie_index/_search
{
  "query": {
    "match_phrase": {
      "source_url": "url"
    }
  }
}
```

### 获取数量

```
GET /<target>/_count

GET /movie_index/_count
{
  "query": {
    "match_phrase": {
      "title": "movie"
    }
  }
}
```

### 聚合

```
# 平均分数。
POST /movie_index/_search?size=0
{
  "aggs": {
    "avg_score": { "avg": { "field": "score"} }
  }
}
```

## 索引操作

### 创建索引

现在我们需要对这个建立索引的过程做更多的控制：我们想要确保这个索引有数量适中的主分片，并且在我们索引任何数据 之前 ，分析器和映射已经被建立好。
为了达到这个目的，我们需要手动创建索引，在请求体里面传入设置或类型映射，如下所示：

```
PUT /my_index
{
    "settings": { ... any settings ... },
    "mappings": {
        "type_one": { ... any mappings ... },
        "type_two": { ... any mappings ... },
        ...
    }
}
```

如果你想禁止自动创建索引，你 可以通过在 config/elasticsearch.yml 的每个节点下添加下面的配置：

```yaml
action.auto_create_index: false
```

#### settings

下属配置项:

- number_of_shards: 每个索引的主分片数，默认值是 5 。这个配置在索引创建后不能修改。
- number_of_replicas: 每个主分片的副本数，默认值是 1 。对于活动的索引库，这个配置可以随时修改。
  ```
  PUT /my_temp_index/_settings
  {
    "number_of_replicas": 1
  }
  ```
- analysis: 来配置已存在的分析器或针对你的索引创建新的自定义分析器。中文一般使用 tk 分词器。
  ```
  PUT /spanish_docs
  {
      "settings": {
          "analysis": {
              "analyzer": {
                  "es_std": {
                      "type":      "standard",
                      "stopwords": "_spanish_"
                  }
              }
          }
      }
  }
  ```
  standard 分析器是用于全文字段的默认分析器，对于大部分西方语系来说是一个不错的选择。它包括了以下几点：
  - standard 分词器，通过单词边界分割输入的文本。
  - standard 语汇单元过滤器，目的是整理分词器触发的语汇单元（但是目前什么都没做）。
  - lowercase 语汇单元过滤器，转换所有的语汇单元为小写。
  - stop 语汇单元过滤器，删除停用词—​ 对搜索相关性影响不大的常用词，如 a ， the ， and ， is 。

### 删除索引

```
<!-- 删除指定索引 -->
DELETE /my_index

<!-- 删除多个 -->
DELETE /index_one,index_two
DELETE /index_*
```

# Golang 客户端

使用`go get`命令下载客户端库文件。这是官方提供的库。

```sh
go get github.com/elastic/go-elasticsearch/v8@latest
```

## 连接

```go
// ES 配置
cfg := elasticsearch.Config{
	Addresses: []string{
		"http://localhost:9200",
	},
}

// 创建客户端连接
client, err := elasticsearch.NewTypedClient(cfg)
if err != nil {
	fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
	return
}
```

## Document

### 创建

```go
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
```

### 检索

```go
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
```

### 更新

```go
resp, err := client.UpdateByQuery("movie").
		Query(&types.Query{MatchPhrase: map[string]types.MatchPhraseQuery{"source": {Query: "souce"}}}).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
```

### 查询数量

```go
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
```

### 删除

```go
resp, err := client.DeleteByQuery("movie").
		Query(&types.Query{Match: map[string]types.MatchQuery{"tags": {Query: "tag2"}}}).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
```

## 索引

### 创建

```go
resp, err := client.Indices.
  Create("my-review-1").
  Do(context.Background())
if err != nil {
  fmt.Printf("create index failed, err:%v\n", err)
  return
}
fmt.Printf("index:%#v\n", resp.Index)
```

### 查找

```go
indices, err := client.Cat.Indices().Do(context.Background())
if err != nil {
  return err
}
for _, index := range indices {
  fmt.Println(*index.Index)
}
return nil
```

### 删除

```go
resp, err := client.Indices.Delete("movie").Do(context.Background())
if err != nil {
  return err
}
fmt.Println(resp.Acknowledged)
return nil
```
