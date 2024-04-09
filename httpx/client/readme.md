# http client
重写的http client. 扩展性更好，使用体验更好。

请求方法:
- [x] PostJson
- [ ] PostXml
- [ ] PostForm
- [x] PutJson
- [ ] PutXml
- [x] Get
- [x] Delete

解析响应方法：
- [x] Response: 返回*http.Response自己处理
- [x] ScanJson: 按json格式解析响应body
- [x] ScanXml: 按xml格式解析响应body

## 示例
```go
package example

import "github.com/zeddy-go/zeddy/http/client"

var req struct{
	Field1 int `json:"field1"`
}
var resp struct{
	Field1 int `json:"field1"`
}
//直接用
client := http.NewClient()

resp, err := client.PostJson("http://www.baidu.com/some/api", req)
if err != nil {
	panic(err)
}
defer resp.Close()
err = resp.ScanJson(&resp)
if err != nil {
	panic(err)
}

//预先设置base url
client := http.NewClient(http.WithBaseUrl(baseUrl))
resp, err = client.PostJson("some/api", req)
if err != nil {
	panic(err)
}
defer resp.Close()
err = resp.ScanJson(&resp)
if err != nil {
	panic(err)
}
```