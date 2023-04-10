package ac

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestAc(t *testing.T) {
	r, _ := regexp.Compile("([\u000E-\u02FF]+)")
	tests := struct {
		name     string
		words    []WordInfos
		contents map[string]string
		want     map[string]*[]AddrWords
		wantErr  bool
	}{
		name: "1",
		words: []WordInfos{
			{Word: "好cool", Tags: []string{"brand"}},
			{Word: "洗面", Tags: []string{"brand"}},
			{Word: "氨基酸洗面奶", Tags: []string{"brand"}},
			{Word: "洗面奶", Tags: []string{"brand"}},
			{Word: "面奶", Tags: []string{"brand"}},
			{Word: "基酸洗面", Tags: []string{"brand"}},
		},
		contents: map[string]string{
			"1": "你真的好cool",
			"2": "你真的好cooler",
			"3": "氨基酸洗面奶",
		},
		want: map[string]*[]AddrWords{
			"1": {
				{Addr: 9, Words: []WordInfos{{Word: "好 cool ", Tags: []string{"brand"}}}},
			},
			"2": {},
			"3": {
				{Addr: 4, Words: []WordInfos{
					{Word: "基酸洗面", Tags: []string{"brand"}},
					{Word: "洗面", Tags: []string{"brand"}},
				}},
				{Addr: 5, Words: []WordInfos{
					{Word: "氨基酸洗面奶", Tags: []string{"brand"}},
					{Word: "洗面奶", Tags: []string{"brand"}},
					{Word: "面奶", Tags: []string{"brand"}},
				}},
			},
		},
		wantErr: false,
	}
	formatWords := make([]WordInfos, 0)
	for _, wordInfo := range tests.words {
		word := r.ReplaceAllStringFunc(wordInfo.Word, func(s string) string {
			return fmt.Sprintf(" %s ", s)
		})
		formatWords = append(formatWords, WordInfos{Word: word, Tags: wordInfo.Tags})
	}
	acLoadMatcher := ReLoadNewMatcher()
	acLoadMatcher.Build(formatWords)
	AcMatcher = acLoadMatcher

	matcher := NewMatcher()
	mss := make(map[string]*[]AddrWords, 0)
	for key, content := range tests.contents {
		content = r.ReplaceAllStringFunc(content, func(s string) string {
			return fmt.Sprintf(" %s ", s)
		})
		ms := matcher.MatchMany(content)
		mss[key] = &ms
	}
	tests.wantErr = true
	for key, v := range mss {
		if !reflect.DeepEqual(tests.want[key], v) {
			t.Errorf("error : key:%s", key)
		}
	}
}
