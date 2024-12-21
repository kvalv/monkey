package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kvalv/monkey/bin/lsp/msg"
	"github.com/kvalv/monkey/bin/lsp/scanner"
)

func getLogger() *slog.Logger {
	fh, err := os.OpenFile("/tmp/monkey.lsplog", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	return slog.New(slog.NewTextHandler(fh, &slog.HandlerOptions{AddSource: true}))
}

func handleConnection(log *slog.Logger, conn net.Conn, conf *Config) {
	sc := scanner.New(conn)
	for sc.Scan() {
		if err := sc.Err(); err != nil {
			log.Info("error scanning", "err", err)
			log.Error("error scanning", "err", err)
			// os.Exit(1)
			continue
		}
		log.Info("got ok")

		got := sc.Next()
		log.Info("received rpc message", "method", got.MethodName())
		body := got.Body
		switch message := body.(type) {
		case *msg.RequestInitialize:
			res := msg.New(msg.NewInitializeResult(message.Id))
			// b := res.Bytes()
			// b[12] = 123 // break it
			log.Info("writing result", "res", res.String())
			if _, err := conn.Write([]byte(res.String())); err != nil {
				log.Error("failed to write rpc", "err", err)
			} else {
				log.Info("wrote initialize result back")
			}
		case *msg.InitializedNotification:
			log.Info("initialized")
		case *msg.DidOpenNotification:
			log.Info("file opened", "uri", message.Params.TextDocument.URI, "version", message.Params.TextDocument.Version)
		case *msg.DidChangeNotification:
			log.Info("file changed")
		case *msg.DidSaveNotification:
			log.Info("file saved")
		case *msg.RequestHover:
			res := msg.New(&msg.ResponseHover{
				Id: message.Id,
				Result: msg.Hover{
					Contents: msg.MarkupContent{
						Kind: "markdown",
						Value: `# Hello
I am a 
- great
- markdown
lsp`,
					},
				},
			})
			fmt.Fprintf(conn, res.String())
		default:
			log.Error(fmt.Sprintf("no handler for %T", body))
		}
	}
	if err := sc.Err(); err != nil {
		log.Error("scan failure", "err", err)
	}
	log.Info("connection closed")
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: false}))
	conf := NewConfig(log)
	log.Info("starting...")

	// we want to have some graceful shutdown, so we'll listen to various
	// signals and clean up after ourselves when it comes
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt, syscall.SIGTERM)

	// time.Sleep(5 * time.Second)

	// teeReader := io.TeeReader(os.Stdin, fh)
	listener, err := net.Listen("unix", "/tmp/monkey.socks")
	if err != nil {
		panic(err)
	}
	go func() {
		<-c
		log.Info("received term signal. closing...")
		listener.Close()
	}()
	defer listener.Close()

	log.Info("listening on unix socket", "addr", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			} else {
				log.Error("error accepting connection", "err", err)
			}
		} else {
			go handleConnection(log, conn, conf)
		}
	}

}
