package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

type attach struct {
	filename    string
	contentType string
	body        io.ReadCloser
	// alreayBase64 bool
}

type Message struct {
	// 主题
	Subject string
	// 内容
	Body string
	//
	From mail.Address
	// 接收者
	To string
	// 附件
	attrs []attach
	// 内容类型
	BodyType string
}

// AddAttach 添加附件
func (m *Message) AddAttach(fname string, r io.Reader) {
	rc, ok := r.(io.ReadCloser)
	if !ok {
		rc = io.NopCloser(r)
	}
	ctype := mime.TypeByExtension(filepath.Ext(fname))
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	m.attrs = append(m.attrs, attach{
		filename:    mime.BEncoding.Encode("UTF-8", fname),
		contentType: ctype,
		body:        rc,
	})
}

func (m *Message) base64data(buf *bytes.Buffer, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b, data)
	for i, l := 0, len(b); i < l; i++ {
		buf.WriteByte(b[i])
		if (i+1)%76 == 0 {
			buf.WriteString("\r\n")
		}
	}
	return nil
}

func (m *Message) Bytes() ([]byte, error) {

	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)

	// 公共头
	baseHeader := make(textproto.MIMEHeader)
	baseHeader.Set("From", m.From.String())
	baseHeader.Set("To", m.To)
	baseHeader.Set("Subject", mime.BEncoding.Encode("UTF-8", m.Subject))
	baseHeader.Set("MIME-Version", "1.0")
	baseHeader.Set("Content-Type", fmt.Sprintf(`multipart/mixed; boundary="%s"`, w.Boundary()))
	baseHeader.Set("Date", time.Now().Format(time.RFC1123Z))
	baseHeader.Set("Message-ID", fmt.Sprintf("<id_%d.ums>", time.Now().UnixNano()))
	for k, v := range baseHeader {
		fmt.Fprintf(buf, "%s: %s\r\n", k, v[0])
	}
	fmt.Fprintf(buf, "\r\n")

	// 内容
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", m.BodyType+"; charset=utf-8")
	header.Set("Content-Transfer-Encoding", "base64")
	if _, err := w.CreatePart(header); err != nil {
		return nil, err
	}
	if err := m.base64data(buf, strings.NewReader(m.Body)); err != nil {
		return nil, err
	}

	// 添加附件
	for _, att := range m.attrs {
		var fheader = make(textproto.MIMEHeader)
		fheader.Set("Content-Transfer-Encoding", "base64")
		fheader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, att.filename))
		fheader.Set("Content-Type", fmt.Sprintf(`%s; name="%s"`, att.contentType, att.filename))
		if _, err := w.CreatePart(fheader); err != nil {
			return nil, err
		}
		if err := m.base64data(buf, att.body); err == nil {
			att.body.Close()
		}
	}
	w.Close()
	return buf.Bytes(), nil
}
