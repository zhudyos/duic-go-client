package duic

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
)

var (
	BaseUri  string // duic 配置中心请求地址
	Name     string // 应用名称
	Profiles string // 环境名称，多个使用英文逗号分隔
	Tokens   string // 认证信息，多个使用英文逗号分隔
)

var (
	hc      *http.Client
	configs map[string]interface{}
	state   string
)

// 初始化配置数据。
func Init() {
	tr := &http.Transport{
		MaxIdleConns:    2,
		IdleConnTimeout: 60 * time.Second,
	}

	hc = &http.Client{
		Transport: tr,
	}

	loadState()
	loadConf()

	go func() {
		for {
			reload()
		}
	}()
}

// 返回一个 Bool 配置，如果配置不存在或者是一个错误的 Bool 值则将返回错误。
func Bool(k string) (bool, error) {
	v, err := getV(k)
	if err != nil {
		return false, err
	}

	switch v.(type) {
	case bool:
		return v.(bool), nil
	case float64:
		return v.(float64) != 0, nil
	case string:
		b, err := strconv.ParseBool(v.(string))
		if err != nil {
			return false, err
		}
		return b, nil
	default:
		return false, nil
	}
}

// 返回一个 bool 配置，如果配置不存在或者是一个错误的 bool 值则将返回默认值。
func Bool2(k string, defVar bool) bool {
	v, err := Bool(k)
	if err != nil {
		return defVar
	}
	return v
}

// 返回一个 int64 配置，如果配置不存在或者是一个错误的 int64 值则将返回错误。
func Int64(k string) (int64, error) {
	v, err := Float64(k)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

// 返回一个 int64 配置，如果配置不存在或者是一个错误的 int64 值则将返回默认值。
func Int642(k string, defVar int64) int64 {
	v, err := Int64(k)
	if err != nil {
		return defVar
	}
	return v
}

// 返回一个 float64 配置，如果配置不存在或者是一个错误的 float64 值则将返回错误。
func Float64(k string) (float64, error) {
	v, err := getV(k)
	if err != nil {
		return 0.0, err
	}

	switch v.(type) {
	case float64:
		return v.(float64), nil
	case string:
		f, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return 0.0, err
		}
		return f, nil
	default:
		return 0.0, fmt.Errorf("值：%v 不能转换为 float64", v)
	}
}

// 返回一个 float64 配置，如果配置不存在或者是一个错误的 float64 值则将返回默认值。
func Float642(k string, defVar float64) float64 {
	v, err := Float64(k)
	if err != nil {
		return defVar
	}
	return v
}

// 返回一个 string 配置，如果配置不存在或者是一个错误的 string 值则返回错误。
func String(k string) (string, error) {
	v, err := getV(k)
	if err != nil {
		return "", err
	}

	switch v.(type) {
	case string:
		return v.(string), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// 返回一个 string 配置，如果配置不存在或者是一个错误的 string 值则将返回默认值。
func String2(k string, defVar string) string {
	v, err := String(k)
	if err != nil {
		return defVar
	}
	return v
}

// 返回一个数组配置，如果配置不存在或者是一个错误类型则将返回错误。
func Array(k string) ([]interface{}, error) {
	v, err := getV(k)
	if err != nil {
		return nil, err
	}

	switch v.(type) {
	case []interface{}:
		return v.([]interface{}), nil
	}
	return nil, fmt.Errorf("%v 错误的数据类型", v)
}

// 返回一个对象配置，如果配置不存在或者是一个错误的类型则将返回错误。
func Object(k string) (map[string]interface{}, error) {
	v, err := getV(k)
	if err != nil {
		return nil, err
	}

	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{}), nil
	}
	return nil, fmt.Errorf("%v 错误的数据类型", v)
}

func getV(k string) (interface{}, error) {
	arr := strings.Split(k, ".")

	var v interface{} = configs
	for _, e := range arr {
		switch v.(type) {
		case map[string]interface{}:
			tv := v.(map[string]interface{})
			v = tv[e]
		default:
			return nil, fmt.Errorf("错误的数据 %v", k)
		}
	}

	return v, nil
}

func reload() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	s, err := watchState()
	if err != nil {
		log.Println(err)
	}

	if s != state {
		state = s
		loadConf()
		log.Printf("DuiC config reload successfully. newState: %s\n", s)
	}
}

func watchState() (string, error) {
	url := fmt.Sprintf("%s/apps/watches/%s/%s?state=%s", BaseUri, Name, Profiles, state)
	req := newReq(url)
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return "", fmt.Errorf("GET %s 失败，httpStatus: %v", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	var m map[string]string
	json.Unmarshal(body, &m)

	return m["state"], nil
}

func loadState() {
	url := fmt.Sprintf("%s/apps/states/%s/%s", BaseUri, Name, Profiles)
	req := newReq(url)
	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		panic(fmt.Errorf("GET %s 失败，httpStatus: %v", url, resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var m map[string]string
	json.Unmarshal(body, &m)

	state = m["state"]

	log.Printf("Get config state from %s [%s]\n", url, state)
}

// 加载配置中心配置数据
func loadConf() {
	url := fmt.Sprintf("%s/apps/%s/%s", BaseUri, Name, Profiles)
	req := newReq(url)
	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		panic(fmt.Errorf("GET %s 失败，httpStatus: %v", url, resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var m map[string]interface{}
	json.Unmarshal(body, &m)
	configs = m

	log.Printf("Fetch config from %s\n", url)
}

func newReq(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	if Tokens != "" {
		req.Header.Add("x-config-token", Tokens)
	}
	return req
}
