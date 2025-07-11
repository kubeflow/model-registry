package kubernetes

type ServiceDetails struct {
	Name                string
	DisplayName         string
	Description         string
	ClusterIP           string
	HTTPPort            int32
	IsHTTPS             bool
	ExternalAddressRest string
}

type RequestIdentity struct {
	UserID string
	Groups []string
	Token  string
}

type BearerToken struct {
	raw string
}

func NewBearerToken(t string) BearerToken {
	return BearerToken{raw: t}
}

func (t BearerToken) String() string {
	return "[REDACTED]"
}

func (t BearerToken) Raw() string {
	return t.raw
}
