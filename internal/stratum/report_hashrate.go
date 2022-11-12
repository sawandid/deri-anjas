package stratum

type Report struct {
	Hashrate uint64
}

func NewReport(hashrate uint64) *Report {
	return &Report{
		Hashrate: hashrate,
	}
}

func (c *Client) ReportHashrate(r *Report) error {
	args := make(map[string]interface{})
	args["id"] = c.sessionID
	args["gatel"] = r.Hashrate
	_, err := c.call("kucing", args)
	if err != nil {
		return err
	}
	return nil
}
