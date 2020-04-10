package main

import (
	"context"
	"fmt"
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

	for _, breed := range breeds {
		subBreeds, err := client.SubBreeds(ctx, breed)
		if err != nil {
			return result, fmt.Errorf("failed while fetching the sub-breeds of '%s': %w", breed, err)
		}

		if len(subBreeds) == 0 {
			result = append(result, breed)
		} else {
			for _, subBreed := range subBreeds {
				result = append(result, fmt.Sprintf("%s/%s", breed, subBreed))
			}
		}
	}

	return result, nil
}
