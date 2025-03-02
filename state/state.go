package state

type Handle string

type Responder func(string)

type Online struct {
	User    Handle
	Respond Responder
}

type Offline struct {
	User Handle
}

type Request struct {
	User Handle
	Text string
}

type Abandon struct {
	User Handle
}
