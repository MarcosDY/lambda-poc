package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/marcosdy/lambda-poc/extension"
	"github.com/marcosdy/lambda-poc/secret"
	"github.com/marcosdy/lambda-poc/util"
)

const (
	outputPath = "/tmp"
)

var (
	extensionClient = extension.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	secretName      = os.Getenv("SECRET_NAME")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sm, err := secret.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create secrets client: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		log.Printf("Exiting: %v", s)
	}()

	if _, err := extensionClient.Register(ctx, filepath.Base(os.Args[0])); err != nil {
		log.Fatal(err)
	}

	if err := initialiseExtension(ctx, sm); err != nil {
		log.Fatalf("Failed to initialize extension: %v", err)
	}

	processEvents(ctx)
}

func initialiseExtension(ctx context.Context, sm secret.SecretManager) error {
	startAt := time.Now()
	defer func() {
		elapse := time.Since(startAt)
		log.Printf("initialize takes: %s", strconv.FormatInt(elapse.Milliseconds(), 10))
	}()

	// Get a Secret binary from Secret Manager
	secretBinary, err := sm.GetSecret(ctx, secretName)
	if err != nil {
		return err
	}

	if err := util.SaveSvid(secretBinary, outputPath); err != nil {
		return err
	}

	return nil
}

func processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			res, err := extensionClient.NextEvent(ctx)
			if err != nil {
				log.Fatalf("Error receiving event: %s", err)
			}
			// Exit if we receive a SHUTDOWN event
			if res.EventType == extension.Shutdown {
				log.Println("Received SHUTDOWN event, exiting")
				return
			}
		}
	}
}
