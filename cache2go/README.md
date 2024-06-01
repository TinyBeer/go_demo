cache2go 支持过期的并发安全 go 缓存库。
cache2go 使用 map 作为存储 item 的容器，使用定时器进行过期检测，同时支持时间触发回调。

# 安装

```sh
go get github.com/muesli/cache2go
```

# 官方使用示例

```go
// Keys & values in cache2go can be of arbitrary types, e.g. a struct.
// 在cache2go中 key & values 可以是任何类型 比如 struc
type myStruct struct {
	text     string
	moreData []byte
}

func main() {
	// Accessing a new cache table for the first time will create it.
    // 第一次访问一个新的cache table会进行创建
	cache := cache2go.Cache("myCache")

	// We will put a new item in the cache. It will expire after
	// not being accessed via Value(key) for more than 5 seconds.
    // 我们向cache中放入一个新的item，他将会在未被访问后5秒过期
	val := myStruct{"This is a test!", []byte{}}
	cache.Add("someKey", 5*time.Second, &val)

	// Let's retrieve the item from the cache.
    // 从cache中获取这个item
	res, err := cache.Value("someKey")
	if err == nil {
		fmt.Println("Found value in cache:", res.Data().(*myStruct).text)
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	// Wait for the item to expire in cache.
    // 等待这个item过期
	time.Sleep(6 * time.Second)
	res, err = cache.Value("someKey")
	if err != nil {
		fmt.Println("Item is not cached (anymore).")
	}

	// Add another item that never expires.
    // 添加另一个永不过期的item
	cache.Add("someKey", 0, &val)

	// cache2go supports a few handy callbacks and loading mechanisms.
    // cache2go支持一些方便的回调和加载机制
	cache.SetAboutToDeleteItemCallback(func(e *cache2go.CacheItem) {
		fmt.Println("Deleting:", e.Key(), e.Data().(*myStruct).text, e.CreatedOn())
	})

	// Remove the item from the cache.
    // 移除这个item
	cache.Delete("someKey")

	// And wipe the entire cache table.
    // 清理整个cache table
	cache.Flush()
}

```

# 使用示例补充

## 注册 item 过期删除回调

`SetAboutToExpireCallback`和`AddAboutToExpireCallback`都可以进行 item 鼓起删除回调注册，区别是`SetAboutToExpireCallback`会删除掉之前注册的过期删除回到，而`AddAboutToExpireCallback`则只会向后追加。
另外，除了超时以外，`Delete`也可以触发过期删除回调操作。

```go
func main() {
	cache := cache2go.Cache("my_cache")
	item := cache.Add("hello", time.Second, "world")
	item.SetAboutToExpireCallback(func(key interface{}) {
		fmt.Println(key, "deleted")
	})
	time.Sleep(time.Second * 2)

	item = cache.Add("hello", time.Second, "world")
	item.AddAboutToExpireCallback(func(key interface{}) {
		fmt.Println(key, "deleted")
	})
	item.AddAboutToExpireCallback(func(key interface{}) {
		fmt.Println(key, "deleted2")
	})
	cache.Delete("hello")
}
```

```sh
hello deleted
hello deleted
hello deleted2
```

# 核心数据结构梳理

## CacheTable

```go
type CacheTable struct {
	sync.RWMutex

	// The table's name.
	name string
	// All cached items.
	items map[interface{}]*CacheItem

	// Timer responsible for triggering cleanup.
	cleanupTimer *time.Timer
	// Current timer duration.
	cleanupInterval time.Duration

	// The logger used for this table.
	logger *log.Logger

	// Callback method triggered when trying to load a non-existing key.
	loadData func(key interface{}, args ...interface{}) *CacheItem
	// Callback method triggered when adding a new item to the cache.
	addedItem []func(item *CacheItem)
	// Callback method triggered before deleting an item from the cache.
	aboutToDeleteItem []func(item *CacheItem)
}
```

1. sync.RWMutex
   一把读写锁，用于控制并发读写
2. name
   缓存表名，用于区分不同表
3. items
   实际存储缓存对象的地方，是一个 map。也就是 key 支持 comparable 对象，而不是所有类型。比如 slice、func、map 等就不支持。
4. cleanupTimer
   用于触发清理任务的定时器
5. cleanupInterval
   上次清理到现在的时间
6. logger
   日志组件
7. loadData
   读取不存在的 key 时触发的回调方法
8. addItem
   当添加新 item 时触发的回调方法
9. aboutToDeleteItem
   删除 item 前触发的回调方法

## CacheItem

```go
// CacheItem is an individual cache item
// Parameter data contains the user-set value in the cache.
type CacheItem struct {
	sync.RWMutex

	// The item's key.
	key interface{}
	// The item's data.
	data interface{}
	// How long will the item live in the cache when not being accessed/kept alive.
	lifeSpan time.Duration

	// Creation timestamp.
	createdOn time.Time
	// Last access timestamp.
	accessedOn time.Time
	// How often the item was accessed.
	accessCount int64

	// Callback method triggered right before removing the item from the cache
	aboutToExpire []func(key interface{})
}
```

- lifeSpan
  未被访问后存活的时长
- createdOn
  创建时间
- accessedOn
  最后一次访问时间，配合 lifeSpan 进行生命周期判定
- accessCount
  访问计数
- aboutToExpire
  在被移除前触发的回调方法

# 执行流程梳理

## Cache

我们可以通过`Cache`获取一个 `cache table`，如这个`cache table`不存在，则会创建一个新的。执行流程如下：

1. 在读锁保护下获取名为 name 的`cache table`信息
2. 判断是否存在，如果存在则直接返回，否则执行 3
3. 在写锁保护下，再次检查是否存在对应`cache table`，如果不存在则创建并注册到 cache 中，否则返回。再次判断是为了避免在未持有锁期间其他协程进行了创建操作引起的并发问题。

> mutex & cache 私有全集变量

```go
var (
  cache = make(map[string]*CacheTable)
  mutex sync.RWMutex
)
```

```go
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			t = &CacheTable{
				name:  table,
				items: make(map[interface{}]*CacheItem),
			}
			cache[table] = t
		}
		mutex.Unlock()
	}

	return t
}
```

## Add

通过`Add`方法我们可以向`CacheTable`存入 item。执行流程如下：

1. 根据参数创建 item
2. 加锁进行表更新操作

   - 更新`table.items`，如果有 key 相同的 item 会被覆盖掉。
   - 缓存` table.cleanupInterval``table.addedItem `后释放锁
   - 执行`table.addedItem`回调函数
   - 如果存入的 item 不是永久有效并且没有触发过过期检测或者 item 过期时长小于下一次触发间隔，则执行更新间隔操作

3. 返回 item，这使得我们可以对 item 添加回调方法。

```go
func (table *CacheTable) Add(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
	item := NewCacheItem(key, lifeSpan, data)

	// Add item to cache.
	table.Lock()
	table.addInternal(item)

	return item
}
```

```go
func (table *CacheTable) addInternal(item *CacheItem) {
	// Careful: do not run this method unless the table-mutex is locked!
	// It will unlock it for the caller before running the callbacks and checks
	table.log("Adding item with key", item.key, "and lifespan of", item.lifeSpan, "to table", table.name)
	table.items[item.key] = item

	// Cache values so we don't keep blocking the mutex.
	expDur := table.cleanupInterval
	addedItem := table.addedItem
	table.Unlock()

	// Trigger callback after adding an item to cache.
	if addedItem != nil {
		for _, callback := range addedItem {
			callback(item)
		}
	}

	// If we haven't set up any expiration check timer or found a more imminent item.
	if item.lifeSpan > 0 && (expDur == 0 || item.lifeSpan < expDur) {
		table.expirationCheck()
	}
}
```

## Value

我们使用`Value`方法获取缓存的 item，执行流程如下：

1. 在读锁保护性获取对应 item 信息，同时混村`table.loadData`回调
2. 判断 item 是否存在，存在则执行 KeepAlive 方法，更新最后一次访问时间和访问计数器，同时返回 item。否则执行 3
3. 判断是否注册了 loadData，如果有则执行 loadData 尝试加载 item，成功则返回 item，否则返回错误`ErrKeyNotFoundOrLoadable`

```go
func (table *CacheTable) Value(key interface{}, args ...interface{}) (*CacheItem, error) {
	table.RLock()
	r, ok := table.items[key]
	loadData := table.loadData
	table.RUnlock()

	if ok {
		// Update access counter and timestamp.
		r.KeepAlive()
		return r, nil
	}

	// Item doesn't exist in cache. Try and fetch it with a data-loader.
	if loadData != nil {
		item := loadData(key, args...)
		if item != nil {
			table.Add(key, item.lifeSpan, item.data)
			return item, nil
		}

		return nil, ErrKeyNotFoundOrLoadable
	}

	return nil, ErrKeyNotFound
}
```

## Delete

1. 获取写锁，并利用 defer 在流程结束后进行一次释放锁操作。
2. 检查 key 是否存在，不存在直接报错。否则执行 3。
3. 缓存`table.aboutToDeleteItem`回调函数，并释放锁资源。
4. 如果`table.aboutToDeleteItem`不为 nil 则执行
5. 获取读锁，如果注册了过期回调则执行。
6. 加写锁，删除 item，并返回。

```go
// Delete an item from the cache.
func (table *CacheTable) Delete(key interface{}) (*CacheItem, error) {
	table.Lock()
	defer table.Unlock()

	return table.deleteInternal(key)
}
```

```go
func (table *CacheTable) deleteInternal(key interface{}) (*CacheItem, error) {
	r, ok := table.items[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	// Cache value so we don't keep blocking the mutex.
	aboutToDeleteItem := table.aboutToDeleteItem
	table.Unlock()

	// Trigger callbacks before deleting an item from cache.
	if aboutToDeleteItem != nil {
		for _, callback := range aboutToDeleteItem {
			callback(r)
		}
	}

	r.RLock()
	defer r.RUnlock()
	if r.aboutToExpire != nil {
		for _, callback := range r.aboutToExpire {
			callback(key)
		}
	}

	table.Lock()
	table.log("Deleting item with key", key, "created on", r.createdOn, "and hit", r.accessCount, "times from table", table.name)
	delete(table.items, key)

	return r, nil
}
```

## Flush

1. 加写锁，使用 defer 实现流程结束后释放锁
2. 创建新的`table.items`覆盖原来的，让 GC 自行清理数据。
3. 重置定时器信息。

```go
// Flush deletes all items from this cache table.
func (table *CacheTable) Flush() {
	table.Lock()
	defer table.Unlock()

	table.log("Flushing table", table.name)

	table.items = make(map[interface{}]*CacheItem)
	table.cleanupInterval = 0
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}
}
```
