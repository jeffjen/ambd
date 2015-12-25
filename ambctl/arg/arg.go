package arg

type Info struct {
	// describe what TLS role this instnace is running
	// Either *leave empty*, which is not engage in TLS
	// or *client*, which is to connect with provided certificate
	// or *server*, which is to listen with provided certificate
	ServerRole string `json:"tlsrole"`

	// certificate root path
	CA string `json:"tlscacert"`

	// certificate public private key pair
	Cert string `json:"tlscert"`
	Key  string `json:"tlskey"`

	Name string `json:"name"`

	Net       string   `json:"net"`
	From      string   `json:"src"`
	FromRange []string `json:"range"`

	// static assignment
	To []string `json:"dst,omitempty"`

	// read from discovery
	Service string `json:"srv,omitempty"`
}
