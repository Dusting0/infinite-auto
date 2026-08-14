package defense

import (
	"infinite-calc/dice"
)

// AttackInput 自定义伤害输入（敌方攻击）
type AttackInput struct {
	Ranged       bool   `json:"ranged"`       // false=近战 true=远程
	AttackDP     int    `json:"attackDP"`     // 敌方攻击 DP
	Speed        int    `json:"speed"`        // 高速
	ArmorPierce  int    `json:"armorPierce"`  // 破甲
	MagicPierce  int    `json:"magicPierce"`  // 破魔
	ExplodeOn    int    `json:"explodeOn"`    // 加骰 8/9/10
	BonusSuccess int    `json:"bonusSuccess"` // 攻击附加成功（可负）
	IsPhysical   bool   `json:"isPhysical"`   // 兼容旧前端：true=物理，false=能量
	DamageKind   string `json:"damageKind"`   // physical | energy | mixed
	DamageLimit  int    `json:"damageLimit"`  // 伤害上限；≤0 不生效
}

// PoolBreakdown 击破前后分项（便于 UI 展示）
type PoolBreakdown struct {
	// 高速池顺序：格挡、闪避、基础
	BlockBefore int `json:"blockBefore"`
	BlockAfter  int `json:"blockAfter"`
	DodgeBefore int `json:"dodgeBefore"`
	DodgeAfter  int `json:"dodgeAfter"`
	BaseBefore  int `json:"baseBefore"`
	BaseAfter   int `json:"baseAfter"`
	// 破甲池顺序：盾牌、盔甲、天生
	ShieldBefore  int `json:"shieldBefore"`
	ShieldAfter   int `json:"shieldAfter"`
	ArmorBefore   int `json:"armorBefore"`
	ArmorAfter    int `json:"armorAfter"`
	NaturalBefore int `json:"naturalBefore"`
	NaturalAfter  int `json:"naturalAfter"`
	// 破魔池顺序：力场、偏斜、洞察、掩蔽、其他2、其他3
	ForceBefore      int `json:"forceBefore"`
	ForceAfter       int `json:"forceAfter"`
	DeflectionBefore int `json:"deflectionBefore"`
	DeflectionAfter  int `json:"deflectionAfter"`
	InsightBefore    int `json:"insightBefore"`
	InsightAfter     int `json:"insightAfter"`
	CoverBefore      int `json:"coverBefore"`
	CoverAfter       int `json:"coverAfter"`
	Other2Before     int `json:"other2Before"`
	Other2After      int `json:"other2After"`
	Other3Before     int `json:"other3Before"`
	Other3After      int `json:"other3After"`

	SpeedPoolBefore int `json:"speedPoolBefore"`
	SpeedPoolAfter  int `json:"speedPoolAfter"`
	ArmorPoolBefore int `json:"armorPoolBefore"`
	ArmorPoolAfter  int `json:"armorPoolAfter"`
	MagicPoolBefore int `json:"magicPoolBefore"`
	MagicPoolAfter  int `json:"magicPoolAfter"`
}

// AttackResult 自定义伤害结算结果（不写入血量）
type AttackResult struct {
	Miss             bool   `json:"miss"`
	MissReason       string `json:"missReason,omitempty"`
	EffectiveDefense int    `json:"effectiveDefense"`
	ActualDP         int    `json:"actualDP"`
	// 抵消后实际生效的击破值
	EffSpeed             int           `json:"effSpeed"`
	EffArmorPierce       int           `json:"effArmorPierce"`
	EffMagicPierce       int           `json:"effMagicPierce"`
	ResistSpeed          int           `json:"resistSpeed"`
	ResistAP             int           `json:"resistAP"`
	ResistMagic          int           `json:"resistMagic"`
	TouchAttack          bool          `json:"touchAttack"`
	PerfectDefense       int           `json:"perfectDefense"`
	PerfectDefenseActive bool          `json:"perfectDefenseActive"`
	Pools                PoolBreakdown `json:"pools"`

	// 掷骰（仅未因 DP<=0 未命中时有）
	Roll *dice.Result `json:"roll,omitempty"`

	NaturalSuccess int `json:"naturalSuccess"`
	FinalSuccess   int `json:"finalSuccess"`
	DefenseBonus   int `json:"defenseBonus"`

	RawDamage     int    `json:"rawDamage"`   // 掷骰最终成功数（上限前）
	AfterLimit    int    `json:"afterLimit"`  // 伤害上限封顶后
	DamageLimit   int    `json:"damageLimit"` // 本次使用的上限（≤0 表示未启用）
	AfterDR       int    `json:"afterDR"`
	AfterAbsorb   int    `json:"afterAbsorb"`
	FinalDamage   int    `json:"finalDamage"`
	DRValue       int    `json:"drValue"`
	DRType        string `json:"drType"`
	ERValue       int    `json:"erValue"`
	ERType        string `json:"erType"`
	AbsorbApplied int    `json:"absorbApplied"`
	AbsorbType    string `json:"absorbType"`
	Summary       string `json:"summary"`
}

// subtract 按顺序从一组值中扣除 pierce，返回 after 切片
func subtractOrdered(pierce int, vals []int) (int, []int) {
	after := make([]int, len(vals))
	copy(after, vals)
	left := pierce
	if left < 0 {
		left = 0
	}
	for i := range after {
		if left <= 0 {
			break
		}
		if after[i] >= left {
			after[i] -= left
			left = 0
		} else {
			left -= after[i]
			after[i] = 0
		}
	}
	return left, after
}

func sum(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func damageKind(in AttackInput) string {
	switch in.DamageKind {
	case "physical", "energy", "mixed":
		return in.DamageKind
	}
	if in.IsPhysical {
		return "physical"
	}
	return "energy"
}

func damageReduction(kind string, drValue, erValue int) int {
	switch kind {
	case "physical":
		return drValue
	case "mixed":
		return min(drValue, erValue)
	default:
		return erValue
	}
}

// ResolveAttack 根据当前防御预设结算一次攻击，不修改预设、不写血量。
func (p *Preset) ResolveAttack(in AttackInput) AttackResult {
	p.mu.Lock()
	defer p.mu.Unlock()

	base := p.effectiveBase()
	dodge := p.effectiveDodge()
	block := p.effectiveBlock()
	nat := p.Natural
	// 不再区分近战/远程：统一使用一套盔甲/盾牌数值
	armor := p.ArmorMelee
	shield := p.effectiveShieldMelee()
	force, defl := p.Force, p.Deflection
	insight := p.Insight
	cover := p.effectiveCover()
	o2, o3 := p.Other2, p.Other3

	// 抵高速/破甲/破魔：先削减进攻方击破值，剩余再击破对应防御池
	effSpeed := in.Speed - p.ResistSpeed
	if effSpeed < 0 {
		effSpeed = 0
	}
	effAP := in.ArmorPierce - p.ResistAP
	if effAP < 0 {
		effAP = 0
	}
	effMP := in.MagicPierce - p.ResistMagic
	if effMP < 0 {
		effMP = 0
	}

	// 高速：格挡 → 闪避 → 基础
	speedBefore := []int{block, dodge, base}
	_, speedAfter := subtractOrdered(effSpeed, speedBefore)

	// 破甲：盾牌 → 盔甲 → 天生
	armorBefore := []int{shield, armor, nat}
	_, armorAfter := subtractOrdered(effAP, armorBefore)

	// 破魔：力场 → 偏斜 → 洞察 → 掩蔽 → 其他2 → 其他3
	magicBefore := []int{force, defl, insight, cover, o2, o3}
	_, magicAfter := subtractOrdered(effMP, magicBefore)

	pools := PoolBreakdown{
		BlockBefore: block, BlockAfter: speedAfter[0],
		DodgeBefore: dodge, DodgeAfter: speedAfter[1],
		BaseBefore: base, BaseAfter: speedAfter[2],
		ShieldBefore: shield, ShieldAfter: armorAfter[0],
		ArmorBefore: armor, ArmorAfter: armorAfter[1],
		NaturalBefore: nat, NaturalAfter: armorAfter[2],
		ForceBefore: force, ForceAfter: magicAfter[0],
		DeflectionBefore: defl, DeflectionAfter: magicAfter[1],
		InsightBefore: insight, InsightAfter: magicAfter[2],
		CoverBefore: cover, CoverAfter: magicAfter[3],
		Other2Before: o2, Other2After: magicAfter[4],
		Other3Before: o3, Other3After: magicAfter[5],
		SpeedPoolBefore: sum(speedBefore), SpeedPoolAfter: sum(speedAfter),
		ArmorPoolBefore: sum(armorBefore), ArmorPoolAfter: sum(armorAfter),
		MagicPoolBefore: sum(magicBefore), MagicPoolAfter: sum(magicAfter),
	}

	// 接触攻击：检定无视破甲池；破甲仍如常击破该池（防弹等）
	armorForDef := pools.ArmorPoolAfter
	if p.TouchAttack {
		armorForDef = 0
	}
	effDef := pools.SpeedPoolAfter + armorForDef + pools.MagicPoolAfter

	// 完美防御：不被击破，勾选生效时直接从攻击 DP 扣除（不计入三池）
	perfect := 0
	if p.PerfectDefenseActive && p.PerfectDefense > 0 {
		perfect = p.PerfectDefense
	}
	actualDP := in.AttackDP - effDef - perfect

	res := AttackResult{
		EffectiveDefense:     effDef,
		ActualDP:             actualDP,
		PerfectDefense:       perfect,
		PerfectDefenseActive: p.PerfectDefenseActive,
		EffSpeed:             effSpeed,
		EffArmorPierce:       effAP,
		EffMagicPierce:       effMP,
		ResistSpeed:          p.ResistSpeed,
		ResistAP:             p.ResistAP,
		ResistMagic:          p.ResistMagic,
		TouchAttack:          p.TouchAttack,
		Pools:                pools,
		DefenseBonus:         p.DefenseBonusSuccess,
		DRValue:              p.DRValue,
		DRType:               p.DRType,
		ERValue:              p.ERValue,
		ERType:               p.ERType,
		AbsorbType:           p.DamageAbsorbType,
	}
	kind := damageKind(in)

	// DP≤0：直接未命中（不走机运骰）
	if actualDP <= 0 {
		res.Miss = true
		res.MissReason = "实际DP≤0"
		res.Summary = "未命中（实际DP≤0）"
		return res
	}

	roll := dice.Roll(actualDP, in.ExplodeOn, in.BonusSuccess)
	res.Roll = &roll
	res.NaturalSuccess = roll.NaturalSuccess
	res.FinalSuccess = roll.FinalSuccess

	// 防御附加成功 > 自然成功 → 未命中（规则书）
	if p.DefenseBonusSuccess > roll.NaturalSuccess {
		res.Miss = true
		res.MissReason = "防御附加成功大于自然成功"
		res.Summary = "未命中（防御附加成功）"
		return res
	}

	// 命中后：最终伤害先减去防御附加成功
	raw := roll.FinalSuccess - p.DefenseBonusSuccess
	if raw < 0 {
		raw = 0
	}
	res.RawDamage = raw

	// 伤害上限：在 DR / 吸收之前封顶；≤0 不生效
	capped := raw
	res.DamageLimit = in.DamageLimit
	if in.DamageLimit > 0 && capped > in.DamageLimit {
		capped = in.DamageLimit
	}
	res.AfterLimit = capped

	// 物理 → DR；能量 → ER；混合 → 吃较低的减免，保留更高伤害。
	reduction := damageReduction(kind, p.DRValue, p.ERValue)
	reduced := capped - reduction
	if reduced < 0 {
		reduced = 0
	}
	res.AfterDR = reduced

	// 伤害吸收：全伤害始终；物理吸收对物理和混合伤害生效。
	absorb := 0
	if p.DamageAbsorb > 0 {
		if p.DamageAbsorbType == "all" || (p.DamageAbsorbType == "physical" && kind != "energy") {
			absorb = p.DamageAbsorb
		}
	}
	if absorb > reduced {
		absorb = reduced
	}
	afterAbs := reduced - absorb
	res.AbsorbApplied = absorb
	res.AfterAbsorb = afterAbs
	res.FinalDamage = afterAbs
	res.Summary = roll.Summary
	return res
}
