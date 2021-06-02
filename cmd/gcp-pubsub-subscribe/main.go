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
  gcp-pubsub-subscribe [-credentialsfile=<...>|-credentialsjson=<...>] -project=<...> -subscription=<...>`)
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
GCP credentials file:
  https://developers.google.com/identity/protocols/application-default-credentials`)
		fmt.Fprintln(os.Stderr)
	}

	ack := flag.Bool("ack", false, "ack messages")
	credentialsFile := flag.String("credentialsfile", "", "path to a GCP credentials file")
	credentialsJSON := flag.String("credentialsjson", "", "json string of a GCP credentials file content")
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

	if *credentialsFile != "" && *credentialsJSON != "" {
		fmt.Fprint(os.Stderr, "conflict: use either credentialsfile or credentialsjson")
		os.Exit(1)
	}

	ctx := context.Background()

	var opts []gcpoption.ClientOption
	if *credentialsFile != "" {
		opts = append(opts, gcpoption.WithCredentialsFile(*credentialsFile))
	}
	if *credentialsJSON != "" {
		opts = append(opts, gcpoption.WithCredentialsJSON([]byte(*credentialsJSON)))
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
		if *ack {
			msg.Ack()
		}
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
