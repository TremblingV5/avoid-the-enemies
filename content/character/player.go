package character

import "time"

type Player struct {
	Base       Character
	IsSkill    bool      // 是否释放技能
	SkillFrame int       // 技能的帧数
	SkillTime  time.Time // 技能释放的时间
	StartTime  time.Time // 游戏开始的时间
}
