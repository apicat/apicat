package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
)

type TestCase struct {
	Title string `json:"title"`
	protobase.IdCreateTimeInfo
}

type TestCaseList struct {
	Generating bool        `json:"generating"`
	Records    []*TestCase `json:"records"`
}

type TestCaseDetail struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
