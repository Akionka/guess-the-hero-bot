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
	"github.com/patrickmn/go-cache"

	_ "golang.org/x/image/webp"
)

type HeroImageFetcher struct {
	client *minio.Client
	cache  *cache.Cache
}

func NewHeroImageFetcher(client *minio.Client, cache *cache.Cache) *HeroImageFetcher {
	return &HeroImageFetcher{
		client: client,
		cache:  cache,
	}
}

func (f HeroImageFetcher) FetchImage(ctx context.Context, shortName string) (data.Image, error) {
	imgKey := fmt.Sprintf("hero_img_%s", shortName)
	v, found := f.cache.Get(imgKey)
	if found {
		return data.Image{Image: v.(image.Image)}, nil
	}

	object, err := f.client.GetObject(ctx, "images", fmt.Sprintf("heroes/%s.png", shortName), minio.GetObjectOptions{})
	if err != nil {
		return data.Image{}, err
	}
	img, _, err := image.Decode(object)
	if err != nil {
		return data.Image{}, err
	}
	f.cache.Set(imgKey, img, cache.NoExpiration)

	return data.Image{Image: img}, nil
}

type ItemImageFetcher struct {
	client *minio.Client
	cache  *cache.Cache
}

func NewItemImageFetcher(client *minio.Client, cache *cache.Cache) *ItemImageFetcher {
	return &ItemImageFetcher{
		client: client,
		cache:  cache,
	}
}

func (f ItemImageFetcher) FetchImage(ctx context.Context, shortName string) (data.Image, error) {
	imgKey := fmt.Sprintf("item_img_%s", shortName)
	v, found := f.cache.Get(imgKey)
	if found {
		return data.Image{Image: v.(image.Image)}, nil
	}

	object, err := f.client.GetObject(ctx, "images", fmt.Sprintf("items/%s.png", shortName), minio.GetObjectOptions{})
	if err != nil {
		return data.Image{}, err
	}
	img, _, err := image.Decode(object)
	if err != nil {
		return data.Image{}, err
	}
	f.cache.Set(imgKey, img, cache.NoExpiration)

	return data.Image{Image: img}, nil
}
