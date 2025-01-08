package lbser

import "github.com/shopspring/decimal"

type LBSer interface {
	ReGeo() (*Address, error)
	SetLocation(lon, lat decimal.Decimal) LBSer
	SetApi(api string) LBSer
}

type Provider string

type Conf struct {
	Provider Provider
	Key      string
	Api      string
}

type Address struct {
	Province string
	City     string
	Region   string
}
