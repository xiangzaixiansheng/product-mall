# Go 速查手册

## 常用命令

```bash
# 初始化项目
go mod init module-name

# 下载依赖
go mod download

# 整理依赖 (删除未使用，添加缺失)
go mod tidy

# 编译
go build ./...
go build -o appname .

# 运行
go run .

# 测试
go test ./...
go test -v ./...
go test -cover ./...

# 代码检查
go vet ./...
go fmt ./...

# 依赖管理
go list -m all          # 列出所有依赖
go get package@version # 添加依赖
go mod why package     # 依赖原因
```

## 基础语法

### 变量声明
```go
// 显式类型
var name string = "test"

// 类型推断
var name = "test"

// 短声明 (函数内)
name := "test"

// 常量
const PI = 3.14

// 多变量
x, y := 1, 2
```

### 数据类型
```go
// 基础类型
bool, string, int, int8/16/32/64, uint, float32/64, complex64/128

// 数组
var arr [5]int
arr := [...]int{1, 2, 3}

// 切片
slice := []int{1, 2, 3}
slice := make([]int, 0, 10)

// Map
m := map[string]int{"a": 1}
m := make(map[string]int)

// 结构体
type User struct {
    Name string
    Age  int
}
```

### 控制流
```go
// if
if x > 0 {
    // ...
} else if x < 0 {
    // ...
} else {
    // ...
}

// for (三种形式)
for i := 0; i < 10; i++ {}
for i < 10 {}
for {} // 无限循环

// switch
switch v {
case 1:
    // ...
case 2:
    // ...
default:
    // ...
}
```

### 函数
```go
// 基本函数
func add(a, b int) int {
    return a + b
}

// 多返回值
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 命名返回值
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // 隐式返回
}

// 变参
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 函数作为参数
func apply(fn func(int) int, value int) int {
    return fn(value)
}
```

## 常用标准库

### context
```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 传递值
ctx := context.WithValue(ctx, "key", "value")
value := ctx.Value("key")
```

### time
```go
now := time.Now()
later := now.Add(24 * time.Hour)
duration := later.Sub(now)

// 格式化
fmt.Println(now.Format("2006-01-02 15:04:05"))

// 解析
t, _ := time.Parse("2006-01-02", "2024-01-01")
```

### strings
```go
strings.Contains(s, "test")
strings.Split(s, ",")
strings.Join([]string{"a", "b"}, "-")
strings.TrimSpace(s)
strings.ToLower(s)
strings.ToUpper(s)
strings.Replace(s, "old", "new", -1)
```

### strconv
```go
strconv.Atoi("123")
strconv.Itoa(123)
strconv.ParseInt("123", 10, 64)
strconv.FormatInt(123, 10)
strconv.ParseFloat("3.14", 64)
```

### json
```go
// 序列化
data, _ := json.Marshal(obj)
jsonStr := string(data)

// 反序列化
var obj Type
json.Unmarshal([]byte(jsonStr), &obj)

// tag
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

### errors
```go
errors.New("message")
err := errors.New("not found")

// 错误包装
fmt.Errorf("failed: %w", err)

// 检查错误
errors.Is(err, io.EOF)
errors.As(err, &myError{})
```

### log/slog
```go
slog.Info("message", "key", value)
slog.Error("error", "err", err)
slog.Debug("debug", "value", 123)
slog.Warn("warning")

// 带 key-value
slog.Info("user action", "user_id", 1, "action", "login")
```

### os
```go
os.ReadFile("file.txt")
os.WriteFile("file.txt", data, 0644)
os.MkdirAll("dir", 0755)
os.Remove("file.txt")
os.Getenv("KEY")
os.Setenv("KEY", "value")
```

### io
```go
io.ReadAll(reader)
io.Copy(writer, reader)
io.ReadFull(buf, reader)

// 常用接口
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

## 并发

### goroutine
```go
go func() {
    // 异步执行
}()

// 带等待
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()
```

### channel
```go
ch := make(chan int)
ch <- 1
v := <-ch
close(ch)

// 带缓冲
ch := make(chan int, 10)

// select
select {
case v := <-ch:
    fmt.Println(v)
case <-time.After(time.Second):
    fmt.Println("timeout")
}
```

### sync
```go
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

var once sync.Once
once.Do(func() {}) // 只执行一次

var cond sync.Cond
cond.L = &mu
cond.Wait()
cond.Signal()
```

## 常用语法糖

### defer
```go
func() {
    file, _ := os.Open("file")
    defer file.Close()
    // file 会在函数结束时关闭
}()
```

### 切片操作
```go
s := []int{1, 2, 3, 4, 5}

// 头部
s = s[1:]

// 尾部
s = s[:len(s)-1]

// 追加
s = append(s, 6)

// 删除索引 i
s = append(s[:i], s[i+1:]...)
```

### range
```go
for i, v := range arr {
    fmt.Println(i, v)
}

for k, v := range m {
    fmt.Println(k, v)
}

for range ch { // 只等待关闭
}
```

### 泛型 (Go 1.18+)
```go
// 定义
func genericFunc[T any](v T) T {
    return v
}

// 结构体泛型
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}
```

### interface
```go
// 空接口
var v any = "string"

// 类型断言
str, ok := v.(string)
switch v.(type) {
case int:
    // ...
case string:
    // ...
}
```

### init 函数
```go
func init() {
    // 在包导入时执行
}
```

## GORM 常用操作

```go
// 创建
db.Create(&user)

// 查询
db.First(&user, id)
db.Where("name = ?", "Tom").First(&user)

// 更新
db.Model(&user).Update("name", "Jerry")
db.Updates(&user)

// 删除
db.Delete(&user)

// 软删除 (带 DeletedAt)
db.Delete(&user) // 只是标记

// 条件
db.Where("age > ?", 18).Find(&users)
db.Not("name = ?", "Tom").Find(&users)
db.Or("age > ?", 30).Find(&users)
```

## Gin 基础

```go
// 路由
r := gin.Default()
r.GET("/path", handler)
r.POST("/path", handler)
r.PUT("/path", handler)
r.DELETE("/path", handler)

// 路由组
v1 := r.Group("/api/v1")
v1.GET("/users", handler)

// 参数获取
c.Query("name")
c.Param("id")
c.PostForm("name")
c.Bind(&obj)

// JSON响应
c.JSON(200, gin.H{"message": "ok"})
c.JSON(400, gin.H{"error": "bad request"})

// 中间件
r.Use(middleware.Logger())
r.Use(middleware.Auth())
```

## 测试

```go
// 基础测试
func TestAdd(t *testing.T) {
    result := add(1, 2)
    if result != 3 {
        t.Errorf("expected 3, got %d", result)
    }
}

// 表驱动测试
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"2+2", 2, 2, 4},
        {"3+5", 3, 5, 8},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("add() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 常见错误处理模式

```go
// 标准错误处理
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// 忽略特定错误
if errors.Is(err, os.ErrNotExist) {
    // 文件不存在
}

// defer 恢复 panic
defer func() {
    if r := recover(); r != nil {
        fmt.Println("recovered:", r)
    }
}()
```