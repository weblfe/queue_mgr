package utils

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net/http"
)

type (
	httpResponseHandler struct {
		ctx *fiber.Ctx
	}
)

func CreateHttpNative(ctx *fiber.Ctx) (http.ResponseWriter, *http.Request, error) {
	var (
		responseWriter = NewHttpHandler(ctx)
		request, err   = createHttpRequest(ctx)
	)
	if err != nil {
		return nil, nil, err
	}
	return responseWriter, request, nil
}

func NewHttpHandler(ctx *fiber.Ctx) *httpResponseHandler {
	var h = new(httpResponseHandler)
	h.ctx = ctx
	return h
}

func createHttpRequest(ctx *fiber.Ctx) (*http.Request, error) {
	var (
		request  = ctx.Request()
		fullURL  = request.URI().String()
		body     = bytes.NewReader(request.Body())
		req, err = http.NewRequest(ctx.Method(), fullURL, body)
	)
	if err == nil {
		return req, nil
	}
	return nil, err
}

func (h *httpResponseHandler) Header() http.Header {
	var values = make(http.Header)
	if res := h.response(); res != nil {
		res.Header.VisitAll(func(key, value []byte) {
			name := string(key)
			if _, ok := values[name]; !ok {
				values[name] = []string{string(value)}
			} else {
				values[name] = append(values[name], string(value))
			}
		})
	}
	return values
}

func (h *httpResponseHandler) Write(bytes []byte) (int, error) {
	return h.ctx.Write(bytes)
}

func (h *httpResponseHandler) WriteHeader(statusCode int) {
	h.response().SetStatusCode(statusCode)
}

func (h *httpResponseHandler) response() *fasthttp.Response {
	return h.ctx.Response()
}
