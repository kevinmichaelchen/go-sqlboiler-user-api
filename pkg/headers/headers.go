package headers

const (
	// UserID is the header name that stores the current user's ID.
	// Our Envoy proxy will forward all requests (except for logging in and signing up)
	// to our Envoy authz server, which will verify the token and forward headers
	// (such as this) to upstream backends.
	UserID = "x-current-user"
)
