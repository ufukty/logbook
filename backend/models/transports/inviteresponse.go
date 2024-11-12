package transports

type InviteResponse string

const (
	Accept = InviteResponse("accept")
	Reject = InviteResponse("reject")
)
