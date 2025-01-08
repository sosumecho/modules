package lbs

import (
	"github.com/sosumecho/modules/exceptions"
	"github.com/sosumecho/modules/lbs/amap"
	"github.com/sosumecho/modules/lbs/lbser"
	"github.com/sosumecho/modules/lbs/nominatim"
)

func New(conf *lbser.Conf) (lbser.LBSer, error) {
	switch conf.Provider {
	case amap.Name:
		return amap.NewAmap(conf), nil
	case nominatim.Name:
		return nominatim.NewNominatim(conf), nil
	}
	return nil, exceptions.ParamsError
}
