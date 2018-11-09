package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	maxRetries                   = 10
	errMaxRetriesReachedTemplate = "exceeded retry limit: %v"
	containerName                = "testdata"
	blobName                     = "testfile.txt"
)

func main() {
	storageAccountID := getEnv("STORAGE_ACCOUNT_ID")

	resourceIDSplit := strings.Split(storageAccountID, "/")
	if len(resourceIDSplit) < 9 {
		log.Fatal("Invalid resource ID.")
	}

	azStorage, err := NewAzStorage(resourceIDSplit[8], resourceIDSplit[4], resourceIDSplit[2], "")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	var blobContents string

	err = DoRetry(func(attempt int) (retry bool, err error) {
		blobContents, err = azStorage.GetBlob(ctx, containerName, blobName)
		return true, err
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Download Blob Contents:")
	time.Sleep(time.Second * 1)
	fmt.Println(blobContents)

	blocker := make(chan struct{})
	<-blocker
}

func getEnv(envName string) string {
	val, ok := os.LookupEnv(envName)
	if !ok {
		log.Fatalf("%s must be set.", envName)
	}

	return val
}

type retryFunc func(attempt int) (retry bool, err error)

// DoRetry keeps trying the function until the second argument
// returns false, or no error is returned.
func DoRetry(fn retryFunc) error {
	var err error
	var cont bool
	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return fmt.Errorf(errMaxRetriesReachedTemplate, err)
		}
	}
	return err
}
