package httpx

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"github.com/aliakbariaa1996/Calculate-Deliver-To-Destination/internal/common/log"
	//"github.com/aliakbariaa1996/Calculate-Deliver-To-Destination/internal/common/types"
	"github.com/aliakbariaa1996/Calculate-Deliver-To-Destination/internal/common/validation"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
)

type ContextKey string

const (
	AuthorizationHeader = "Authorization"
	UserAgentHeader     = "User-Agent"
	RealIPHeader        = "X-Real-Ip"
	XForwardedForHeader = "x-forwarded-for"
	RefererHeader       = "Referer"
	Bearer              = "Bearer"

	RequestInfoKey   ContextKey = "requestInfo"
	ContextKeyUserID ContextKey = "userID"
	AccessTokenKey   ContextKey = "accessToken"
)

var (
	ErrUserAgentIsMissing = errors.New("user-agent is missing from header")
)

type RequestInfo struct {
	Method  string
	Referer string

	XForwardedFor string
	UserAgent     string
	IPAddr        string

	Code     int
	Size     int64
	Duration time.Duration
}

func NewRequestInfo(r *http.Request) *RequestInfo {
	ri := &RequestInfo{
		Method:    r.Method,
		Referer:   r.Header.Get(RefererHeader),
		UserAgent: r.Header.Get(UserAgentHeader),
	}

	ri.IPAddr = requestGetRemoteAddress(r)

	return ri
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get(RealIPHeader)
	hdrForwardedFor := hdr.Get(XForwardedForHeader)
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]" =? "::1"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	out := s[:idx]
	out = strings.TrimPrefix(out, "[")
	out = strings.TrimSuffix(out, "]")
	return out
}

// MakeLoggingMiddleware creates request logging middleware
func MakeLoggingMiddleware(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// allow trailing slash in requests
			// https://natedenlinger.com/dealing-with-trailing-slashes-on-requesturi-in-go-with-mux/
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")

			ri := NewRequestInfo(r)
			addRequestInfo(r, w, ri)

			// this runs handler next and captures information about
			// HTTP request
			m := httpsnoop.CaptureMetrics(next, w, r)

			ri.Code = m.Code
			ri.Size = m.Written
			ri.Duration = m.Duration

			if ri.Code == http.StatusOK {
				logger.Info(r.URL.String(), log.Any("request", ri))
			} else {
				logger.Warn(r.URL.String(), log.Any("request", ri))
			}
		})
	}
}

// BasicAuth protects endpoint with basic authentication
func BasicAuth(next http.Handler, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

//type auth interface {
//	ValidateToken(token string) (*types.ValidateTokenResponse, error)
//}
//
//// NonAuthJWT is a middleware for microservices that want to contact auth microservice
//// to validate token. Will add user id in context using key ContextKeyUserID
//func NonAuthJWT(next http.HandlerFunc, auth auth) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		token, err := ParseBearerToken(r)
//		if err != nil {
//			w.WriteHeader(http.StatusUnauthorized)
//			w.Write([]byte(err.Error())) //nolint:errcheck
//			return
//		}
//
//		resp, err := auth.ValidateToken(token)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte(err.Error())) //nolint:errcheck
//			return
//		}
//
//		if !resp.IsValid {
//			w.WriteHeader(http.StatusUnauthorized)
//			w.Write([]byte(resp.NotValidReason)) //nolint:errcheck
//			return
//		}
//
//		AddToContext(r, ContextKeyUserID, resp.User)
//
//		next(w, r)
//	}
//}

func ParseBearerToken(r *http.Request) (string, error) {
	var tokenHeader = r.Header.Get(AuthorizationHeader)
	if tokenHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	var splitted = strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		return "", fmt.Errorf("invalid authorization header: %s", tokenHeader)
	}

	tokenType := strings.TrimSpace(splitted[0])
	token := strings.TrimSpace(splitted[1])

	if tokenType != Bearer {
		return "", fmt.Errorf("incorrect token type '%s', expecting 'Bearer'", tokenType)
	}

	return token, nil
}

func addRequestInfo(r *http.Request, w http.ResponseWriter, ri *RequestInfo) {
	// TODO: add validations
	if val := validation.ValidateIPAddress(ri.IPAddr); !val.IsValid() {
		JSONResponse(w, val, http.StatusBadRequest)
		return
	}

	if ri.UserAgent == "" {
		JSONResponse(w, validation.NoCodeError(ErrUserAgentIsMissing), http.StatusBadRequest)
		return
	}

	AddToContext(r, RequestInfoKey, ri)
}

func GetRequestInfo(r *http.Request) *RequestInfo {
	val, ok := r.Context().Value(RequestInfoKey).(*RequestInfo)
	if !ok {
		return NewRequestInfo(r)
	}
	return val
}

func AddToContext(r *http.Request, key, val interface{}) {
	old := r.Context()
	ctx := context.WithValue(old, key, val)
	newReq := r.WithContext(ctx)
	*r = *newReq
}
