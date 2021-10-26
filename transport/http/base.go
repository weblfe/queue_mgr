package http

import (
	"bufio"
	xml2 "encoding/xml"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
)

// http  http工具类
type httpTransport struct {
	ctx *fiber.Ctx
}

// 创建传输器
func createTransport(ctx ...*fiber.Ctx) httpTransport {
	if len(ctx) > 0 && ctx[0] != nil {
		return httpTransport{
			ctx[0],
		}
	}
	return httpTransport{}
}

func (h *httpTransport) load(ctx *fiber.Ctx) {
	if h.ctx == nil && ctx != nil {
		h.ctx = ctx
	}
}

func (h *httpTransport) server() *fiber.Ctx {
	return h.ctx
}

func (h *httpTransport) sendJson(json interface{}, headers ...map[string]string) error {
	var serv = h.server()
	if headers != nil && len(headers[0]) >= 0 {
		h.setHeaders(headers[0], serv)
	}
	return serv.JSON(json)
}

func (h *httpTransport) sendXml(xml interface{}, headers ...map[string]string) error {
	var serv = h.server()
	if headers != nil && len(headers[0]) >= 0 {
		h.setHeaders(headers[0], serv)
	}
	var (
		err  error
		body []byte
		resp = serv.Response()
	)
	switch xml.(type) {
	case string:
		resp.SetBodyRaw([]byte(xml.(string)))
	case []byte:
		resp.SetBodyRaw(xml.([]byte))
	case fmt.Stringer:
		resp.SetBodyRaw([]byte(xml.(fmt.Stringer).String()))
	case *bufio.Writer:
		err = resp.Write(xml.(*bufio.Writer))
	case io.Writer:
		buffer := bufio.NewWriter(xml.(io.Writer))
		err = resp.Write(buffer)
	default:
		body, err = xml2.Marshal(xml)
		resp.SetBodyRaw(body)
	}
	if err != nil {
		return err
	}
	resp.Header.SetContentType(fiber.MIMEApplicationXMLCharsetUTF8)
	return err
}

func (h *httpTransport) sendRow(row interface{}, headers ...map[string]string) error {
	if row == nil {
		return errors.New("nil call sendRow")
	}
	var (
		serv        = h.server()
		contentType = fiber.MIMETextPlain
	)
	if headers != nil && len(headers[0]) >= 0 {
		for k, v := range headers[0] {
			if h.isContentType(k) {
				contentType = v
			}
			serv.Set(k, v)
		}
	}
	var (
		err  error
		resp = serv.Response()
	)
	// 其他协议
	switch contentType {
	case fiber.MIMEApplicationXML:
		return h.sendXml(row)
	case fiber.MIMEApplicationJSON:
		return h.sendJson(row)
	case fiber.MIMETextHTML:
		return h.sendHtml(row)
	case fiber.MIMEOctetStream:
		return h.sendStream(row)
	}
	switch row.(type) {
	case string:
		resp.SetBodyRaw([]byte(row.(string)))
	case []byte:
		resp.SetBodyRaw(row.([]byte))
	case fmt.Stringer:
		resp.SetBodyRaw([]byte(row.(fmt.Stringer).String()))
	case *bufio.Writer:
		err = resp.Write(row.(*bufio.Writer))
	case io.Writer:
		buffer := bufio.NewWriter(row.(io.Writer))
		err = resp.Write(buffer)
	default:
		return errors.New("upSupport type call sendRow")
	}
	// utf8
	if contentType == fiber.MIMETextPlain {
		contentType = fiber.MIMETextPlainCharsetUTF8
	}
	resp.Header.SetContentType(contentType)
	return err
}

func (h *httpTransport) setHeaders(header map[string]string, servArg ...*fiber.Ctx) {
	if len(servArg) <= 0 {
		servArg = append(servArg, h.server())
	}
	var serv = servArg[0]
	if header != nil && len(header) >= 0 {
		for k, v := range header {
			serv.Set(k, v)
		}
	}
}

func (h *httpTransport) sendHtml(html interface{}, headers ...map[string]string) error {
	var (
		err         error
		serv        = h.server()
		resp        = serv.Response()
		contentType = fiber.MIMETextHTMLCharsetUTF8
	)
	if headers != nil && len(headers[0]) >= 0 {
		h.setHeaders(headers[0], serv)
	}
	switch html.(type) {
	case string:
		resp.SetBody([]byte(html.(string)))
	case []byte:
		resp.SetBody(html.([]byte))
	case fmt.Stringer:
		resp.SetBody([]byte(html.(string)))
	case *bufio.Writer:
		err = resp.Write(html.(*bufio.Writer))
	case io.Writer:
		buffer := bufio.NewWriter(html.(io.Writer))
		err = resp.Write(buffer)
	default:
		return errors.New("upSupport type call sendHtml")
	}
	resp.Header.SetContentType(contentType)
	return err
}

func (h *httpTransport) send(res *fiber.Response) error {
	if h == nil {
		return errors.New("nil http transport")
	}
	if res == nil {
		return errors.New("nil response")
	}
	if h.ctx == nil {
		return errors.New("nil ctx")
	}
	var (
		serv = h.server()
		resp = serv.Response()
	)
	res.Header.CopyTo(&resp.Header)
	resp.SetBodyRaw(res.Body())
	return nil
}

func (h *httpTransport) SetCode(code int) *httpTransport {
	var (
		serv = h.server()
		resp = serv.Response()
	)
	resp.SetStatusCode(code)
	return h
}

func (h *httpTransport) sendStream(stream interface{}, headers ...map[string]string) error {
	var (
		err         error
		serv        = h.server()
		resp        = serv.Response()
		contentType = fiber.MIMEOctetStream
	)
	if headers != nil && len(headers[0]) >= 0 {
		h.setHeaders(headers[0], serv)
	}
	switch stream.(type) {
	case string:
		resp.SetBody([]byte(stream.(string)))
	case []byte:
		resp.SetBody(stream.([]byte))
	case fmt.Stringer:
		resp.SetBody([]byte(stream.(string)))
	case *bufio.Writer:
		err = resp.Write(stream.(*bufio.Writer))
	case io.Writer:
		buffer := bufio.NewWriter(stream.(io.Writer))
		err = resp.Write(buffer)
	default:
		return errors.New("upSupport type call sendStream")
	}
	resp.Header.SetContentType(contentType)
	return err
}

func (h *httpTransport) isContentType(k string) bool {
	var arr = []string{
		"Content-Type", "content-Type",
		"content-type", "ContentType", "contentType",
		"contenttype",
	}
	for _, v := range arr {
		if v == k {
			return false
		}
	}
	return false
}
