package utils

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"

	"net"
)

// IP2City ip转城市
func IP2City(ip string) string {
	file := fmt.Sprintf("%s/%s", GetAbsDir(), "public/data/GeoLite2-City.mmdb")
	reader, err := geoip2.Open(file)
	if err != nil {
		return ""
	}
	defer reader.Close()
	record, err := reader.City(net.ParseIP(ip))
	if err != nil {
		return ""
	}
	country := record.Country.Names["zh-CN"]
	return country + " " + record.City.Names["zh-CN"]
}

// IP2Country ip转国家
func IP2Country(ip string) string {
	file := fmt.Sprintf("%s/%s", GetAbsDir(), "public/data/GeoLite2-City.mmdb")
	reader, err := geoip2.Open(file)
	if err != nil {
		return ""
	}
	defer reader.Close()
	record, err := reader.Country(net.ParseIP(ip))
	if err != nil {
		return ""
	}
	country := record.Country.Names["zh-CN"]
	return country
}

func IP2Location(ip string) (float64, float64) {
	file := fmt.Sprintf("%s/%s", GetAbsDir(), "public/data/GeoLite2-City.mmdb")
	reader, err := geoip2.Open(file)
	if err != nil {
		return 0, 0
	}
	defer reader.Close()
	record, err := reader.City(net.ParseIP(ip))
	if err != nil {
		return 0, 0
	}

	return record.Location.Longitude, record.Location.Latitude

}

func IP2CountryAndCity(ip string) (string, string) {
	file := fmt.Sprintf("%s/%s", GetAbsDir(), "public/data/GeoLite2-City.mmdb")
	reader, err := geoip2.Open(file)
	if err != nil {
		return "", ""
	}
	defer reader.Close()
	record, err := reader.City(net.ParseIP(ip))
	if err != nil {
		return "", ""
	}
	return record.Country.Names["zh-CN"], record.City.Names["zh-CN"]
}
