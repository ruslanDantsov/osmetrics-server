package client

import "github.com/go-resty/resty/v2"

type RestClient interface {
	R() *resty.Request
}
