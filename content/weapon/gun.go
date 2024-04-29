package weapon

import (
	"avoid-the-enemies/content/suspend"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/math/f64"
)

type Gun struct {
	Type         string        // 武器类型
	Image        *ebiten.Image // 加载武器的图片
	bullet       *ebiten.Image // 子弹图片
	speed        float64       // 子弹的速度
	distance     float64       // 子弹的射程
	damage       float64       // 子弹的伤害值
	LastFireTime time.Time     // 上次开火的时间
	shotPlayer   *audio.Player // 射击音效
}

func (w *Gun) GetImage() *ebiten.Image {
	return w.Image
}

func (w *Gun) Copy() *Gun {
	return &Gun{
		Type:       w.Type,
		Image:      w.Image,
		bullet:     w.bullet,
		speed:      w.speed,
		distance:   w.distance,
		damage:     w.damage,
		shotPlayer: w.shotPlayer,
	}
}

func (w *Gun) Fire(g game, player player) {
	if err := w.shotPlayer.Rewind(); err != nil {
		return
	}
	w.shotPlayer.Play()
	// 每次开火生成一颗子弹，移动的距离为 distance 速度为 speed 图片为 bullet
	x, y := player.GenerateBullet()
	bullet := &suspend.Bullet{
		Pos:         f64.Vec2{x, y}, // 子弹的起始位置
		From:        f64.Vec2{x, y},
		DirectIndex: player.GetDirectIndex(),
	}
	g.IncrUniqueId()
	g.AddSuspendForUniqueId(g.GetUniqueId(), bullet)
}
