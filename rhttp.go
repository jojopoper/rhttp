package rhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// ReturnType 返回类型定义
type ReturnType int

// DecodeFunction 解码消息体的方法定义
type DecodeFunction func(b []byte) (interface{}, error)

const (
	ReturnMap        ReturnType = 1
	ReturnSlice      ReturnType = 2
	ReturnString     ReturnType = 3
	ReturnSliceByte  ReturnType = 4
	ReturnCustomType ReturnType = 5
)

// TimeStatics 时间统计定义
type TimeStatics struct {
	StartSend    int64
	CompleteSend int64
}

// CHttp custom http
type CHttp struct {
	RClient
	TimeStatics
	client     *http.Client
	clientConn *httputil.ClientConn
	request    *http.Request
	customFun  DecodeFunction
}

// SetClient 设置http client
func (ths *CHttp) SetClient(c *http.Client) {
	ths.client = c
}

// SetClientConn 设置http client connection
func (ths *CHttp) SetClientConn(c *httputil.ClientConn) {
	ths.clientConn = c
}

// SetDecodeFunc 设置自定义解码消息体函数
func (ths *CHttp) SetDecodeFunc(f DecodeFunction) {
	ths.customFun = f
}

// Get 使用原生http get
func (ths *CHttp) Get(addr string, retType ReturnType) (interface{}, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, fmt.Errorf("[ CHttp:Get ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// ClientGet 通过Http client 获取
func (ths *CHttp) ClientGet(addr string, retType ReturnType) (interface{}, error) {
	resp, err := ths.client.Get(addr)
	if err != nil {
		ths.client = ths.GetClient(ths.timeout)
		return nil, fmt.Errorf("[ CHttp:ClientGet ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// PostForm 使用原生http post form
func (ths *CHttp) PostForm(address string, retType ReturnType, data string) (interface{}, error) {
	resp, err := http.Post(address, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("[ CHttp:PostForm ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// PostJSON 使用原生http post json
func (ths *CHttp) PostJSON(address string, retType ReturnType, data string) (interface{}, error) {
	resp, err := http.Post(address, "application/json", strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("[ CHttp:PostJson ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// ClientPostForm 使用http client post form
func (ths *CHttp) ClientPostForm(address string, retType ReturnType, data string) (interface{}, error) {
	resp, err := ths.client.Post(address, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		ths.client = ths.GetClient(ths.timeout)
		return nil, fmt.Errorf("[ CHttp:ClientPostForm ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// ClientPostJSON 使用 client post json
func (ths *CHttp) ClientPostJSON(address string, retType ReturnType, data string) (interface{}, error) {
	resp, err := ths.client.Post(address, "application/json", strings.NewReader(data))
	if err != nil {
		ths.client = ths.GetClient(ths.timeout)
		return nil, fmt.Errorf("[ CHttp:ClientPostJson ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// ClientPostFormWithHeader 使用 client post form with header
func (ths *CHttp) ClientPostFormWithHeader(address string, retType ReturnType, data string, header map[string]string) (ret interface{}, err error) {
	ths.request, err = http.NewRequest("POST", address, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		ths.request.Header.Set(k, v)
	}
	ths.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := ths.client.Do(ths.request)
	if err != nil {
		ths.client = ths.GetClient(ths.timeout)
		return nil, fmt.Errorf("[ CHttp:ClientPostFormWithHeader ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// ClientPostJsonWithHeader 使用 client post json with header
func (ths *CHttp) ClientPostJsonWithHeader(address string, retType ReturnType, data string, header map[string]string) (ret interface{}, err error) {
	ths.request, err = http.NewRequest("POST", address, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		ths.request.Header.Set(k, v)
	}
	ths.request.Header.Set("Content-Type", "application/json")
	resp, err := ths.client.Do(ths.request)
	if err != nil {
		ths.client = ths.GetClient(ths.timeout)
		return nil, fmt.Errorf("[ CHttp:ClientPostJsonWithHeader ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

// ClientConnGet 通过Http client connection获取
func (ths *CHttp) ClientConnGet(addr string, header map[string]string) (err error) {
	ths.request, err = http.NewRequest("GET", addr, nil)
	if err != nil {
		return err
	}
	for k, v := range header {
		ths.request.Header.Set(k, v)
	}

	ths.StartSend = time.Now().UnixNano()
	return ths.clientConn.Write(ths.request)
}

// ClientConnPostForm 通过Http client connection post form
func (ths *CHttp) ClientConnPostForm(addr, data string, header map[string]string) (err error) {
	ths.request, err = http.NewRequest("POST", addr, strings.NewReader(data))
	if err != nil {
		return
	}
	for k, v := range header {
		ths.request.Header.Set(k, v)
	}
	ths.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ths.StartSend = time.Now().UnixNano()
	return ths.clientConn.Write(ths.request)
}

// ClientConnPostJSON 通过Http client connection post json
func (ths *CHttp) ClientConnPostJSON(addr, data string, header map[string]string) (err error) {
	ths.request, err = http.NewRequest("POST", addr, strings.NewReader(data))
	if err != nil {
		return
	}
	for k, v := range header {
		ths.request.Header.Set(k, v)
	}
	ths.request.Header.Set("Content-Type", "application/json")

	ths.StartSend = time.Now().UnixNano()
	return ths.clientConn.Write(ths.request)
}

// ClientConnResponse 通过Http client connection response
func (ths *CHttp) ClientConnResponse(retType ReturnType) (interface{}, error) {
	resp, err := ths.clientConn.Read(ths.request)
	ths.CompleteSend = time.Now().UnixNano()
	if err != nil {
		return nil, fmt.Errorf("[ CHttp:ClientConnResponse ] Has error:\r\n%+v", err)
	}
	return ths.decode(resp, retType)
}

func (ths *CHttp) decode(resp *http.Response, retType ReturnType) (interface{}, error) {
	if resp == nil {
		return nil, fmt.Errorf("[ CHttp:decode ] http.Response is nil")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ CHttp:decode ] Has error[0]:\n%+v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("[ CHttp:decode ] Has error:\n Statues code : %d\n Status : %s\n Body : %s", resp.StatusCode, resp.Status, string(body))
	}

	switch retType {
	case ReturnMap:
		result := make(map[string]interface{})
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, fmt.Errorf("[ CHttp:decode ] Unmarshal map[string]interface{} has error[0]:\r\nerror is : %+v\r\nreturn body is : %s", err, string(body))
		}
		return result, nil
	case ReturnSlice:
		result := make([]interface{}, 0)
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, fmt.Errorf("[ CHttp:decode ] Unmarshal []interface{} has error[0]:\r\nerror is : %+v\r\nreturn body is : %s", err, string(body))
		}
		return result, nil
	case ReturnString:
		return string(body), nil
	case ReturnSliceByte:
		return body, nil
	case ReturnCustomType:
		if ths.customFun != nil {
			return ths.customFun(body)
		}
		return nil, fmt.Errorf("[ CHttp:decode ] Custom return type has not valid execute function.\r\nreturn body is : %s", string(body))
	}
	return nil, fmt.Errorf("[ CHttp:decode ] Unknown return type [ %d ]\r\nreturn body is : %s", retType, string(body))
}
