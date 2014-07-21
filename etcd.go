package etcd

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var (
	EtcdLookupFailure = errors.New("Failed to communicate with etcd")
)

type EtcdKeyNode struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	ModifiedIndex int    `json:"modifiedIndex"`
	CreatedIndex  int    `json:"createdIndex"`
}

type EtcdKeyResponse struct {
	Action   string       `json:"action"`
	Node     EtcdKeyNode  `json:"node"`
	PrevNode *EtcdKeyNode `json:"prevNode"`
}

type Etcd struct {
	host string
}

func (e Etcd) GetValue(key, defaultValue string) (string, error) {
	response, err := http.Get(e.host + "/v2/keys/" + key)
	if err != nil {
		return "", EtcdLookupFailure
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return defaultValue, nil
	}

	keyResponse := &EtcdKeyResponse{}
	if err := json.NewDecoder(response.Body).Decode(keyResponse); err != nil {
		return "", EtcdLookupFailure
	}

	return keyResponse.Node.Value, nil
}

func (e Etcd) SetValue(key, value string) error {
	urlValues := url.Values{"value": {value}}

	request, err := http.NewRequest("PUT", e.host+"/v2/keys/"+key, strings.NewReader(urlValues.Encode()))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if _, err := http.DefaultClient.Do(request); err != nil {
		return err
	}

	return nil
}

func NewEtcd(host string) Etcd {
	return Etcd{
		host: host,
	}
}
