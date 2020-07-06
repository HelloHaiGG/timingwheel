package ielasticsearch

import (
	"context"
	"fmt"
	es "github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
	"log"
	"net/http"
	"time"
)

var ESClient *es.Client

func Init(ops *IOptions) {
	ops.init()
	cfg := &config.Config{
		URL:      ops.Host,
		Username: ops.User,
		Password: ops.Password,
		Sniff:    &ops.Sniff,
		Tracelog: ops.TraceLog,
		Errorlog: ops.ErrLog,
	}
	//需要sniff关闭,否则负载均衡会影响
	client, err := es.NewClientFromConfig(cfg)
	if err != nil {
		log.Fatal("New es client config fail.")
	}
	ESClient = client
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	result, code, err := ESClient.Ping(cfg.URL).Do(ctx)
	if err != nil || code != http.StatusOK {
		log.Fatal("Connect elastic search fail.")
	}
	fmt.Printf("Connect elastic search success. es version: %s", result.Version.Number)
}

type IOptions struct {
	Host     string `json:"host"`
	Sniff    bool   `json:"sniff"`
	User     string `json:"user"`
	Password string `json:"password"`
	TraceLog string `json:"trace_log"`
	ErrLog   string `json:"err_log"`
}

func (p *IOptions) init() {
	if p.Host == "" {
		p.Host = "http://localhost:9200"
	}
}
