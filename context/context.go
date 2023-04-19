package context

type contextKey struct {
	name string
}

var HttpWriterKey = &contextKey{"httpWriter"}
var HttpReaderKey = &contextKey{"httpReader"}
var CookieRefreshTokenKey = &contextKey{"refreshToken"}
