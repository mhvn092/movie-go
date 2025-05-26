package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

var (
	Logger       = requestLogger()
	AuthUser     = isUserAuthorized()
	AuthAdmin    = isAdminAuthorized()
	RecoverPanic = recoverPanic()
)
