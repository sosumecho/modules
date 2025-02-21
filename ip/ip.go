package ip

import (
	"github.com/oschwald/geoip2-golang"
	"net"
	"sync"
)

var (
	once  sync.Once
	geoIP *GeoIP
)

type GeoIP struct {
	reader *geoip2.Reader
}

func (g *GeoIP) City(ip string) string {
	record, err := g.reader.City(net.ParseIP(ip))
	if err != nil {
		return ""
	}
	country := record.Country.Names["zh-CN"]
	if country == "" || record.City.Names["zh-CN"] == "" {
		return ""
	}
	return country + " " + record.City.Names["zh-CN"]
}

func (g *GeoIP) CityRaw(ip string) string {
	record, err := g.reader.City(net.ParseIP(ip))
	if err != nil {
		return ""
	}
	return record.City.Names["zh-CN"]
}

func (g *GeoIP) Country(ip string) (string, string) {
	record, err := g.reader.City(net.ParseIP(ip))
	if err != nil {
		return "", ""
	}
	return record.Country.Names["zh-CN"], record.Country.IsoCode
}

func (g *GeoIP) Location(ip string) (float64, float64) {
	record, err := g.reader.City(net.ParseIP(ip))
	if err != nil {
		return 0, 0
	}
	return record.Location.Longitude, record.Location.Latitude
}

func (g *GeoIP) CountryAndCity(ip string) (string, string) {
	record, err := g.reader.City(net.ParseIP(ip))
	if err != nil {
		return "", ""
	}
	return record.Country.Names["zh-CN"], record.City.Names["zh-CN"]
}

// Continent 大陆
func (g *GeoIP) Continent(ip string) string {
	record, err := g.reader.City(net.ParseIP(ip))
	if err != nil {
		return ""
	}
	return record.Continent.Names["zh-CN"]
}

type Conf struct {
	Path string `mapstructure:"path"`
}

// NewGeoIP 实例化geo ip
func NewGeoIP(conf *Conf) *GeoIP {
	if geoIP == nil {
		once.Do(func() {
			reader, err := geoip2.Open(conf.Path)
			if err != nil {
				geoIP = nil
			} else {
				geoIP = &GeoIP{
					reader: reader,
				}
			}
		})
	}
	return geoIP
}
