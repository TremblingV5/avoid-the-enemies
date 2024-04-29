package weapon

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
)

type Melee struct {
	Type  string        // 武器类型
	Image *ebiten.Image // 加载武器的图片
	angle float64       // 武器的旋转角度
	spin  float64       // 武器的旋转速度
	Trail []f64.Vec2    // 武器的轨迹
}

func (w *Melee) GetImage() *ebiten.Image {
	return w.Image
}

func (w *Melee) Spin() {
	w.angle += w.spin
	w.angle = math.Mod(w.angle, 2*math.Pi)
}

func (w *Melee) Copy() *Melee {
	return &Melee{
		Type:  w.Type,
		Image: w.Image,
		angle: w.angle,
		spin:  w.spin,
		Trail: w.Trail,
	}
}

// DrawTrail 在绘制时，绘制轨迹效果
func (w *Melee) DrawTrail(screen *ebiten.Image) {
	// 绘制轨迹效果
	for i := 1; i < len(w.Trail); i++ {
		prevPos := w.Trail[i-1]
		currPos := w.Trail[i]
		// 绘制当前位置与前一位置之间的轨迹线段
		// 根据需要设置线段的颜色、粗细等属性
		// 例如使用 ebitenutil.DrawLine() 函数
		ebitenutil.DrawLine(screen, prevPos[0], prevPos[1], currPos[0], currPos[1], color.RGBA{255, 255, 255, 128})
	}
}
