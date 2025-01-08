package nominatim

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"github.com/sosumecho/modules/exceptions"
	"github.com/sosumecho/modules/lbs/lbser"
)

const Name lbser.Provider = "nominatim"

type Nominatim struct {
	lon  decimal.Decimal
	lat  decimal.Decimal
	conf *lbser.Conf
}

func NewNominatim(conf *lbser.Conf) *Nominatim {
	return &Nominatim{conf: conf}
}

type ReGeoResponse struct {
	Type     string `json:"type"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Geocoding struct {
				Label string `json:"label"`
				Admin struct {
					Province string `json:"level4"`
					City     string `json:"level5"`
					Region   string `json:"level6"`
				} `json:"admin"`
			} `json:"geocoding"`
		} `json:"properties"`
	} `json:"features"`
}

func (n *Nominatim) ReGeo() (*lbser.Address, error) {
	client := resty.New()
	var rs ReGeoResponse
	_, err := client.R().
		SetQueryParams(map[string]string{
			"lat":    n.lat.String(),
			"lon":    n.lon.String(),
			"format": "geocodejson",
			"zoom":   "8",
		}).
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6").
		SetResult(&rs).Get(fmt.Sprintf("%s/reverse", n.conf.Api))
	if err != nil {
		return nil, err
	}
	if len(rs.Features) == 0 {
		return nil, exceptions.ParamsError
	}
	admin := rs.Features[0].Properties.Geocoding.Admin
	if admin.City == "" {
		admin.City = admin.Province
	}
	return &lbser.Address{
		Province: admin.Province,
		City:     admin.City,
		Region:   admin.Region,
	}, nil
}

func (n *Nominatim) SetLocation(lon, lat decimal.Decimal) lbser.LBSer {
	n.lon = lon
	n.lat = lat
	return n
}

func (n *Nominatim) SetApi(api string) lbser.LBSer {
	n.conf.Api = api
	return n
}
