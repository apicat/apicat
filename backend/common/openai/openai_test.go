package openai

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/apicat/apicat/backend/config"
)

func loadOpenAIConfig() config.OpenAI {
	raw, err := os.ReadFile("./testdata/openai_config.json")
	if err != nil {
		panic(err)
	}

	var configs config.OpenAI
	err = json.Unmarshal(raw, &configs)
	if err != nil {
		panic(err)
	}

	return configs
}

func TestCreateApi(t *testing.T) {
	configs := loadOpenAIConfig()
	o := NewOpenAI(configs, "en")
	res, err := o.CreateApi("user list")
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}

func TestCreateSchema(t *testing.T) {
	configs := loadOpenAIConfig()
	o := NewOpenAI(configs, "zh")
	res, err := o.CreateSchema("用户列表")
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}

func TestListApiBySchema(t *testing.T) {
	configs := loadOpenAIConfig()
	o := NewOpenAI(configs, "en")
	o.SetMaxTokens(3000)
	res, err := o.ListApiBySchema("Customer")
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}

func TestCreateApiBySchema(t *testing.T) {
	configs := loadOpenAIConfig()

	schema, err := os.ReadFile("./testdata/customer_schema.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	o := NewOpenAI(configs, "en")
	o.SetMaxTokens(3000)
	res, err := o.CreateApiBySchema("CreateCustomer", "/customers", "POST", string(schema))
	if err != nil || res == "" {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(res)
}
