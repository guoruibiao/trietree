package trietree

import (
	"strings"
	"os"
	"io/ioutil"
	"sync"
	)
// 参考链接
// https://www.cnblogs.com/sunlong88/p/11980046.html

type trieNode struct {
	char rune
	Data interface{}
	parent *trieNode
	Depth int
	childs map[rune]*trieNode
	term bool
	// TODO 考虑加上hooks 可以在某些情况下做一些回调处理
}

type TrieTree struct {
	root *trieNode
	size int
	lock *sync.RWMutex
	// 做一些统计之类的事
	statistic map[string]int
}

func newNode() *trieNode {
	return &trieNode{
		childs:make(map[rune]*trieNode, 32),
	}
}

func NewTrie() *TrieTree {
	return &TrieTree{
		root:newNode(),
		lock:new(sync.RWMutex),
	}
}

func (p *TrieTree) BuildTrieTreeFromFile(filepath, separator string) (err error) {
	// 拿到文件句柄
	fileHandler, err := os.Open(filepath)
	defer fileHandler.Close()
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(fileHandler)
	if err != nil {
		return
	}

	// 解析词库
	words := strings.Split(string(data), separator)
	if len(words) >= 1 {
		for _, word := range words {
			// 排除尾行空串的影响，不过对发言内容来说，空内容其实是不合法的
			word = strings.TrimSpace(word)
			if word != "" {
				p.AddWord(word, nil)
			}
		}
	}

	return
}

func (p *TrieTree) AddWord(key string, data interface{}) (err error) {
	// 只为写操作加锁处理，读操作原则上来不加锁最多影响点精度
	p.lock.Lock()
	defer p.lock.Unlock()

	key = strings.TrimSpace(key)
	node := p.root
	runes := []rune(key)

	// 逐个rune进行添加
	for _, r := range runes {
		ret, ok := node.childs[r]
		if !ok {
			ret = newNode()
			ret.Depth = ret.Depth + 1
			ret.char = r
			node.childs[r] = ret
		}

		node = ret
	}

	node.term = true
	node.Data = data
	return
}

func (p *TrieTree) WordExists(key string) (exists bool) {
	key = strings.TrimSpace(key)
	chars := []rune(key)
	node := p.root

	for _, char := range chars {
		ret, ok := node.childs[char]
		if !ok {
			return
		}

		node = ret
	}
	return true
}

// 内部统计方法调用
func (p *TrieTree) incrWordFreq(key string) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	
	old, ok := p.statistic[key]
	if !ok {
		p.statistic[key] = 1
		return 1
	}
	
	p.statistic[key] = p.statistic[key] + 1
	return old + 1
}

func (p *TrieTree) ExportStatistic() map[string]int {
	return p.statistic
}

func (p *TrieTree)Filter(text, replace string) (result string, hit bool) {
	chars := []rune(text)
	replace = strings.TrimSpace(replace)

	if p.root == nil {
		return
	}

	var left []rune
	node := p.root
	start := 0

	for index, char := range chars {
		ret, ok := node.childs[char]
		if !ok {
			left = append(left, chars[start:index+1]...)
			start = index + 1
			node = p.root
			continue
		}

		node = ret

		// start 到 index+1 部分刚好为遍历过程走过的长度，用于替换即可
		if node.term == true {
			left = append(left, []rune(replace)...)
			hit = true
			node = p.root
			start = index + 1
			continue
		}
	}
	return string(left), hit
}
