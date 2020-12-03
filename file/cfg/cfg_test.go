package cfg

import (
	"strings"
	"testing"
)

const test1 = `

[hello]
integer  = 100
floating = 100.0
string   = Some text here


`

func TestBasic1(t *testing.T) {
	config := New()
	config.Set("hello", "integer", Int(-13))

	if config.Get("hello", "integer").String() != "-13" {
		t.Errorf("Get as string failed")
	}

	config.Set("boi", "bar", String("100500"))
	if value, err := config.Get("boi", "bar").Int() ; value != 100500 || err != nil {
		t.Errorf("Get as int failed")
	}

	if config.Parse(strings.NewReader(test1)) != nil {
		t.Errorf("Parse failed")
	}
}
