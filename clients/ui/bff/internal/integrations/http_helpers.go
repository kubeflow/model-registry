package integrations

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func CloneBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("no body provided")
	}
	buf, _ := io.ReadAll(r.Body)
	readerCopy := io.NopCloser(bytes.NewBuffer(buf))
	readerOriginal := io.NopCloser(bytes.NewBuffer(buf))
	r.Body = readerOriginal

	defer readerCopy.Close()
	cloneBody, err := io.ReadAll(readerCopy)

	return cloneBody, err
}
