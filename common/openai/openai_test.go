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

func TestListApiBySchema(t *testing.T) {
	token, err := os.ReadFile("./testdata/openai_token.txt")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	schema, err := os.ReadFile("./testdata/customer_schema.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	o := NewOpenAI(string(token), "en")
	o.SetMaxTokens(3000)
	res, err := o.ListApiBySchema("Customer", string(schema))
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}

func TestCreateApiBySchema(t *testing.T) {
	token, err := os.ReadFile("./testdata/openai_token.txt")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	schema, err := os.ReadFile("./testdata/customer_schema.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	o := NewOpenAI(string(token), "en")
	o.SetMaxTokens(3000)
	res, err := o.CreateApiBySchema("CreateCustomer", string(schema))
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}
