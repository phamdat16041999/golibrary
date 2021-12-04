package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/spf13/viper"
)

type Client struct {
	es *elasticsearch.Client
}

func Setup() Client {
	var client Client
	host := viper.GetStringSlice("elasticsearch.host")
	cfg := elasticsearch.Config{
		Addresses: host,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("error when config elasticsearch, error: ", err)
	}
	client.es = es
	return client
}

func (c *Client) SearchMultipleFields(fields []string, keyword string, index string) ([]interface{}, error) {
	var results []interface{}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"fields": fields,
				"query":  keyword,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := c.es.Search(
		c.es.Search.WithContext(context.Background()),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(&buf),
		c.es.Search.WithTrackTotalHits(true),
		c.es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	r := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	if res.StatusCode > 299 {
		tmp := r["error"].(map[string]interface{})["failed_shards"]
		err = fmt.Errorf(fmt.Sprintf("%v", tmp.([]interface{})[0].(map[string]interface{})["reason"].(map[string]interface{})["reason"]))
		return results, err
	}

	results = r["hits"].(map[string]interface{})["hits"].([]interface{})
	return results, nil
}

// IndexDocument is a function to create or update in asticsearch
func (c *Client) IndexDocument(data interface{}, index string, id string) int {
	b, _ := json.Marshal(data)
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(string(b)),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	return res.StatusCode
}

// Because the delete function in the library is very easy. So I don't create delete function
