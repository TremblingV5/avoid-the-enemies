package resources

import (
	"avoid-the-enemies/resources/images"
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AkImage     *ebiten.Image
	BulletImage *ebiten.Image
	RunnerImage *ebiten.Image
	SickleImage *ebiten.Image
	SwordImage  *ebiten.Image
	SkillImage  *ebiten.Image
	FireImage   *ebiten.Image
)

func InitImage() {
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	RunnerImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Sickle_png))
	if err != nil {
		log.Fatal(err)
	}
	SickleImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Sword_png))
	if err != nil {
		log.Fatal(err)
	}
	SwordImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.AK_png))
	if err != nil {
		log.Fatal(err)
	}
	AkImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Bullet_png))
	if err != nil {
		log.Fatal(err)
	}
	BulletImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Skill_png))
	if err != nil {
		log.Fatal(err)
	}
	SkillImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Fire_png))
	if err != nil {
		log.Fatal(err)
	}
	FireImage = ebiten.NewImageFromImage(img)
}
