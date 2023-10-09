package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

var (
	Q         = ""
	T         = ""
	phpsessid = ""
	phptime   = 0
)

type Conf struct {
	T         string `json:"T"`
	Q         string `json:"Q"`
	Phpsessid string `json:"phpsessid"`
	Phptime   int    `json:"phptime"`
}

func main() {
	http.HandleFunc("/get_phpsessid", get_phpsessid)
	ticker := time.Tick(20 * time.Minute)
	go func() {
		for range ticker {
			keep()
		}
	}()
	http.ListenAndServe(":14911", nil)
}
func get_phpsessid(w http.ResponseWriter, r *http.Request) {
	phpsessid = run()
	w.Write([]byte(phpsessid))
}
func run() string {
	fmt.Println("当前时间戳：", time_new())
	read_conf() //读取配置文件
	if Q == "" && T == "" {
		fmt.Println("配置参数为空，请补充配置文件")
		return "配置参数为空，请补充配置文件"
	} else if phpsessid == "" && phptime == 0 {
		fmt.Println("php参数数据为空")
		login360() //使用360cookie登录
		return phpsessid
	} else if phptime < int(time_new()) {
		fmt.Println("phpsessid过期")
		login360() //使用360cookie登录
		return phpsessid

	} else if phpsessid != "" && phptime > int(time_new()) {
		fmt.Println("时间戳未过期，直接登录")
		write_conf()
		return phpsessid
	} else {
		fmt.Println("读取配置文件参数时错误")
		return phpsessid
	}
}

// http请求保持session
func keep() {
	fmt.Println("保持登录")
	post := http.Client{}
	url := "https://wangzhan.qianxin.com/totalview/index"
	http, _ := http.NewRequest("GET", url, nil)
	http.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	http.Header.Set("Cookie", "PHPSESSID="+phpsessid)
	res, err := post.Do(http)
	if err != nil {
		fmt.Println("http错误:", err)
	}
	defer res.Body.Close()
	body, _ := (io.ReadAll(res.Body))
	html := string(body)
	if strings.Contains(html, "用户唯一标识：") {
		fmt.Println("session保持成功")
		write_conf()
	} else {
		fmt.Println("session保持失败，正重新获取")
		run()
	}

}

// 获取当前时间戳
func time_new() int64 {
	time_new := time.Now().Unix()
	return time_new
}

// 获取过期时间戳
func time_25() int64 {
	return time.Now().Add(25 * time.Minute).Unix()
}

// 判断go_conf.json配置文件是否存在并读取配置文件
func read_conf() {
	filename := "go_conf.json"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println("配置文件不存在")
		data := map[string]interface{}{
			"//提示":      "下方T、Q的值请填写360登录的cookie里对应的T、Q值,剩下的phpsessid和phptime请忽略,这两个由程序自动生成读写",
			"T":         "",
			"Q":         "",
			"phpsessid": "",
			"phptime":   0,
		}

		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Println("JSON编码失败:", err)
			return
		}

		err = os.WriteFile("go_conf.json", jsonData, 0644)
		if err != nil {
			fmt.Println("写入配置文件失败：", err)
			return
		}

		fmt.Println("配置文件已创建并成功写入基础数据,请在该程序目录下查看名为【go_conf.json的文件】并补充完数据后再重新运行")
		os.Exit(0)
	} else if err == nil {
		file, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("读取文件时出现错误：", err)
			return
		}
		var config Conf
		err = json.Unmarshal(file, &config)
		if err != nil {
			fmt.Println("解析JSON时出现错误: ", err)
			return
		}
		//赋值变量
		Q = config.Q
		T = config.T
		phpsessid = config.Phpsessid
		phptime = config.Phptime
	} else {
		fmt.Println("配置文件存在")
	}
}

// 写入配置文件
func write_conf() {
	if_phpsessid()
	data := map[string]interface{}{
		"//提示":      "下方T、Q的值请填写360登录的cookie里对应的T、Q值,剩下的phpsessid和phptime请忽略,这两个由程序自动生成读写",
		"T":         T,
		"Q":         Q,
		"phpsessid": phpsessid,
		"phptime":   int(time_25()),
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("JSON编码失败:", err)
		return
	}

	err = os.WriteFile("go_conf.json", jsonData, 0644)
	if err != nil {
		fmt.Println("写入配置文件失败：", err)
		return
	}

	fmt.Println("配置更新写入完成")
}

// 使用360cookie登录
func login360() {
	fmt.Println("使用360登录")
	jar, _ := cookiejar.New(nil)
	session := &http.Client{Jar: jar}
	// 创建http请求
	url := "https://openapi.360.cn/oauth2/authorize"
	data := "client_id=02f55d8f4dd80ac05e0a16617df49e26&response_type=code&redirect_uri=https%3A%2F%2Fuser.skyeye.qianxin.com%2F360oauth_redirect&state=http%3A%2F%2Fwangzhan.qianxin.com%2Flogin%2Flogin&scope=&display=default&mid=&version=&DChannel="
	http, _ := http.NewRequest("POST", url, strings.NewReader(data))
	http.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	http.Header.Set("Cookie", "Q="+Q+"; T="+T)
	http.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.188")
	// 发送http请求
	res, err := session.Do(http)
	if err != nil {
		fmt.Println("http错误:", err)
		return
	}
	// 读取响应内容
	defer res.Body.Close()
	body, _ := (io.ReadAll(res.Body))
	html := string(body)
	if strings.Contains(html, "用户唯一标识：") {
		fmt.Println("网站卫士登录成功")
		cookies := session.Jar.Cookies(res.Request.URL)
		for _, cookie := range cookies {
			if cookie.Name == "PHPSESSID" {
				fmt.Println("获取到的PHPSESSID:", cookie.Value)
				phpsessid = cookie.Value
				break
			}
		}
		write_conf()
	} else if strings.Contains(html, "您正在访问的应用暂时无法正常提供服务") {
		fmt.Println("360Cooke可能已经失效，请更新cookie再次尝试")
	} else if strings.Contains(html, "什么都没有发现啊") {
		fmt.Println("360或奇安信响应过忙，等下重试看看")
	} else {
		fmt.Println("网站卫士登录失败")
		phpsessid = ""
	}

}

// 验证phpsessid格式
func if_phpsessid() {
	if len(phpsessid) == 32 && phpsessid != "" {
		fmt.Println("phpsessid验证成功")
	} else {
		fmt.Println("phpsessid验证失败,已停止运行")
		return
	}

}
