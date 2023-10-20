package wafflepaw

import (
	"fmt"
	"time"

	"github.com/infinitybotlist/sysmanage-web/plugins/wafflepaw/types"
)

func startWafflepawCluster() {
	for {
		initCluster()

		fmt.Println("wafflepaw: restarting cluster due to early return!!!")
	}
}

func initCluster() {
	for i := range config {
		startProjectMonitor(config[i])
	}
}

func startProjectMonitor(c *types.Project) {
}

func StartMonitor(c *types.Component) {
	if c.Every == 0 {
		c.Every = 30 * time.Second
	}

	//ticker := time.NewTicker(c.Every)

}
