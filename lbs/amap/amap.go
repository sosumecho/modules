package amap

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"github.com/sosumecho/modules/lbs/lbser"
)

const (
	Name lbser.Provider = "amap"
	url                 = "https://restapi.amap.com/v3/geocode/regeo"
)

type Amap struct {
	conf     *lbser.Conf
	location string
}

func (a Amap) ReGeo() (*lbser.Address, error) {
	client := resty.New()
	var rs ReGeoResponse
	_, err := client.R().SetQueryParams(map[string]string{
		"key":      a.conf.Key,
		"location": a.location,
	}).SetResult(&rs).Get(a.conf.Api)
	if err != nil {
		return nil, err
	}
	return &lbser.Address{
		Province: rs.ReGeoCode.AddressComponent.Province,
		City:     rs.ReGeoCode.AddressComponent.City,
		Region:   rs.ReGeoCode.AddressComponent.Township,
	}, nil
}

func (a Amap) SetLocation(lon, lat decimal.Decimal) lbser.LBSer {
	a.location = fmt.Sprintf("%s,%s", lon.String(), lat.String())
	return a
}

func (a Amap) SetApi(api string) lbser.LBSer {
	a.conf.Api = api
	return a
}

func NewAmap(conf *lbser.Conf) *Amap {
	if conf.Api == "" {
		conf.Api = url
	}
	return &Amap{
		conf: conf,
	}
}

type ReGeoResponse struct {
	Status    string    `json:"status"`
	ReGeoCode ReGeoCode `json:"regeocode"`
	InfoCode  string    `json:"infocode"`
	Info      string    `json:"info"`
}

type ReGeoCode struct {
	AddressComponent struct {
		City     string `json:"city"`
		Province string `json:"province"`
		District string `json:"district"`
		Country  string `json:"country"`
		Township string `json:"township"`
	} `json:"addressComponent"`
	FormattedAddress string `json:"formatted_address"`
}
