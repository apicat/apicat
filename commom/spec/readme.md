# the ApiCat Specification

主要包含apicat的协议定义 和 对应的插件

### spec

```go
specObj,err:=spec.ParseJSON([]byte(...))
jsonraw,err:=specObj.ToJSON(spec.{})
```


### plugin
* openapi 支持与openapi2(swagger)以及3.x之间的相互转换

