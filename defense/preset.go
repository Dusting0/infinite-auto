// Package defense 防御预设模块。
// 成分与角卡「防禦預設」及规则书防御构成对齐：
// 基础 / 闪避 / 格挡 / 天生 / 盔甲(近远程) / 盾牌(近远程) / 力场 / 偏斜 / 洞察 / 掩蔽 / 自定义
// 击破分类：高速作用于基础·闪避·格挡；破甲作用于盾牌·盔甲·天生；破魔作用于力场·偏斜等魔性防御。
package defense

import (
	"fmt"
	"sync"
)

// Category 防御成分对应的击破类型
type Category string

const (
	CatSpeed Category = "speed" // 可被【高速】击破：基础、闪避、格挡
	CatArmor Category = "armor" // 可被【破甲】击破：盾牌、盔甲、天生
	CatMagic Category = "magic" // 可被【破魔】击破：力场、偏斜等
	CatOther Category = "other" // 其他（洞察、掩蔽、自定义等，默认不被高速/破甲/破魔直接击破）
)

// CustomComp 用户自定义防御成分
type CustomComp struct {
	Name     string   `json:"name"`
	Value    int      `json:"value"`
	Category Category `json:"category"` // speed/armor/magic/other
	// Scope: both | melee | ranged
	Scope string `json:"scope"`
	// OnlyWhenBlocking 仅格挡时生效（类似盾牌/格挡防御）
	OnlyWhenBlocking bool `json:"onlyWhenBlocking"`
	// LostOnFlatFooted 措手不及时失去
	LostOnFlatFooted bool `json:"lostOnFlatFooted"`
}

// Preset 一份防御预设
type Preset struct {
	mu sync.Mutex
	presetData
}

// presetData 包含 Preset 的全部数据字段（不含 mu）。独立出来是为了让 Reset
// 可以整体复制数据而不触碰已上锁的 mu——直接 *p = *NewPreset(name) 会用新预设
// 零值的 mu 覆盖当前已上锁的 mu，导致 defer Unlock 解锁一把未上锁的 mutex 而 panic。
type presetData struct {
	Name string `json:"name"`

	// 状态开关
	FlatFooted  bool `json:"flatFooted"`  // 措手不及
	FullDefense bool `json:"fullDefense"` // 全力防御（基础防御翻倍）
	Blocking    bool `json:"blocking"`    // 是否处于格挡
	Cover       bool `json:"cover"`       // 是否有掩蔽
	TouchAttack bool `json:"touchAttack"` // 接触攻击：检定无视天生/盔甲/盾牌

	// 标准成分
	Base         int `json:"base"`         // 基础防御
	Dodge        int `json:"dodge"`        // 闪避防御
	Block        int `json:"block"`        // 格挡防御（仅格挡时）
	Natural      int `json:"natural"`      // 天生防御
	ArmorMelee   int `json:"armorMelee"`   // 盔甲（近战）
	ArmorRanged  int `json:"armorRanged"`  // 盔甲（远程）
	ShieldMelee  int `json:"shieldMelee"`  // 盾牌（近战，仅格挡）
	ShieldRanged int `json:"shieldRanged"` // 盾牌（远程，仅格挡）
	Force        int `json:"force"`        // 力场
	Deflection   int `json:"deflection"`   // 偏斜
	Insight      int `json:"insight"`      // 洞察
	CoverBonus   int `json:"coverBonus"`   // 掩蔽加值（开启掩蔽时计入，角卡默认 4）

	Other2     int    `json:"other2"`     // 其他1 数值（字段名历史兼容）
	Other3     int    `json:"other3"`     // 其他2 数值
	Other2Name string `json:"other2Name"` // 其他1 显示名
	Other3Name string `json:"other3Name"` // 其他2 显示名

	// 完美防御：不被高速/破甲/破魔击破；不计入三池合计；勾选生效时直接从攻击 DP 扣除
	PerfectDefense       int  `json:"perfectDefense"`
	PerfectDefenseActive bool `json:"perfectDefenseActive"`

	// 抵消击破（保留字段以兼容已有单测与旧数据，UI 不再展示）
	ResistSpeed int `json:"resistSpeed"` // 抵高速
	ResistAP    int `json:"resistAP"`    // 抵破甲
	ResistMagic int `json:"resistMagic"` // 抵破魔

	// 防御附加成功 / 伤害相关
	DefenseBonusSuccess int    `json:"defenseBonusSuccess"` // 防御附加成功
	DamageAbsorb        int    `json:"damageAbsorb"`        // 伤害吸收数值
	DamageAbsorbType    string `json:"damageAbsorbType"`    // physical | energy | all
	DRValue             int    `json:"drValue"`             // 物理伤害减免数值（DR X）
	DRType              string `json:"drType"`              // DR 弱点枚举：""=DR/- 无弱点 | magic | divine | slashing | piercing | bludgeoning
	ERValue             int    `json:"erValue"`             // 能量抗力数值（ER X）
	ERType              string `json:"erType"`              // 能量子类型枚举：all=全能量抗力 | pure | fire | cold | lightning | corrosion | light | dark | sonic | luminous

	Customs []CustomComp `json:"customs"`
	Notes   string       `json:"notes"`
	Traits  []string     `json:"traits"` // 特性列表（多行，可排序）
}

func otherName(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// 弱点/子类型/吸收类型的合法枚举。未知值归一：DRType→""（DR/-，对所有物理生效）、
// ERType→"all"（全能量抗力，保持旧行为）、DamageAbsorbType→"physical"。
var validDRTypes = map[string]bool{"": true, "magic": true, "divine": true, "slashing": true, "piercing": true, "bludgeoning": true}
var validERTypes = map[string]bool{"": true, "all": true, "pure": true, "fire": true, "cold": true, "lightning": true, "corrosion": true, "light": true, "dark": true, "sonic": true, "luminous": true}
var validAbsorbTypes = map[string]bool{"": true, "physical": true, "energy": true, "all": true}

func normalizeDRType(s string) string {
	if validDRTypes[s] {
		return s
	}
	return ""
}
func normalizeERType(s string) string {
	if validERTypes[s] {
		return s
	}
	return "all"
}
func normalizeAbsorbType(s string) string {
	if validAbsorbTypes[s] {
		return s
	}
	return "physical"
}

func NewPreset(name string) *Preset {
	if name == "" {
		name = "默认预设"
	}
	return &Preset{
		presetData: presetData{
			Name:             name,
			CoverBonus:       4, // 角卡掩蔽默认 +4
			Other2Name:       "其他1",
			Other3Name:       "其他2",
			DamageAbsorbType: "physical",
			ERType:           "all", // 默认全能量抗力（对所有能量生效）
			Customs:          []CustomComp{},
		},
	}
}

// effectiveBase 计算生效的基础防御（措手不及→0；全力防御→翻倍）
func (p *Preset) effectiveBase() int {
	if p.FlatFooted {
		return 0
	}
	v := p.Base
	if p.FullDefense {
		v *= 2
	}
	return v
}

func (p *Preset) effectiveDodge() int {
	if p.FlatFooted {
		return 0
	}
	return p.Dodge
}

func (p *Preset) effectiveBlock() int {
	if p.FlatFooted || !p.Blocking {
		return 0
	}
	return p.Block
}

func (p *Preset) effectiveShieldMelee() int {
	if !p.Blocking {
		return 0
	}
	return p.ShieldMelee
}

func (p *Preset) effectiveShieldRanged() int {
	if !p.Blocking {
		return 0
	}
	return p.ShieldRanged
}

func (p *Preset) effectiveCover() int {
	if !p.Cover {
		return 0
	}
	return p.CoverBonus
}

// Totals 近战/远程防御合计，以及按击破类型分类的合计（便于对照高速/破甲/破魔）
type Totals struct {
	Melee  int `json:"melee"`
	Ranged int `json:"ranged"`

	// 可被高速击破的部分（近战视角，格挡已按状态计入）
	SpeedPoolMelee  int `json:"speedPoolMelee"`
	SpeedPoolRanged int `json:"speedPoolRanged"`
	// 可被破甲击破
	ArmorPoolMelee  int `json:"armorPoolMelee"`
	ArmorPoolRanged int `json:"armorPoolRanged"`
	// 可被破魔击破
	MagicPool int `json:"magicPool"`
	// 其他
	OtherPool int `json:"otherPool"`
}

func (p *Preset) Compute() Totals {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.computeLocked()
}

func (p *Preset) computeLocked() Totals {
	base := p.effectiveBase()
	dodge := p.effectiveDodge()
	block := p.effectiveBlock()
	nat := p.Natural
	am, ar := p.ArmorMelee, p.ArmorRanged
	sm, sr := p.effectiveShieldMelee(), p.effectiveShieldRanged()
	force, defl := p.Force, p.Deflection
	insight := p.Insight
	cover := p.effectiveCover()
	o2, o3 := p.Other2, p.Other3

	speedM := base + dodge + block
	speedR := base + dodge + block
	armorM := sm + am + nat
	armorR := sr + ar + nat
	magic := force + defl
	other := insight + cover + o2 + o3

	// 接触攻击：检定无视天生/盔甲/盾牌，合计中不计入破甲池
	armorMEff, armorREff := armorM, armorR
	if p.TouchAttack {
		armorMEff, armorREff = 0, 0
	}

	melee := speedM + armorMEff + magic + other
	ranged := speedR + armorREff + magic + other

	for _, c := range p.Customs {
		v := c.Value
		if c.LostOnFlatFooted && p.FlatFooted {
			v = 0
		}
		if c.OnlyWhenBlocking && !p.Blocking {
			v = 0
		}
		scope := c.Scope
		if scope == "" {
			scope = "both"
		}
		addM, addR := scope != "ranged", scope != "melee"
		if scope == "melee" {
			addM, addR = true, false
		} else if scope == "ranged" {
			addM, addR = false, true
		}
		switch c.Category {
		case CatSpeed:
			if addM {
				speedM += v
				melee += v
			}
			if addR {
				speedR += v
				ranged += v
			}
		case CatArmor:
			if addM {
				armorM += v
				if !p.TouchAttack {
					melee += v
				}
			}
			if addR {
				armorR += v
				if !p.TouchAttack {
					ranged += v
				}
			}
		case CatMagic:
			magic += v
			if addM {
				melee += v
			}
			if addR {
				ranged += v
			}
		default:
			other += v
			if addM {
				melee += v
			}
			if addR {
				ranged += v
			}
		}
	}

	return Totals{
		Melee:           melee,
		Ranged:          ranged,
		SpeedPoolMelee:  speedM,
		SpeedPoolRanged: speedR,
		ArmorPoolMelee:  armorM,
		ArmorPoolRanged: armorR,
		MagicPool:       magic,
		OtherPool:       other,
	}
}

// Snapshot 供 API / 前端展示
func (p *Preset) Snapshot() map[string]interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	t := p.computeLocked()
	customs := make([]CustomComp, len(p.Customs))
	copy(customs, p.Customs)
	return map[string]interface{}{
		"name":         p.Name,
		"flatFooted":   p.FlatFooted,
		"fullDefense":  p.FullDefense,
		"blocking":     p.Blocking,
		"cover":        p.Cover,
		"touchAttack":  p.TouchAttack,
		"base":         p.Base,
		"dodge":        p.Dodge,
		"block":        p.Block,
		"natural":      p.Natural,
		"armorMelee":   p.ArmorMelee,
		"armorRanged":  p.ArmorRanged,
		"shieldMelee":  p.ShieldMelee,
		"shieldRanged": p.ShieldRanged,
		"force":        p.Force,
		"deflection":   p.Deflection,
		"insight":      p.Insight,
		"coverBonus":   p.CoverBonus,
		"other2":               p.Other2,
		"other3":               p.Other3,
		"other2Name":           otherName(p.Other2Name, "其他1"),
		"other3Name":           otherName(p.Other3Name, "其他2"),
		"perfectDefense":       p.PerfectDefense,
		"perfectDefenseActive": p.PerfectDefenseActive,
		"resistSpeed":  p.ResistSpeed,
		"resistAP":     p.ResistAP,
		"resistMagic":  p.ResistMagic,
		"defenseBonusSuccess": p.DefenseBonusSuccess,
		"damageAbsorb":        p.DamageAbsorb,
		"damageAbsorbType":    p.DamageAbsorbType,
		"drValue":             p.DRValue,
		"drType":              p.DRType,
		"erValue":             p.ERValue,
		"erType":              p.ERType,
		"customs":      customs,
		"notes":        p.Notes,
		"traits":       p.Traits,
		"totals":       t,
		// 生效后的分项（便于 UI 展示）
		"effective": map[string]int{
			"base":         p.effectiveBase(),
			"dodge":        p.effectiveDodge(),
			"block":        p.effectiveBlock(),
			"shieldMelee":  p.effectiveShieldMelee(),
			"shieldRanged": p.effectiveShieldRanged(),
			"cover":        p.effectiveCover(),
		},
	}
}

// Update 用前端提交的字段整体更新（保留未传字段需由前端全量提交）
func (p *Preset) Update(in *Preset) error {
	if in == nil {
		return fmt.Errorf("参数为空")
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if in.Name != "" {
		p.Name = in.Name
	}
	p.FlatFooted = in.FlatFooted
	p.FullDefense = in.FullDefense
	p.Blocking = in.Blocking
	p.Cover = in.Cover
	p.TouchAttack = in.TouchAttack

	p.Base = in.Base
	p.Dodge = in.Dodge
	p.Block = in.Block
	p.Natural = in.Natural
	p.ArmorMelee = in.ArmorMelee
	p.ArmorRanged = in.ArmorRanged
	p.ShieldMelee = in.ShieldMelee
	p.ShieldRanged = in.ShieldRanged
	p.Force = in.Force
	p.Deflection = in.Deflection
	p.Insight = in.Insight
	if in.CoverBonus > 0 {
		p.CoverBonus = in.CoverBonus
	}
	p.Other2 = in.Other2
	p.Other3 = in.Other3
	if in.Other2Name != "" {
		p.Other2Name = in.Other2Name
	}
	if in.Other3Name != "" {
		p.Other3Name = in.Other3Name
	}
	p.PerfectDefense = in.PerfectDefense
	p.PerfectDefenseActive = in.PerfectDefenseActive
	p.ResistSpeed = in.ResistSpeed
	p.ResistAP = in.ResistAP
	p.ResistMagic = in.ResistMagic
	p.DefenseBonusSuccess = in.DefenseBonusSuccess
	p.DamageAbsorb = in.DamageAbsorb
	p.DamageAbsorbType = normalizeAbsorbType(in.DamageAbsorbType)
	p.DRValue = in.DRValue
	p.DRType = normalizeDRType(in.DRType)
	p.ERValue = in.ERValue
	p.ERType = normalizeERType(in.ERType)
	p.Notes = in.Notes
	if in.Traits != nil {
		p.Traits = append([]string(nil), in.Traits...)
	}
	if in.Customs != nil {
		p.Customs = append([]CustomComp(nil), in.Customs...)
	}
	return nil
}

func (p *Preset) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	name := p.Name
	// 只复制数据字段，保留已上锁的 mu。presetData 不含 mutex，整体复制安全。
	p.presetData = NewPreset(name).presetData
}
