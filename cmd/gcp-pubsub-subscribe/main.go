package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	gcppubsub "cloud.google.com/go/pubsub"
	gcpoption "google.golang.org/api/option"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage:
  gcp-pubsub-subscribe [-credentials=<...>] -project=<...> -subscription=<...>`)
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
GCP credentials file:
  https://developers.google.com/identity/protocols/application-default-credentials`)
		fmt.Fprintln(os.Stderr)
	}

	credentialsFile := flag.String("credentials", "", "path to a GCP credentials file")
	projectID := flag.String("project", "", "Pub/Sub project ID")
	quiet := flag.Bool("quiet", false, "do not print messages (side-effect: more speed)")
	subscriptionName := flag.String("subscription", "", "Pub/Sub subscription name")

	flag.Parse()

	var missing []string

	if *projectID == "" {
		missing = append(missing, "project")
	}

	if *subscriptionName == "" {
		missing = append(missing, "subscription")
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "missing: %s\n\n", strings.Join(missing, ", "))
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	var opts []gcpoption.ClientOption
	if *credentialsFile != "" {
		opts = append(opts, gcpoption.WithCredentialsFile(*credentialsFile))
	}
	gcpClient, err := gcppubsub.NewClient(ctx, *projectID, opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	subscription := gcpClient.Subscription(*subscriptionName)

	var mu sync.Mutex
	err = subscription.Receive(ctx, func(ctx context.Context, msg *gcppubsub.Message) {
		if !*quiet {
			mu.Lock()
			os.Stdout.Write(msg.Data)
			fmt.Println()
			mu.Unlock()
		}
		msg.Ack()
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
