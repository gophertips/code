package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gopher.tips/repository"
)

type httpClient struct {
	client http.Client
}

func (hc httpClient) Do(req *http.Request) (*http.Response, error) {
	return hc.client.Do(req)
}

func main() {
	client := repository.Client{
		HTTPClient: httpClient{
			client: *http.DefaultClient,
		},
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		ctxCancel()
	}()

	t1 := time.Now()
	result, err := process(ctx, client)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to fetch the breeds: %w", err).Error())
		fmt.Println("Quantity: ", len(result))
		os.Exit(1)
	}

	fmt.Println(result)
	fmt.Println("Quantity: ", len(result))
	fmt.Println("Processing time: ", time.Now().Sub(t1).String())
}
