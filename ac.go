package ac

import (
	"container/list"
)

type trieNode struct {
	count    int
	fail     *trieNode
	children map[rune]*trieNode
	words    []WordInfos
}

func newTrieNode() *trieNode {
	return &trieNode{
		count:    0,
		fail:     nil,
		children: make(map[rune]*trieNode),
		words:    []WordInfos{},
	}
}

type Matcher struct {
	root  *trieNode
	size  int
	mark  map[string]bool
	words []string
}

type WordInfos struct {
	Word string
	Tags []string
}

var AcMatcher = &Matcher{
	root:  newTrieNode(),
	size:  0,
	mark:  make(map[string]bool, 0),
	words: make([]string, 0),
}

func NewMatcher() *Matcher {
	return AcMatcher
}

func ReLoadNewMatcher() *Matcher {
	return &Matcher{
		root:  newTrieNode(),
		size:  0,
		mark:  make(map[string]bool, 0),
		words: make([]string, 0),
	}

}

// initialize the ahocorasick
func (this *Matcher) Build(dictionary []WordInfos) {
	words := make([]string, 0)
	words = append(words, this.words...)
	for i := range dictionary {
		this.words = append(this.words, dictionary[i].Word)
		words = append(words, dictionary[i].Word)
		this.insert(dictionary[i])
	}
	this.build()
	this.mark = make(map[string]bool, this.size)
}

func (this *Matcher) GetSize() (size int) {
	size = this.size
	return
}

func (this *Matcher) GetWords() []string {
	return this.words
}

// string match search
// return all strings matched as indexes into the original dictionary
func (this *Matcher) Match(s string) []string {
	curNode := this.root
	this.resetMark()
	var p *trieNode = nil
	ret := make([]string, 0)
	rs := []rune(s)
	flag := false
	for index := 0; index < len(rs); index++ {
		v := rs[index]

		for curNode.children[v] == nil && curNode != this.root {
			curNode = curNode.fail
		}
		curNode = curNode.children[v]
		if curNode == nil {
			curNode = this.root
		}
		p = curNode
		if flag && p == this.root {
			index--
			flag = false
		}
		for _, index := range p.words {
			for p != this.root && p.count > 0 && !this.mark[index.Word] {
				this.mark[index.Word] = true

				for i := 0; i < p.count; i++ {
					ret = append(ret, index.Word)
				}
				p = p.fail
				flag = true
			}
		}

	}

	return ret
}

type AddrWords struct {
	Addr  int
	Words []WordInfos
}

func (this *Matcher) MatchMany(s string) []AddrWords {
	curNode := this.root
	var p *trieNode = nil

	ret := make([]AddrWords, 0)
	ss := []rune(s)
	for index, v := range ss {
		for curNode.children[v] == nil && curNode != this.root {
			curNode = curNode.fail
		}
		curNode = curNode.children[v]
		if curNode == nil {
			curNode = this.root
		}
		p = curNode

		words := make([]WordInfos, 0)
		if p.count == 0 {
			p = p.fail
		}
		for p != nil && p.count > 0 {
			words = append(words, p.words...)
			p = p.fail
		}
		if len(words) == 0 {
			continue
		}
		ret = append(ret, AddrWords{
			Addr:  index,
			Words: words,
		})
	}

	return ret
}

// just return the number of len(Match(s))
func (this *Matcher) GetMatchResultSize(s string) int {

	curNode := this.root
	this.resetMark()
	var p *trieNode = nil

	num := 0

	for _, v := range s {
		for curNode.children[v] == nil && curNode != this.root {
			curNode = curNode.fail
		}
		curNode = curNode.children[v]
		if curNode == nil {
			curNode = this.root
		}

		p = curNode
		for _, word := range p.words {
			for p != this.root && p.count > 0 && !this.mark[word.Word] {
				this.mark[word.Word] = true
				num += p.count
				p = p.fail
			}
		}

	}

	return num
}

func (this *Matcher) build() {
	ll := list.New()
	// 后入 前出
	ll.PushBack(this.root)
	for ll.Len() > 0 {
		//从前面拿出数据并且从队列中删除
		temp := ll.Remove(ll.Front()).(*trieNode)
		var p *trieNode = nil
		for i, v := range temp.children {
			if temp == this.root {
				v.fail = this.root
			} else {
				p = temp.fail
				for p != nil {
					if p.children[i] != nil {
						v.fail = p.children[i]
						break
					}
					p = p.fail
				}

				if p == nil {
					v.fail = this.root
				}
			}
			ll.PushBack(v)
		}
	}

}

func (this *Matcher) insert(w WordInfos) {
	s := w.Word
	curNode := this.root
	// 按字数分解存入
	i := 0
	for _, v := range s {
		if curNode.children[v] == nil {
			curNode.children[v] = newTrieNode()
		}
		// 造成类似于递归插入的方式
		curNode = curNode.children[v]
		i++
		if len([]rune(s)) == i {
			curNode.words = append(curNode.words, w)
		}
	}
	// count 是用来确定此处是不是结束
	curNode.count++
	// 需要查找的词在数组中的未知
	this.size++
}

func (this *Matcher) resetMark() {
	this.mark = make(map[string]bool, this.size)
}
