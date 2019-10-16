package openapi

import (
	"bytes"
	"errors"
	"free5gc/lib/openapi/common"
	"io/ioutil"
	"net/http"
	"regexp"
)

type MultipartRelatedBinding struct{}

func (MultipartRelatedBinding) Name() string {
	return "json"
}

func (MultipartRelatedBinding) Bind(req *http.Request, obj interface{}) error {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return common.Decode(obj, b, req.Header.Get("Content-Type"))
}

func (MultipartRelatedBinding) BindBody(body []byte, obj interface{}) error {
	re, _ := regexp.Compile(`--([a-zA-Z0-9+\-_]+)--`)
	submatch := re.FindSubmatch(body)
	if len(submatch) < 1 {
		return errors.New("cannot parse multipart boundary")
	}
	return common.Decode(obj, body, "multipart/related; boundary="+string(submatch[1]))
}

type MultipartRelatedRender struct {
	Data        interface{}
	contentType string
}

func (r MultipartRelatedRender) Render(w http.ResponseWriter) (err error) {
	payloadBuf := &bytes.Buffer{}
	ct, err := common.MultipartEncode(r.Data, payloadBuf)
	if err != nil {
		panic(err)
	}
	r.contentType = ct
	w.Header().Set("Content-Type", r.contentType)
	_, err = payloadBuf.WriteTo(w)
	if err != nil {
		panic(err)
	}
	return
}

func (r MultipartRelatedRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", r.contentType)
}
