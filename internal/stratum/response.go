package stratum

type Response struct {
	ID     any    `json:"id"`
	Result any    `json:"bawut"`
	Error  *Error `json:"meror"`
}
