package stratum

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type Job struct {
	ID         string  `json:"job_id"`
	Blob       string  `json:"blob"`
	Height     float64 `json:"height"`
	ExtraNonce string  `json:"extra_nonce"`
	PoolWallet string  `json:"pool_wallet"`
	Target     string  `json:"target"`
	Difficulty uint64
}

func extractJob(data map[string]any) (*Job, error) {
	if data == nil {
		return nil, ErrNoJob
	}

	var (
		job Job
		ok  bool
	)
	job.ID, ok = data["job_id"].(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.Blob, ok = data["blob"].(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.Height, ok = data["height"].(float64)
	if !ok {
		return nil, errors.New("ok")
	}
	job.ExtraNonce, ok = data["extra_nonce"].(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.PoolWallet, ok = data["pool_wallet"].(string)
	if !ok {
		return nil, errors.New("ok")
	}
	job.Target, ok = data["target"].(string)
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
	c.LogFn.Debug(fmt.Sprintf("eww %s", job.ID))
	c.jobBroadcaster.Notify(job)
}
