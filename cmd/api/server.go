package main

import (
	"context"
	"dx-mock/pkg/bgp"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (a *application) serve() error {
	var myHandler http.Handler = http.HandlerFunc(a.handleRequest)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.port),
		Handler:      myHandler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		a.logger.Info("shuting down server, signal ", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Call Shutdown() on the server like before, but now we only send on the
		// shutdownError channel if it returns an error.
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		if a.createBGP {
			a.logger.Info("removing BGP peers")
			for _, bgpPeer := range a.bgpPeers {
				// Check if the BGP peer is still active
				if bgpPeer.BGPPeerState != "active" {
					continue
				}
				err = bgp.DeleteBGPPeer(a.serverBgp, bgpPeer.ASN, net.ParseIP(bgpPeer.CustomerAddress))
				if err != nil {
					a.logger.Info("error in deleting BGP peer ", err)
				}
			}
			a.logger.Info("stopping the BGP service")
			a.serverBgp.Stop()
		}

		// Log a message to say that we're waiting for any background goroutines to
		// complete their tasks.
		a.logger.Info("completing background tasks", srv.Addr)

		// Call Wait() to block until our WaitGroup counter is zero --- essentially
		// blocking until the background goroutines have finished. Then we return nil on
		// the shutdownError channel, to indicate that the shutdown completed without
		// any issues.
		a.wg.Wait()
		shutdownError <- nil

	}()

	a.logger.Info("starting server on port ", a.config.port)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	a.logger.Info("stopped server")

	return nil
}
