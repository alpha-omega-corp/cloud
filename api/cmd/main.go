package main

import (
	"fmt"
	"github.com/alpha-omega-corp/cloud/api/pkg/user"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/alpha-omega-corp/cloud/core/httputils"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := bunrouter.New(
		bunrouter.WithMiddleware(reqlog.NewMiddleware(
			reqlog.WithEnabled(true),
			reqlog.WithVerbose(true),
		)))

	envUser, err := config.NewHandler().Environment("user")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		svcUser := user.NewClient(envUser.Host)
		user.RegisterClient(svcUser, router)
	}

	listenAndServe(router)
}

func listenAndServe(r *bunrouter.Router) {
	var handler http.Handler
	handler = httputils.ExitOnPanicHandler{Next: r}

	srv := &http.Server{
		Addr:         "0.0.0.0:3000",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !isServerClosed(err) {
			log.Printf("ListenAndServe failed: %s", err)
		}
	}()

	fmt.Printf("listening on http://%s\n", srv.Addr)
	awaitSignal()
}

func isServerClosed(err error) bool {
	return err.Error() == "http: Server closed"
}

func awaitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}
