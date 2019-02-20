package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	gcppubsub "cloud.google.com/go/pubsub"
	gcpoption "google.golang.org/api/option"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage:
  gcp-pubsub-publish [-credentialsfile=<...>|-credentialsjson=<...>] -project=<...> -topic=<...>`)
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
GCP credentials file:
  https://developers.google.com/identity/protocols/application-default-credentials`)
		fmt.Fprintln(os.Stderr)
	}

	credentialsFile := flag.String("credentialsfile", "", "path to a GCP credentials file")
	credentialsJSON := flag.String("credentialsjson", "", "json string of a GCP credentials file content")
	projectID := flag.String("project", "", "Pub/Sub project ID")
	topicName := flag.String("topic", "", "Pub/Sub topic")

	flag.Parse()

	var missing []string

	if *projectID == "" {
		missing = append(missing, "project")
	}

	if *topicName == "" {
		missing = append(missing, "topic")
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
	topic := gcpClient.Topic(*topicName)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := make([]byte, len(scanner.Bytes()))
		copy(data, scanner.Bytes())
		topic.Publish(ctx, &gcppubsub.Message{Data: data})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	topic.Stop()
}
