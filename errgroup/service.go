package main

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type fetcher interface {
	Breeds(ctx context.Context) (breeds []string, err error)
	SubBreeds(ctx context.Context, parentBreed string) (subBreeds []string, err error)
}

func process(ctx context.Context, client fetcher) ([]string, error) {
	var result []string

	breeds, err := client.Breeds(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to fetch the breeds: %w", err)
	}

	var (
		g, gctx = errgroup.WithContext(ctx)
		mutex   sync.Mutex
	)

	for _, breed := range breeds {
		breed := breed
		g.Go(func() error {
			subBreeds, err := client.SubBreeds(gctx, breed)
			if err != nil {
				return fmt.Errorf("failed while fetching the sub-breeds of '%s': %w", breed, err)
			}

			mutex.Lock()
			defer mutex.Unlock()

			if len(subBreeds) == 0 {
				result = append(result, breed)
			} else {
				for _, subBreed := range subBreeds {
					result = append(result, fmt.Sprintf("%s/%s", breed, subBreed))
				}
			}

			return nil
		})
	}

	return result, g.Wait()
}
