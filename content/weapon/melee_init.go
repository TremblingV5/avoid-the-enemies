package weapon

import (
	"avoid-the-enemies/content/resources"
	"math"
)

func NewSickle() *Melee {
	return &Melee{
		Type:  "sickle",
		Image: resources.SickleImage,
		angle: 0,
		spin:  1.75 * math.Pi / 60, // 每帧转动的角度（弧度）
	}
}

func NewSword() *Melee {
	return &Melee{
		Type:  "sword",
		Image: resources.SwordImage,
		angle: 0,
		spin:  2 * math.Pi / 60, // 每帧转动的角度（弧度）
	}
}
