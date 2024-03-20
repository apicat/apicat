package sendcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apicat/apicat/backend/module/mail/common"
)

type sendcloud struct {
	apiUser  string
	apiKey   string
	from     string
	fromName string
}

func NewSendCloud(cfg map[string]interface{}) (*sendcloud, error) {
	for _, v := range []string{"ApiUser", "ApiKey", "From", "FromName"} {
		if _, ok := cfg[v]; !ok {
			return nil, fmt.Errorf("sendcloud config %s is required", v)
		}
	}
	return &sendcloud{
		apiUser:  cfg["ApiUser"].(string),
		apiKey:   cfg["ApiKey"].(string),
		from:     cfg["From"].(string),
		fromName: cfg["FromName"].(string),
	}, nil
}

func (s *sendcloud) Send(msg *common.Message, to ...string) error {
	param := url.Values{}
	param.Add("to", strings.Join(to, ","))
	param.Add("subject", msg.Subject)
	param.Add("html", msg.Body)
	param.Add("apiUser", s.apiUser)
	param.Add("apiKey", s.apiKey)
	param.Add("from", s.from)
	param.Add("fromName", s.fromName)
	ctx, cfn := context.WithTimeout(context.Background(), time.Second*10)
	defer cfn()
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost,
		"https://api.sendcloud.net/apiv2/mail/send",
		strings.NewReader(param.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var ret Result
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return err
	}
	if ret.Ok {
		return nil
	}
	return ret
}

type Result struct {
	Ok         bool   `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (r Result) Error() string {
	return fmt.Sprintf("[%d] %s", r.StatusCode, r.Message)
}
