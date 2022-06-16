package gofound

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type GOFound struct {
	addr string
}

func NewGOFound(addr string) *GOFound {
	return &GOFound{addr: addr}
}

func (g *GOFound) request(api string, method string, body []byte) ([]byte, error) {
	url := g.addr + "/api/" + api
	payload := bytes.NewReader(body)
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func (g *GOFound) PutIndex(db string, id uint32, text string, document interface{}) error {
	b, _ := json.Marshal(&PutIndexRequest{
		ID:       id,
		Text:     text,
		Document: document,
	})
	respBody, err := g.request("index?database="+db, http.MethodPost, b)
	if err != nil {
		return err
	}
	resp := &PutIndexResponse{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return err
	}
	if resp.State {
		return nil
	}
	return errors.New(resp.Message)
}

func (g *GOFound) QueryIndex(db string, queryParam *QueryIndexRequest) (*QueryIndexInfo[interface{}], error) {
	return QueryIndex[interface{}](g, db, queryParam)
}

func QueryIndex[T any](g *GOFound, db string, queryParam *QueryIndexRequest) (*QueryIndexInfo[T], error) {
	b, _ := json.Marshal(queryParam)
	respBody, err := g.request("query?database="+db, http.MethodPost, b)
	if err != nil {
		return nil, err
	}
	resp := &QueryIndexResponse[T]{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}
	if resp.State {
		return resp.Data, nil
	}
	return nil, errors.New(resp.Message)
}

func (g *GOFound) RemoveIndex(db string, id uint32) error {
	b, _ := json.Marshal(&RemoveIndexRequest{ID: id})
	respBody, err := g.request("remove?database="+db, http.MethodPost, b)
	if err != nil {
		return err
	}
	resp := &RemoveIndexResponse{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return err
	}
	if resp.State {
		return nil
	}
	return errors.New(resp.Message)
}

func (g *GOFound) DropDatabase(db string) error {
	respBody, err := g.request("db/drop?database="+db, http.MethodGet, nil)
	if err != nil {
		return err
	}
	resp := &Response{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return err
	}
	if resp.State {
		return nil
	}
	return errors.New(resp.Message)
}
