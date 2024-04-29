package character

import (
	"avoid-the-enemies/content/base"
	"avoid-the-enemies/content/weapon"
	"time"
)

type Character struct {
	Base              base.Element
	Weapon            weapon.WeaponComponents // 角色所拥有的武器
	WeaponX, WeaponY  float64                 // 武器相对于角色中心的偏移量
	Health            float64                 // 角色的生命值
	LastCollisionTime time.Time               // 上次碰撞发生的时间
	DirectIdx         int                     // 角色的方向
}
