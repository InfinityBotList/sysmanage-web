package nginx

import (
	"fmt"
	"sysmanage-web/core/state"
)

var zoneMap = make(map[string]string)

func setupCf() {
	// List zpnes
	zones, err := cf.ListZones(state.Context)

	if err != nil {
		panic(err)
	}

	for _, zone := range zones {
		fmt.Println("CF: Zone added =>", zone.Name, "("+zone.ID+")")
		zoneMap[zone.Name] = zone.ID
	}
}
