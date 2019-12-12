# trietree
一个敏感词检测服务，包含替换功能

## 简易使用

```golang
package main

import (
    "net/http"
    "flag"
    "fmt"
    "log"
    "github.com/guoruibiao/trietree"
    )

var tree *trietree.TrieTree

func init() {
    tree = trietree.NewTrie()
}

func index(writer http.ResponseWriter, request *http.Request) {
    indexhtml := `
    请求格式：
        1. http://ip:port/
            首页
        2. http://ip:port/addword?word=xxx
            添加违禁词
        3. http://ip:port/reload?filepath=xxx
            重新载入
        4. http://ip:port/filter?text=xxx&replace=xxx
            检测是否有违禁词，并进行替换
`
    fmt.Fprintln(writer, indexhtml)
}

func addword(writer http.ResponseWriter, request *http.Request) {
    word := request.URL.Query().Get("word")
    if word == "" {
        fmt.Fprintln(writer, "您未添加单词哦")
    }else {
        tree.AddWord(word, nil)
        fmt.Fprintln(writer, "单词【" + word + "】已添加到词库")
    }
}

func exists(writer http.ResponseWriter, request *http.Request) {
    word := request.URL.Query().Get("word")
    if word == "" {
        fmt.Fprintln(writer, "您未添加单词哦")
    }else {
        if tree.WordExists(word) {
            fmt.Fprintln(writer, "单词【" + word + "】在词库中")
        }else{
            fmt.Fprintln(writer, "单词【" + word + "】不在词库中")
        }
    }
}

func reload(writer http.ResponseWriter, request *http.Request) {
    filepath := request.URL.Query().Get("filepath")
    separator := request.URL.Query().Get("separator")
    fmt.Println("filepath="+filepath+", separator="+separator)
    if filepath == "" || separator == ""{
        fmt.Fprintln(writer, "filepath,separator均不能为空")
    }else{
        tree = nil
        tree = trietree.NewTrie()
        tree.BuildTrieTreeFromFile(filepath, separator)
        fmt.Fprintln(writer, "新的TrieTree已重建")
    }
}

func filter(writer http.ResponseWriter, request *http.Request) {
    text := request.URL.Query().Get("text")
    replace := request.URL.Query().Get("replace")

    if replace == "" {
        replace = "***"
    }

    if text == "" {
        fmt.Fprintln(writer, "待检测文本不能为空")
    }else {
        ret, hit := tree.Filter(text, replace)
        if hit == true {
            fmt.Fprintln(writer, "hit, final result: " + ret)
        }else{
            fmt.Fprintln(writer, "not hit, final result: " + ret)
        }
    }
}

func main() {
    // 整几个命令行参数
    port := flag.String("p", "8080", "监听端口")
    flag.Parse()

    fmt.Println("start listening at `http://localhost:" + *port + "/` ...")
    http.HandleFunc("/filter", filter)
    http.HandleFunc("/addword", addword)
    http.HandleFunc("/exists", exists)
    http.HandleFunc("/", index)
    http.HandleFunc("/reload", reload)
    err := http.ListenAndServe(":" + *port, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```

## 测试demo
```
➜  /tmp cat words.txt
SB|SAD|FUCK
➜  /tmp curl http://localhost:8080/

    请求格式：
        1. http://ip:port/
            首页
        2. http://ip:port/addword?word=xxx
            添加违禁词
        3. http://ip:port/reload?filepath=xxx
            重新载入
        4. http://ip:port/filter?text=xxx&replace=xxx
            检测是否有违禁词，并进行替换
➜  /tmp curl "http://localhost:8080/addword?word=傻瓜"
单词【傻瓜】已添加到词库
➜  /tmp curl "http://localhost:8080/exists?word=傻瓜"
单词【傻瓜】在词库中
➜  /tmp curl "http://localhost:8080/exists?word=你好"
单词【你好】不在词库中
➜  /tmp curl "http://localhost:8080/reload?filepath=/tmp/words.txt&separator=|"
新的TrieTree已重建
➜  /tmp curl "http://localhost:8080/exists?word=SAD"
单词【SAD】在词库中
➜  /tmp
➜  /tmp
➜  /tmp
➜  /tmp curl "http://localhost:8080/filter?text=别SAD了，你妈喊你回家吃饭了，小傻瓜&replace=***"
hit, final result: 别***了，你妈喊你回家吃饭了，小傻瓜
➜  /tmp curl "http://localhost:8080/exists?word=傻瓜"
单词【傻瓜】不在词库中
➜  /tmp curl "http://localhost:8080/addword?word=傻瓜"
单词【傻瓜】已添加到词库
➜  /tmp curl "http://localhost:8080/filter?text=别SAD了，你妈喊你回家吃饭了，小傻瓜&replace=***"
hit, final result: 别***了，你妈喊你回家吃饭了，小***
```