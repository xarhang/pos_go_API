package middleware

import (
	"bytes"
	"io"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// now := time.Now()
		wb := &toolBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
		}
		c.Writer = wb
		requestDump, _ := httputil.DumpRequest(c.Request, true)
		requestID := c.Request.URL.Path + ".." + uuid.New().String()
		os.WriteFile("./uploads/request/"+strings.Replace(requestID, "/", "=", -1), requestDump, 0755)
		// urls := string(c.Request.Host + c.Request.URL.Path)
		ioRequest := c.Request.Body
		requestBody := new(bytes.Buffer)
		requestBody.ReadFrom(ioRequest)
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(requestBody.Bytes())))

		c.Next()

		originBytes := wb.body

		// tbLogger := dto.TbLogger{
		// 	Url:         urls,
		// 	Method:      c.Request.Method,
		// 	Request:     string(requestDump),
		// 	RequestBody: string(requestBody.String()),
		// 	Code:        c.Writer.Status(),
		// 	Response:    originBytes.String(),
		// 	Accesstime:  now.Format(time.DateTime),
		// 	Handletime:  time.Since(now).String(),
		// }

		// db.Create(&tbLogger)
		os.WriteFile("./uploads/response/"+strings.Replace(requestID, "/", "=", -1), originBytes.Bytes(), 0755)
		wb.body = &bytes.Buffer{}
		json := `{"response_data":` + originBytes.String() +
			`,"transaction_id":"` + strings.Replace(requestID, "/", "=", -1) + `"}`
		wb.Write([]byte(json))
		wb.ResponseWriter.Write(wb.body.Bytes())

	}
}

type toolBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r toolBodyWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}
