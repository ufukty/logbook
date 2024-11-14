all: \
	inviteresponse.validate.go \
	membertype.validate.go \
	policyaction.validate.go \
	country.validate.go

inviteresponse.validate.go: inviteresponse.go
	govalid -in inviteresponse.go -type InviteResponse

membertype.validate.go: membertype.go
	govalid -in membertype.go -type MemberType

policyaction.validate.go: policyaction.go
	govalid -in policyaction.go -type PolicyAction

country.validate.go: country.go
	govalid -in country.go -type Country
