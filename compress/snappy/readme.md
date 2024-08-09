[golang/snappy](https://github.com/golang/snappy) 是go语言实现的snappy压缩库。

# snappy简介

Snappy 是一个 C++ 的用来压缩和解压缩的开发包。其目标不是最大限度压缩或者兼容其他压缩格式，而是旨在提供高速压缩速度和合理的压缩率。Snappy 比 zlib 更快，但文件相对要大 20% 到 100%。在 64位模式的 Core i7 处理器上，可达每秒 250~500兆的压缩速度。
Snappy采用新BSD协议开源。

# 使用

```sh
go get github.com/golang/snappy
```

```go
compressed := snappy.Encode(nil, origin)

decompressed, err := snappy.Decode(nil, compressed)
```

# 与其他算法的对比
