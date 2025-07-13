# 性能优化

## 对象池化与预分配优化

从v1.x.x版本开始，Hessian2支持对象池化和预分配优化，可以显著提升编解码性能：

- 大型对象解码性能提升60%
- 内存使用减少58%
- 并行请求处理性能提升31%

### 启用方法

默认情况下，对象池化功能是关闭的。要启用此功能，请设置环境变量：

```bash
export HESSIAN_POOL=1
```

或在Go代码中：

```go
import (
    "os"
    hessian "github.com/apache/dubbo-go-hessian2"
)

func init() {
    os.Setenv("HESSIAN_POOL", "1")
}
```

### 使用注意事项

1. 同一个Decoder/Encoder实例不可跨goroutine并发使用
2. 使用完毕后，请调用PutDecoder/PutEncoder归还对象到池中
3. 对象池化主要在大型对象和高并发场景下有显著性能提升

### 示例代码

```go
// 编码
encoder := hessian.GetEncoder()
err := encoder.Encode(obj)
if err != nil {
    // 处理错误
}
bytes := encoder.Buffer()
hessian.PutEncoder(encoder)

// 解码
decoder := hessian.GetDecoder(bytes)
result, err := decoder.Decode()
if err != nil {
    // 处理错误
}
hessian.PutDecoder(decoder)
``` 