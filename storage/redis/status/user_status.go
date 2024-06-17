package status

type Status int

const (
	Unknown Status = iota
	Exhausted
	Free
	Unlimited
	Premium
)
