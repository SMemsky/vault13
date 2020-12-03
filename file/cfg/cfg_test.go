package cfg

import (
	"strings"
	"testing"
)

const test1 = `

foobar = 10015

[hello]
integer  = 100
floating = 100.3
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

	if value, err := config.Get("hello", "integer").Int() ; value != 100 || err != nil {
		t.Errorf("Get integer failed %s", config.Get("hello", "integer").String())
	}

	if value, err := config.Get("hello", "floating").Float() ; value != 100.3 || err != nil {
		t.Errorf("Get floating failed")
	}

	if value := config.Get("hello", "string").String() ; value != "Some text here" {
		t.Errorf("Get string failed")
	}

	if value, err := config.Get("unknown", "foobar").Int() ; value != 10015 || err != nil {
		t.Errorf("Get unknown failed")
	}
}
