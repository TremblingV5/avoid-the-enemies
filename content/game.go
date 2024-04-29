package main

import (
	"avoid-the-enemies/content/character"
	"avoid-the-enemies/content/weapon"
	"bytes"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/math/f64"
)

type Game struct {
	mode           Mode
	player         *Player
	player1        *character.Player
	uniqueId       int
	monsters       map[int]*Player
	monsterTarget  map[int]f64.Vec2 // 记录每个怪物的目标位置
	monsterTimer   map[int]int      // 记录每个怪物的计时器
	weaponTimer    time.Time        // 武器刷新时间
	weaponPosition map[int]f64.Vec2 // 武器位置
	weapons        map[int]Weapon
	suspends       map[int]*Suspend
	hitPlayer      *audio.Player
}

func (g *Game) init() {
	g.mode = ModeTitle
	g.player = &Player{
		x:                 screenWidth/2 - frameWidth/2,
		y:                 screenHeight/2 - frameHeight/2,
		speed:             2.0, // 您可以根据需要调整这个值
		weaponX:           frameWidth / 2,
		weaponY:           frameHeight / 2,
		health:            100,
		lastCollisionTime: time.Now(),
		directIdx:         0,
		id:                1,
		score:             0,
		isSkill:           false,
		skillFrame:        0,
		startTime:         time.Now(),
	}
	g.monsters = make(map[int]*Player)
	g.monsterTarget = make(map[int]f64.Vec2)
	g.monsterTimer = make(map[int]int)
	g.weaponTimer = time.Now()
	g.weaponPosition = make(map[int]f64.Vec2)
	g.weapons = make(map[int]Weapon)
	g.suspends = make(map[int]*Suspend)
	g.uniqueId = 1

	if audioContext == nil {
		audioContext = audio.NewContext(48000)
	}
	jabD, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		log.Fatal(err)
	}
	g.hitPlayer, err = audioContext.NewPlayer(jabD)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.mode = ModeGame
			g.player.startTime = time.Now()
		}
	case ModeGame:
		if err := g.gameRender(); err != nil {
			return err
		}
	case ModeGameOver:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.init()
			g.mode = ModeTitle
		}
	}

	return nil
}

func (g *Game) gameRender() error {
	g.player.count++

	// 响应按键事件
	g.resolveKeyPress()

	// 如果 GIF 正在播放且播放完成，则停止播放
	if g.player.isSkill && time.Since(g.player.skillTime) > time.Second*3 {
		g.player.isSkill = false
	}

	// 更新所有远程武器的发射产物位置
	SuspendMove(g)

	// 生成怪物
	GenerateMonster(g)

	// 武器在地图上随机位置刷新
	GenerateWeapon(g)

	g.resolvePickupWeapons()

	if err := g.resolvePlayWeapon(); err != nil {
		return err
	}

	if err := g.resolveNpc(); err != nil {
		return err
	}

	return nil
}

func (g *Game) resolveKeyPress() {
	// 检查键盘输入，人物移动
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.Move(-g.player.speed, 0)
		g.player.directIdx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.Move(g.player.speed, 0)
		g.player.directIdx = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.Move(0, -g.player.speed)
		g.player.directIdx = 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.Move(0, g.player.speed)
		g.player.directIdx = 1
	}
	// 按下 q 键可以释放技能 && 距离上一次释放技能时间大于技能冷却时间
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && time.Since(g.player.skillTime) > time.Second*5 {
		if g.player.score >= 20 {
			g.player.isSkill = true
			g.player.skillTime = time.Now()
			g.player.score -= 20
			g.player.skillFrame = 0
		}
	}

	// 如果人物有武器，且是远程武器，按空格键开火
	if g.player.weapon != nil && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// 如果上次未开火过，执行开火操作
		switch g.player.weapon.(type) {
		case *RangedWeapon:
			weapon := g.player.weapon.(*RangedWeapon)
			weapon.Fire(g, g.player)
		}
	}
}

func (g *Game) resolvePickupWeapons() {
	for id, weapon := range g.weapons {
		// 玩家移动到武器位置可以获得武器
		if IsTouch(g.player.x, g.player.y, g.weaponPosition[id][0], g.weaponPosition[id][1]) {
			g.player.weapon = weapon
			delete(g.weapons, id)
			delete(g.weaponPosition, id)
			break
		}
		// 怪物移动到武器位置可以获得武器
		for _, monster := range g.monsters {
			if IsTouch(monster.x, monster.y, g.weaponPosition[id][0], g.weaponPosition[id][1]) {
				monster.weapon = weapon
				delete(g.weapons, id)
				delete(g.weaponPosition, id)
				break
			}
		}
	}
}

func (g *Game) resolveNpc() error {
	// 怪物移动
	for id, monster := range g.monsters {
		// 根据玩家位置，确定怪物 directIdx
		if g.player.x < monster.x {
			if g.player.y < monster.y {
				monster.directIdx = 2
			} else {
				monster.directIdx = 1
			}
		} else {
			if g.player.y < monster.y {
				monster.directIdx = 3
			} else {
				monster.directIdx = 0
			}
		}

		target := g.monsterTarget[id]
		timer := g.monsterTimer[id]

		// 更新计时器
		timer++
		g.monsterTimer[id] = timer

		// 每隔一定时间更新一次目标位置
		if timer >= 60 {
			// 以玩家为目标
			g.monsterTarget[id] = f64.Vec2{g.player.x, g.player.y}
			g.monsterTimer[id] = 0
		}

		// 计算当前位置到目标位置的方向向量
		directionX := target[0] - monster.x
		directionY := target[1] - monster.y

		// 根据怪物的速度进行插值
		monster.Move(directionX*monster.speed, directionY*monster.speed)

		if monster.weapon != nil {
			switch monster.weapon.(type) {
			// 怪物武器旋转
			case *MeleeWeapon:
				weapon := monster.weapon.(*MeleeWeapon)
				weapon.Spin()
				// 武器碰撞到玩家，降低生命值
				weaponCenterOffsetX := monster.weaponX // 武器中心相对于怪物中心的 X 坐标偏移
				weaponCenterOffsetY := monster.weaponY // 武器中心相对于怪物中心的 Y 坐标偏移

				// 考虑武器的旋转角度，将偏移向量旋转到合适的位置
				weaponCenterX := monster.x + weaponCenterOffsetX + weaponCenterOffsetX*math.Cos(weapon.angle) - weaponCenterOffsetY*math.Sin(weapon.angle)
				weaponCenterY := monster.y + weaponCenterOffsetY + weaponCenterOffsetX*math.Sin(weapon.angle) + weaponCenterOffsetY*math.Cos(weapon.angle)

				// 武器的轨迹
				weapon.Trail = append(weapon.Trail, f64.Vec2{weaponCenterX, weaponCenterY})
				if len(weapon.Trail) >= 20 {
					weapon.Trail = weapon.Trail[1:]
				}

				// 角色的中心位置
				playerCenterX := g.player.x + frameWidth/2
				playerCenterY := g.player.y + frameHeight/2
				// 并非无敌状态，且碰撞到角色，降低角色生命值
				if !g.player.Invincible() && IsTouch(weaponCenterX, weaponCenterY, playerCenterX, playerCenterY) {
					if time.Since(g.player.lastCollisionTime) < time.Second {
						continue
					}
					if err := g.hitPlayer.Rewind(); err != nil {
						return err
					}
					g.hitPlayer.Play()
					g.player.health -= 25
					g.player.lastCollisionTime = time.Now()
					if g.player.health <= 0 {
						g.mode = ModeGameOver
					}
				}
			case *RangedWeapon:
				weapon := monster.weapon.(*RangedWeapon)
				// 每秒钟发射一颗子弹
				if time.Since(weapon.LastFireTime) > time.Second {
					weapon.LastFireTime = time.Now()
					weapon.Fire(g, monster)
				}
			}
		}

		// 并非无敌状态，怪物碰撞到人物，降低生命值
		if !g.player.Invincible() && IsTouch(g.player.x, g.player.y, monster.x, monster.y) {
			if time.Since(g.player.lastCollisionTime) < time.Second {
				continue
			}
			if err := g.hitPlayer.Rewind(); err != nil {
				return err
			}
			g.hitPlayer.Play()
			g.player.health -= 25
			g.player.lastCollisionTime = time.Now()
			if g.player.health <= 0 {
				g.mode = ModeGameOver
			}
		}
	}

	return nil
}

func (g *Game) resolvePlayWeapon() error {
	if g.player.weapon != nil {
		switch g.player.weapon.(type) {
		case *MeleeWeapon:
			// 角色武器旋转
			weapon := g.player.weapon.(*MeleeWeapon)
			weapon.Spin()
			// 武器碰撞到敌人可以消灭敌人
			weaponCenterOffsetX := g.player.weaponX // 武器中心相对于角色中心的 X 坐标偏移
			weaponCenterOffsetY := g.player.weaponY // 武器中心相对于角色中心的 Y 坐标偏移
			// 考虑武器的旋转角度，将偏移向量旋转到合适的位置
			weaponCenterX := g.player.x + weaponCenterOffsetX + weaponCenterOffsetX*math.Cos(weapon.angle) - weaponCenterOffsetY*math.Sin(weapon.angle)
			weaponCenterY := g.player.y + weaponCenterOffsetY + weaponCenterOffsetX*math.Sin(weapon.angle) + weaponCenterOffsetY*math.Cos(weapon.angle)
			// 武器的轨迹
			weapon.Trail = append(weapon.Trail, f64.Vec2{weaponCenterX, weaponCenterY})
			if len(weapon.Trail) >= 20 {
				weapon.Trail = weapon.Trail[1:]
			}
			for id, monster := range g.monsters {
				// 怪物的中心位置
				monsterCenterX := monster.x + frameWidth/2
				monsterCenterY := monster.y + frameHeight/2
				if IsTouch(weaponCenterX, weaponCenterY, monsterCenterX, monsterCenterY) {
					if err := g.hitPlayer.Rewind(); err != nil {
						return err
					}
					g.hitPlayer.Play()
					g.player.score++
					delete(g.monsters, id)
					delete(g.monsterTarget, id)
					delete(g.monsterTimer, id)
				}
			}
		case *RangedWeapon:
			// TODO
		}
	}

	return nil
}

func (g *Game) resolvePlayWeapon1() error {
	// 调用角色的武器自动处理方法，让武器自身决定如何处理运动轨迹，而不在主要逻辑中处理；非自旋武器此方法为空
	g.player1.Base.Weapon.AutoResolve()

	weaponType := g.player1.Base.Weapon.GetWeaponType()
	if weaponType == weapon.SpinWeaponType {

	}

	if g.player.weapon != nil {
		switch g.player.weapon.(type) {
		case *MeleeWeapon:
			// 角色武器旋转
			weapon := g.player.weapon.(*MeleeWeapon)
			weapon.Spin()
			// 武器碰撞到敌人可以消灭敌人
			weaponCenterOffsetX := g.player.weaponX // 武器中心相对于角色中心的 X 坐标偏移
			weaponCenterOffsetY := g.player.weaponY // 武器中心相对于角色中心的 Y 坐标偏移
			// 考虑武器的旋转角度，将偏移向量旋转到合适的位置
			weaponCenterX := g.player.x + weaponCenterOffsetX + weaponCenterOffsetX*math.Cos(weapon.angle) - weaponCenterOffsetY*math.Sin(weapon.angle)
			weaponCenterY := g.player.y + weaponCenterOffsetY + weaponCenterOffsetX*math.Sin(weapon.angle) + weaponCenterOffsetY*math.Cos(weapon.angle)
			// 武器的轨迹
			weapon.Trail = append(weapon.Trail, f64.Vec2{weaponCenterX, weaponCenterY})
			if len(weapon.Trail) >= 20 {
				weapon.Trail = weapon.Trail[1:]
			}
			for id, monster := range g.monsters {
				// 怪物的中心位置
				monsterCenterX := monster.x + frameWidth/2
				monsterCenterY := monster.y + frameHeight/2
				if IsTouch(weaponCenterX, weaponCenterY, monsterCenterX, monsterCenterY) {
					if err := g.hitPlayer.Rewind(); err != nil {
						return err
					}
					g.hitPlayer.Play()
					g.player.score++
					delete(g.monsters, id)
					delete(g.monsterTarget, id)
					delete(g.monsterTimer, id)
				}
			}
		case *RangedWeapon:
			// TODO
		}
	}

	return nil
}
