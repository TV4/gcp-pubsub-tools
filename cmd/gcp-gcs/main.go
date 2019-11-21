package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	gcpstorage "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	gcpoption "google.golang.org/api/option"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage:
  gcp-gcs [-credentialsfile=<...>|-credentialsjson=<...>] -bucket=<...> <command>`)
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
  Commands:
    ls       [<prefix>]              Lists bucket objects, optionally filtered by the given prefix
    download <object> [<object>...]  Downloads given object(s) from the bucket
    upload   <file> [<file>...]      Uploads given file(s) to the bucket
    read     <object>                Reads a bucket object to stdout
    write    <object>                Writes a bucket object from stdin
    rm       <object> [<object>...]  Deletes the given object(s) from the bucket`)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, `
GCP credentials file:
  https://developers.google.com/identity/protocols/application-default-credentials`)
		fmt.Fprintln(os.Stderr)
	}

	credentialsFile := flag.String("credentialsfile", "", "path to a GCP credentials file")
	credentialsJSON := flag.String("credentialsjson", "", "json string of a GCP credentials file content")
	bucket := flag.String("bucket", "", "bucket name")

	flag.Parse()

	var missing []string

	if *bucket == "" {
		missing = append(missing, "bucket")
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "missing flag: %s\n\n", strings.Join(missing, ", "))
		flag.Usage()
		os.Exit(1)
	}

	if *credentialsFile != "" && *credentialsJSON != "" {
		fmt.Fprint(os.Stderr, "conflict: use either credentialsfile or credentialsjson")
		os.Exit(1)
	}

	var opts []gcpoption.ClientOption
	if *credentialsFile != "" {
		opts = append(opts, gcpoption.WithCredentialsFile(*credentialsFile))
	}
	if *credentialsJSON != "" {
		opts = append(opts, gcpoption.WithCredentialsJSON([]byte(*credentialsJSON)))
	}

	client, err := gcpstorage.NewClient(context.Background(), opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bucketHandle := client.Bucket(*bucket)

	switch cmd := flag.Arg(0); cmd {
	case "ls":
		cmdLs(flag.Arg(1), bucketHandle)
	case "download":
		args := flag.Args()[1:]
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "missing object name\n\n")
			flag.Usage()
			os.Exit(1)
		}
		cmdDownload(args, bucketHandle)
	case "upload":
		args := flag.Args()[1:]
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "missing file name\n\n")
			flag.Usage()
			os.Exit(1)
		}
		cmdUpload(args, bucketHandle)
	case "read":
		objName := flag.Arg(1)
		if objName == "" {
			fmt.Fprint(os.Stderr, "missing object name\n\n")
			flag.Usage()
			os.Exit(1)
		}
		cmdRead(objName, bucketHandle)
	case "write":
		objName := flag.Arg(1)
		if objName == "" {
			fmt.Fprint(os.Stderr, "missing object name\n\n")
			flag.Usage()
			os.Exit(1)
		}
		cmdWrite(objName, bucketHandle)
	case "rm":
		args := flag.Args()[1:]
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "missing object name\n\n")
			flag.Usage()
			os.Exit(1)
		}
		cmdRm(args, bucketHandle)
	case "":
		fmt.Fprint(os.Stderr, "missing command\n\n")
		flag.Usage()
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func cmdLs(prefix string, bucket *gcpstorage.BucketHandle) {
	query := &gcpstorage.Query{Prefix: prefix}

	it := bucket.Objects(context.Background(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listing objects: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(attrs.Name)
	}
}

func cmdDownload(objects []string, bucket *gcpstorage.BucketHandle) {
	var count int
	printDownloaded := func() {
		fmt.Fprintf(os.Stderr, "Downloaded %d object(s)\n", count)
	}
	defer printDownloaded()

	for _, objName := range objects {
		if err := func() error {
			obj := bucket.Object(objName)
			r, err := obj.NewReader(context.Background())
			if err != nil {
				return fmt.Errorf("error opening object: %v", err)
			}
			defer r.Close()

			f, err := os.Create(objName)
			if err != nil {
				return fmt.Errorf("error opening file for writing: %v", err)
			}
			defer f.Close()

			if _, err := io.Copy(f, r); err != nil {
				return fmt.Errorf("error downloading object: %v", err)
			}

			count++

			return nil
		}(); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] %v\n", objName, err)
		}
	}
}

func cmdUpload(files []string, bucket *gcpstorage.BucketHandle) {
	var count int
	printUploaded := func() {
		fmt.Fprintf(os.Stderr, "Uploaded %d file(s)\n", count)
	}
	defer printUploaded()

	for _, fileName := range files {
		if err := func() error {
			f, err := os.Open(fileName)
			if err != nil {
				return fmt.Errorf("error opening file: %v", err)
			}
			defer f.Close()

			obj := bucket.Object(fileName)
			w := obj.NewWriter(context.Background())
			defer w.Close()

			if _, err := io.Copy(w, f); err != nil {
				return fmt.Errorf("error uploading file: %v", err)
			}

			count++

			return nil
		}(); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] %v\n", fileName, err)
			printUploaded()
			os.Exit(1)
		}
	}
}

func cmdRead(objName string, bucket *gcpstorage.BucketHandle) {
	obj := bucket.Object(objName)
	r, err := obj.NewReader(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s] error opening object: %v\n", objName, err)
		os.Exit(1)
	}
	defer r.Close()

	if _, err := io.Copy(os.Stdout, r); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] error reading object: %v\n", objName, err)
		os.Exit(1)
	}
}

func cmdWrite(objName string, bucket *gcpstorage.BucketHandle) {
	obj := bucket.Object(objName)
	w := obj.NewWriter(context.Background())
	defer w.Close()

	if _, err := io.Copy(w, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] error writing object: %v\n", objName, err)
		os.Exit(1)
	}
}

func cmdRm(objects []string, bucket *gcpstorage.BucketHandle) {
	var count int
	printDeleted := func() {
		fmt.Fprintf(os.Stderr, "Deleted %d object(s)\n", count)
	}
	defer printDeleted()

	for _, objName := range objects {
		if err := func() error {
			obj := bucket.Object(objName)

			if err := obj.Delete(context.Background()); err != nil {
				return fmt.Errorf("error deleting object: %v", err)
			}

			count++

			return nil
		}(); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] %v\n", objName, err)
		}
	}
}
