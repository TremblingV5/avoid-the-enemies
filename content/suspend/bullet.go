package suspend

import "golang.org/x/image/math/f64"

type Bullet struct {
	Pos         f64.Vec2 // 当前位置
	From        f64.Vec2 // 发射子弹的位置
	DirectIndex int      // 子弹的方向
	Time        int      // 子弹的生命周期
	TeamIdx     int      // 属于哪个阵营
}
