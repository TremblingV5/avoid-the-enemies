package weapon

import (
	"avoid-the-enemies/content/suspend"

	"github.com/hajimehoshi/ebiten/v2"
)

type WeaponType int

const (
	SpinWeaponType WeaponType = iota
	RangeWeaponType
)

type WeaponComponents interface {
	GetImage() *ebiten.Image
	Copy() WeaponComponents
	PlayAudio()   // 播放武器音效
	AutoResolve() // 发挥武器功能
	GetWeaponType() WeaponType
}

type SpinWeapon interface {
	WeaponComponents
	Spin()
	DrawTrail(screen *ebiten.Image)
}

type RangeWeapon interface {
	WeaponComponents
	Fire(game game, player player)
}

type game interface {
	IncrUniqueId()
	GetUniqueId() int
	AddSuspendForUniqueId(int, *suspend.Bullet)
}

type player interface {
	GenerateBullet() (float64, float64)
	GetId() int
	GetDirectIndex() int
}
