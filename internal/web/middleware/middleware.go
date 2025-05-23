package middleware

import "net/http"

type MiddlewareHandler func(w http.ResponseWriter, r *http.Request, next func(*http.Request))

type Middleware struct {
	handlers []MiddlewareHandler
}

func New(next ...MiddlewareHandler) *Middleware {
	return &Middleware{handlers: next}
}

func (m *Middleware) Append(other *Middleware) *Middleware {
	m.handlers = append(m.handlers, other.handlers...)
	return m
}

func (m *Middleware) HandleFunc(pattern string, handler http.HandlerFunc) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		var invokeNextHandler = true
		next := func(_r *http.Request) {
			invokeNextHandler = true
			r = _r
		}

		for _, h := range m.handlers {
			if invokeNextHandler {
				invokeNextHandler = false
				h(w, r, next)
			}
		}

		if invokeNextHandler {
			handler(w, r)
		}
	})
}
