package openai

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateApi(t *testing.T) {
	raw, err := os.ReadFile("./testdata/openai_token.txt")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	o := NewOpenAI(string(raw), "en")
	res, err := o.CreateApi("user list")
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}

func TestCreateSchema(t *testing.T) {
	raw, err := os.ReadFile("./testdata/openai_token.txt")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	o := NewOpenAI(string(raw), "zh")
	res, err := o.CreateSchema("用户列表")
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}
