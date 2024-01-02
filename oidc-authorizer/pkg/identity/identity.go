package identity

import (
	"encoding/base64"
	"encoding/json"
)

type Header struct {
	User string `json:"user"`
}

func (h *Header) Base64() (string, error) {
	data, err := json.Marshal(h)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func FromBase64(data string) (*Header, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var header Header
	err = json.Unmarshal(jsonBytes, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}
