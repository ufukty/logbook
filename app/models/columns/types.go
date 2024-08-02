package columns

// commons
type (
	NonNegativeNumber int
)

// accounts
type (
	AccessId     string
	Email        string
	HumanName    string
	LoginId      string
	SessionId    string
	SessionToken string
	UserId       string
	Username     string
)

// objectives
type (
	LinkId      string
	ObjectiveId string
	OperationId string
	VersionId   string
)

// tags
type (
	LinkType string
	TagId    string
)

const (
	Primary = LinkType("PRIMARY") // eg. When task owner break downs it
	Private = LinkType("PRIVATE") //
	Remote  = LinkType("REMOTE")  // eg. Collaborated objective attached to local objectives
)
