// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func Init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	InitImage()
	InitFont()
	InitWeapon()
}

// Draw 每次绘制都会调用这个函数，重新设置画面元素的内容
func (g *Game) Draw(screen *ebiten.Image) {
	var titleTexts string
	var texts string
	switch g.mode {
	case ModeTitle:
		titleTexts = "Avoid the Enemies"
		texts = "PRESS SPACE KEY TO START"
	case ModeGameOver:
		titleTexts = "Game Over"
		texts = "PRESS SPACE KEY TO RESTART"
	}

	// 绘制标题
	op := &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, 5*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = titleFontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, titleTexts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   titleFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, 7*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, texts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)

	if g.mode == ModeGame {
		// 绘制分数
		op = &text.DrawOptions{}
		op.GeoM.Translate(3, 3)
		op.ColorScale.ScaleWithColor(color.White)
		op.LineSpacing = fontSize
		text.Draw(screen, "Score: "+strconv.Itoa(g.player.score), &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   fontSize,
		}, op)

		// 绘制游戏时间
		op = &text.DrawOptions{}
		op.GeoM.Translate(screenWidth/2, 3)
		op.ColorScale.ScaleWithColor(color.White)
		op.LineSpacing = fontSize
		text.Draw(screen, "SurvivalTime: "+strconv.Itoa(int(time.Since(g.player.startTime).Seconds()))+"s", &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   fontSize,
		}, op)

		// 绘制技能效果
		if g.player.isSkill {
			g.player.skillFrame++
			op := &ebiten.DrawImageOptions{}
			// 位于血条上方，血条高度为 5
			op.GeoM.Translate(g.player.x-16, g.player.y-5-16)
			//op.GeoM.Translate(g.player.x-8, g.player.y-5-36)
			i := (g.player.skillFrame / 5) % 4
			//i := (g.player.skillFrame / 5) % 90
			sx, sy := i*64, 0
			//sx, sy := i*48, 0
			screen.DrawImage(fireImage.SubImage(image.Rect(sx, sy, sx+64, sy+64)).(*ebiten.Image), op)
			//screen.DrawImage(skillImage.SubImage(image.Rect(sx, sy, sx+48, sy+36)).(*ebiten.Image), op)
		}

		// 绘制角色
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.player.x, g.player.y)
		i := (g.player.count / 5) % frameCount
		sx, sy := frameOX+i*frameWidth, frameOY
		screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

		// 绘制角色武器
		if g.player.weapon != nil {
			switch g.player.weapon.(type) {
			case *MeleeWeapon:
				weapon := g.player.weapon.(*MeleeWeapon)
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Rotate(weapon.angle)
				op.GeoM.Translate(g.player.x+g.player.weaponX, g.player.y+g.player.weaponY)
				screen.DrawImage(weapon.Image.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image), op)
				weapon.DrawTrail(screen)
			case *RangedWeapon:
				weapon := g.player.weapon.(*RangedWeapon)
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Rotate(directions[g.player.directIdx].spin)
				op.GeoM.Translate(rotateAdjust[g.player.directIdx].dx*frameWidth, rotateAdjust[g.player.directIdx].dy*frameHeight)
				op.GeoM.Translate(g.player.x, g.player.y)
				screen.DrawImage(weapon.Image.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image), op)
			}
		}

		// 绘制武器发射产物
		for _, suspend := range g.suspends {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Rotate(directions[g.player.directIdx].spin)
			op.GeoM.Translate(suspend.pos[0], suspend.pos[1])
			screen.DrawImage(suspend.rangeWeapon.bullet.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image), op)
		}

		// 绘制怪物
		for _, monster := range g.monsters {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(monster.x, monster.y)
			screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
			// 绘制怪物武器
			if monster.weapon != nil {
				switch monster.weapon.(type) {
				case *MeleeWeapon:
					weapon := monster.weapon.(*MeleeWeapon)
					op = &ebiten.DrawImageOptions{}
					op.GeoM.Rotate(weapon.angle)
					op.GeoM.Translate(monster.x+monster.weaponX, monster.y+monster.weaponY)
					screen.DrawImage(weapon.Image.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image), op)
					weapon.DrawTrail(screen)
				case *RangedWeapon:
					weapon := monster.weapon.(*RangedWeapon)
					op = &ebiten.DrawImageOptions{}
					op.GeoM.Rotate(directions[monster.directIdx].spin)
					op.GeoM.Translate(rotateAdjust[monster.directIdx].dx*frameWidth, rotateAdjust[monster.directIdx].dy*frameHeight)
					op.GeoM.Translate(monster.x, monster.y)
					screen.DrawImage(weapon.Image.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image), op)
				}
			}
		}

		// 设置血条的位置和尺寸
		x := g.player.x
		y := g.player.y - 5                                  // 位于角色头顶上方
		width := float64(frameWidth) * g.player.health / 100 // 血条宽度根据当前血量动态变化
		height := 5                                          // 血条高度
		// 绘制血条底部
		ebitenutil.DrawRect(screen, x, y, float64(frameWidth), float64(height), color.Gray{0x80})
		// 绘制血条
		ebitenutil.DrawRect(screen, x, y, float64(width), float64(height), color.RGBA{0xFF, 0x00, 0x00, 0xFF})

		// 地图上的武器
		for id, weapon := range g.weapons {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(g.weaponPosition[id][0], g.weaponPosition[id][1])
			screen.DrawImage(weapon.GetImage().SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image), op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	Init()
	ebiten.SetWindowSize(screenWidth*3, screenHeight*3)
	ebiten.SetWindowTitle("Avoid the Enemies")
	g := &Game{}
	g.init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
