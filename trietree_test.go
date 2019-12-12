package trietree

import (
	"testing"
	"fmt"
	)

var trietree *TrieTree

func init() {
	trietree = NewTrie()
	// 预定义几个违禁词
	trietree.AddWord("尼玛", nil)
	trietree.AddWord("傻瓜", nil)
}


func TestTrieTrie_BuildFromFile(t *testing.T) {
	filepath := "/tmp/words.txt"
	err := trietree.BuildTrieTreeFromFile(filepath, "\n")
	if err != nil {
		t.Error(err)
	}
	t.Log(trietree.root.childs)
}


func TestTrieTree_WordExists(t *testing.T) {
	testcases := []string{"傻瓜", "你好", "早上好"}
	var exists bool
	for _, testcase := range testcases {
		exists = trietree.WordExists(testcase)
		fmt.Println("Word: " + testcase + ", exists: ", exists)
	}
}


func TestTrieTree_Filter(t *testing.T) {
	TestTrieTrie_BuildFromFile(t)
	result, hit := trietree.Filter("你是SB吗？傻瓜", "***")
	if hit != true {
		t.Error(result)
	}
	t.Log(result)
}