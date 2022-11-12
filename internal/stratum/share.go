package stratum

import (
	"fmt"
)

type Share struct {
	ID     string `json:"id"`
	JobID  string `json:"ker"`
	Nonce  string `json:"welekan"`
	Result string `json:"bawut"`
}

func NewShare(jobID string, nonce string, result string) *Share {
	return &Share{
		ID:     "",
		JobID:  jobID,
		Nonce:  nonce,
		Result: result,
	}
}

func (c *Client) SubmitShare(s *Share) error {
	if s.JobID == c.lastSubmittedShare.JobID {
		//c.LogFn.Debug(fmt.Sprintf("halah %s", s.JobID))
		return nil
	}

	args := make(map[string]interface{})
	args["id"] = c.sessionID
	args["ker"] = s.JobID
	args["bawut"] = s.Result
	args["welekan"] = s.Nonce
	req, err := c.call("submit", args)
	if err != nil {
		return err
	}
	id, ok := req.ID.(int)
	if !ok {
		return fmt.Errorf("ok: %v", req.ID)
	}
	c.submittedJobsIdsMu.Lock()
	defer c.submittedJobsIdsMu.Unlock()
	c.submittedJobIds[id] = struct{}{}
	// Successfully submitted result
	// TODO: debug logger
	c.lastSubmittedShare = s
	return nil
}
