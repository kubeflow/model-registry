package integrations

import (
	"bytes"
	"io"
	"net/http"
)

func StreamToString(stream io.Reader) string {
	if stream == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return ""
	}
	return buf.String()
}

func CloneBody(r *http.Request) ([]byte, error) {
	buf, _ := io.ReadAll(r.Body)
	rdr1 := io.NopCloser(bytes.NewBuffer(buf))
	rdr2 := io.NopCloser(bytes.NewBuffer(buf))
	r.Body = rdr2 // OK since rdr2 implements the io.ReadCloser interface

	defer rdr1.Close()
	cloneBody, err := io.ReadAll(rdr2)

	return cloneBody, err
}
