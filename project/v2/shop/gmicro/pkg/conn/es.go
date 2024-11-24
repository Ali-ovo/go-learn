package conn

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

type EsOptions struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEsClient(opts *EsOptions) (*elasticsearch.Client, error) {
	host := fmt.Sprintf("http://%s:%d", opts.Host, opts.Port)
	return elasticsearch.NewClient(
		elasticsearch.Config{
			//Logger: &elastictransport.ColorLogger{
			//	Output:             os.Stdout,
			//	EnableRequestBody:  false,
			//	EnableResponseBody: false,
			//},
			Addresses: []string{host},
			Username:  opts.Username,
			Password:  opts.Password,
		},
	)
}
