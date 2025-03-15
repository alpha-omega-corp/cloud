package httputils

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type ExitOnPanicHandler struct {
	Next http.Handler
}

func (h ExitOnPanicHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<20)
			n := runtime.Stack(buf, true)
			_, err := fmt.Fprintf(os.Stderr, "panic: %v\n\n%s", err, buf[:n])
			if err != nil {
				return
			}
			os.Exit(1)
		}
	}()
	h.Next.ServeHTTP(w, req)
}
