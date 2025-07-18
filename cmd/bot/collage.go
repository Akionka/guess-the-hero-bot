package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/anthonynsimon/bild/transform"
	"github.com/patrickmn/go-cache"
)

type Collager interface {
	Collage(options []data.Option, items []data.Item, choice *data.Option) (image.Image, error)
}

type DefaultCollager struct {
	cache *cache.Cache
}

func NewDefaultCollager(cache *cache.Cache) *DefaultCollager {
	return &DefaultCollager{
		cache: cache,
	}
}

func (c *DefaultCollager) Collage(options []data.Option, items []data.Item, choice *data.Option) (image.Image, error) {
	const (
		canvasWidth     int     = 900
		canvasHeight    int     = 900
		heroRoundRadius float64 = 25
		itemRoundRadius float64 = 10

		heroX0, heroY0        int = 110, 175
		heroGapX, heroGapY    int = 40, 40
		heroWidth, heroHeight int = 256 * 1.25, 144 * 1.25
		heroOutline           int = 5

		itemX0, itemY0        int = 276, 645
		itemGapX, itemGapY    int = 42, 30
		itemWidth, itemHeight int = 88, 64
	)

	correctOptionOutlineColor := color.RGBA{0x00, 0xFF, 0x00, 0xFF}     // Green
	wrongOptionOutlineColor := color.RGBA{0xFF, 0x00, 0x00, 0xFF}       // Red
	correctUserAnswerOutlineColor := color.RGBA{0xFF, 0xCC, 0x00, 0xFF} // Gold
	backgroundColor := color.RGBA{0x3a, 0x3a, 0x3a, 0xFF}

	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	var correctOptionOutline draw.Image = image.NewRGBA(image.Rect(0, 0, heroWidth+heroOutline*2, heroHeight+heroOutline*2))
	draw.Draw(correctOptionOutline, correctOptionOutline.Bounds(), image.NewUniform(correctOptionOutlineColor), image.Point{}, draw.Src)
	correctOptionOutline = roundedCorners(correctOptionOutline, heroRoundRadius)

	var wrongOptionOutline draw.Image = image.NewRGBA(image.Rect(0, 0, heroWidth+heroOutline*2, heroHeight+heroOutline*2))
	draw.Draw(wrongOptionOutline, wrongOptionOutline.Bounds(), image.NewUniform(wrongOptionOutlineColor), image.Point{}, draw.Src)
	wrongOptionOutline = roundedCorners(wrongOptionOutline, heroRoundRadius)

	var correctUserAnswerOutline draw.Image = image.NewRGBA(image.Rect(0, 0, heroWidth+heroOutline*2, heroHeight+heroOutline*2))
	draw.Draw(correctUserAnswerOutline, correctUserAnswerOutline.Bounds(), image.NewUniform(correctUserAnswerOutlineColor), image.Point{}, draw.Src)
	correctUserAnswerOutline = roundedCorners(correctUserAnswerOutline, heroRoundRadius)

	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(backgroundColor), image.Point{}, draw.Src)

	for i, option := range options {
		var roundedHero image.Image
		heroKey := fmt.Sprintf("img/hero_rounded_%s", option.Hero.ShortName)

		v, found := c.cache.Get(heroKey)
		if !found {
			roundedHero = roundedCorners(transform.Resize(option.Hero.Image, heroWidth, heroHeight, transform.Lanczos), heroRoundRadius)
			c.cache.Set(heroKey, roundedHero, -1)
		} else {
			roundedHero = v.(image.Image)
		}

		col, row := i%2, i/2
		x := heroX0 + col*(heroWidth+heroGapX)
		y := heroY0 + row*(heroHeight+heroGapY)

		if choice != nil && choice.Hero.ID == option.Hero.ID {
			if option.IsCorrect {
				draw.Draw(canvas, image.Rect(x-heroOutline, y-heroOutline, x+heroWidth+heroOutline, y+heroHeight+heroOutline), correctUserAnswerOutline, image.Point{}, draw.Over)
			} else {
				draw.Draw(canvas, image.Rect(x-heroOutline, y-heroOutline, x+heroWidth+heroOutline, y+heroHeight+heroOutline), wrongOptionOutline, image.Point{}, draw.Over)
			}
		} else if choice != nil {
			if option.IsCorrect {
				draw.Draw(canvas, image.Rect(x-heroOutline, y-heroOutline, x+heroWidth+heroOutline, y+heroHeight+heroOutline), correctOptionOutline, image.Point{}, draw.Over)
			}
		}

		draw.Draw(canvas, image.Rect(x, y, x+heroWidth, y+heroHeight), roundedHero, image.Point{}, draw.Over)
	}

	for i, item := range items {
		var roundedItem image.Image
		itemKey := fmt.Sprintf("img/item_rounded_%s", item.ShortName)

		v, found := c.cache.Get(itemKey)
		if !found {
			roundedItem = roundedCorners(transform.Resize(item.Image, itemWidth, itemHeight, transform.Lanczos), itemRoundRadius)
			c.cache.Set(itemKey, roundedItem, -1)
		} else {
			roundedItem = v.(image.Image)
		}

		col, row := i%3, i/3
		x := itemX0 + col*(itemWidth+itemGapX)
		y := itemY0 + row*(itemHeight+itemGapY)
		draw.Draw(canvas, image.Rect(x, y, x+itemWidth, y+itemHeight), roundedItem, image.Point{}, draw.Over)
	}

	return canvas, nil
}

func roundedCorners(img image.Image, radius float64) draw.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dst := image.NewRGBA(bounds)

	mask := image.NewAlpha(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if inRoundedCorner(x, y, width, height, radius) {
				mask.SetAlpha(x, y, color.Alpha{A: 255})
			} else {
				mask.SetAlpha(x, y, color.Alpha{A: 0})
			}
		}
	}

	draw.DrawMask(dst, bounds, img, image.Point{}, mask, image.Point{}, draw.Over)

	return dst
}

func inRoundedCorner(x, y, width, height int, radius float64) bool {
	corners := [4]struct{ cx, cy float64 }{
		{radius, radius},
		{float64(width - int(radius)), radius},
		{radius, float64(height - int(radius))},
		{float64(width - int(radius)), float64(height - int(radius))},
	}

	for _, corner := range corners {
		dx := float64(x) - corner.cx
		dy := float64(y) - corner.cy

		if dx*dx+dy*dy <= radius*radius {
			return true
		}
	}

	return x > int(radius) && x < width-int(radius) || y > int(radius) && y < height-int(radius)
}
