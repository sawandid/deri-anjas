package miner

import (
	"fmt"
	"runtime"

	"github.com/go-logr/logr"
	"github.com/whalesburg/dero-stratum-miner/internal/version"
)

func (c *Client) setLogger(logger logr.Logger) {
	c.logger.Info("OK")
}
