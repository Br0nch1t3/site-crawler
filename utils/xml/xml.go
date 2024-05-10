package utilsxml

import (
	"bytes"
	"encoding/xml"
)

// Appends xml header to decoded xml bytes
func WithHeader(decoded []byte) []byte {
	return bytes.Join([][]byte{[]byte(xml.Header), decoded}, []byte{})
}
