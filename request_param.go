package rhttp

// PostType Post方法类型定义
type PostType int

const (
	PostJson PostType = 0
	PostForm PostType = 1
)

// CRequestParam 请求参数定义
type CRequestParam struct {
	Address          string
	PostDatas        string
	PostDataType     PostType
	ProxyAddr        string
	ProxyPort        string
	ProxyUserName    string
	ProxyPassword    string
	ConnectionHeader map[string]string
	Timeout          int
}

// OrigConnectHeader 初始化Connection Header
func (ths *CRequestParam) OrigConnectHeader() {
	ths.ConnectionHeader = make(map[string]string)
	ths.ConnectionHeader["accept"] = "application/json, text/plain, */*"
	ths.ConnectionHeader["accept-encoding"] = "gzip, deflate"
	ths.ConnectionHeader["accept-language"] = "zh-CN,zh;q=0.8"
	ths.ConnectionHeader["cache-control"] = "no-cache"
	ths.ConnectionHeader["pragma"] = "no-cache"
	ths.ConnectionHeader["user-agent"] = "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
}
