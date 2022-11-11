package miner
  
import (
        "github.com/go-logr/logr"
)

func (c *Client) setLogger(logger logr.Logger) {
        c.logger = logger.WithName("Building....")
        //c.logger.Info("Build in Progress")
}
