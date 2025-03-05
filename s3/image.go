package s3

import (
	"context"
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/akionka/akionkabot/data"
	"github.com/minio/minio-go/v7"

	_ "golang.org/x/image/webp"
)

type HeroImageFetcher struct {
	client *minio.Client
}

func NewHeroImageFetcher(client *minio.Client) *HeroImageFetcher {
	return &HeroImageFetcher{
		client: client,
	}
}

func (f HeroImageFetcher) FetchImage(ctx context.Context, shortName string) (data.Image, error) {
	object, err := f.client.GetObject(ctx, "images", fmt.Sprintf("heroes/%s.png", shortName), minio.GetObjectOptions{})
	if err != nil {
		return data.Image{}, err
	}
	img, _, err := image.Decode(object)
	if err != nil {
		return data.Image{}, err
	}
	return data.Image{Image: img}, nil
}

type ItemImageFetcher struct {
	client *minio.Client
}

func NewItemImageFetcher(client *minio.Client) *ItemImageFetcher {
	return &ItemImageFetcher{
		client: client,
	}
}

func (f ItemImageFetcher) FetchImage(ctx context.Context, shortName string) (data.Image, error) {
	object, err := f.client.GetObject(ctx, "images", fmt.Sprintf("items/%s.png", shortName), minio.GetObjectOptions{})
	if err != nil {
		return data.Image{}, err
	}
	img, _, err := image.Decode(object)
	if err != nil {
		return data.Image{}, err
	}
	return data.Image{Image: img}, nil
}
