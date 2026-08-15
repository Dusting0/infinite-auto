package defense

import "testing"

func TestResolveMissWhenDPNonPositive(t *testing.T) {
	p := NewPreset("t")
	p.Base = 10
	r := p.ResolveAttack(AttackInput{AttackDP: 5, IsPhysical: true})
	if !r.Miss || r.MissReason != "实际DP≤0" {
		t.Fatalf("期望实际DP≤0未命中，得 miss=%v reason=%s", r.Miss, r.MissReason)
	}
	if r.FinalDamage != 0 {
		t.Fatalf("期望伤害0，得%d", r.FinalDamage)
	}
	if r.Roll != nil {
		t.Fatal("DP≤0 不应掷骰")
	}
}

func TestResolvePierceOrderSpeed(t *testing.T) {
	p := NewPreset("t")
	p.Blocking = true
	p.Base = 3
	p.Dodge = 2
	p.Block = 4
	// 高速 5：格挡4→0，闪避2→1，基础3→3
	r := p.ResolveAttack(AttackInput{AttackDP: 100, Speed: 5, IsPhysical: true})
	if r.Pools.BlockAfter != 0 || r.Pools.DodgeAfter != 1 || r.Pools.BaseAfter != 3 {
		t.Fatalf("高速顺序错误: block=%d dodge=%d base=%d",
			r.Pools.BlockAfter, r.Pools.DodgeAfter, r.Pools.BaseAfter)
	}
}

func TestResolvePierceOrderArmor(t *testing.T) {
	p := NewPreset("t")
	p.Blocking = true
	p.ShieldMelee = 3
	p.ArmorMelee = 2
	p.Natural = 1
	// 破甲 4：盾3→0，甲2→1，天1→1
	r := p.ResolveAttack(AttackInput{AttackDP: 100, ArmorPierce: 4, IsPhysical: true})
	if r.Pools.ShieldAfter != 0 || r.Pools.ArmorAfter != 1 || r.Pools.NaturalAfter != 1 {
		t.Fatalf("破甲顺序错误: shield=%d armor=%d nat=%d",
			r.Pools.ShieldAfter, r.Pools.ArmorAfter, r.Pools.NaturalAfter)
	}
}

func TestResolvePierceMagicIncludesInsight(t *testing.T) {
	p := NewPreset("t")
	p.Force = 1
	p.Deflection = 1
	p.Insight = 2
	p.Other2 = 1
	// 破魔 3：力场1→0，偏斜1→0，洞察2→1
	r := p.ResolveAttack(AttackInput{AttackDP: 100, MagicPierce: 3, IsPhysical: true})
	if r.Pools.ForceAfter != 0 || r.Pools.DeflectionAfter != 0 || r.Pools.InsightAfter != 1 {
		t.Fatalf("破魔顺序错误: force=%d defl=%d insight=%d",
			r.Pools.ForceAfter, r.Pools.DeflectionAfter, r.Pools.InsightAfter)
	}
	if r.Pools.MagicPoolBefore != 5 {
		t.Fatalf("破魔池期望5，得%d", r.Pools.MagicPoolBefore)
	}
}

func TestResolveDefenseBonusMiss(t *testing.T) {
	p := NewPreset("t")
	p.DefenseBonusSuccess = 100
	r := p.ResolveAttack(AttackInput{AttackDP: 3, ExplodeOn: 10, IsPhysical: true})
	if !r.Miss {
		t.Fatalf("防御附加100应未命中，自然=%d", r.NaturalSuccess)
	}
	if r.MissReason != "防御附加成功大于自然成功" {
		t.Fatalf("原因错误: %s", r.MissReason)
	}
}

// 抵高速只削减进攻方高速，不直接减 DP
// DP50 高速5，仅抵高速10 → 有效高速0，无防御池 → 实际DP50
func TestResistSpeedOnlyNoPool(t *testing.T) {
	p := NewPreset("t")
	p.ResistSpeed = 10
	r := p.ResolveAttack(AttackInput{AttackDP: 50, Speed: 5, IsPhysical: true})
	if r.EffSpeed != 0 {
		t.Fatalf("有效高速期望0，得%d", r.EffSpeed)
	}
	if r.ActualDP != 50 {
		t.Fatalf("实际DP期望50，得%d", r.ActualDP)
	}
}

// 高速防御10、无抵高速、进攻高速5 → 池剩5 → 实际DP45
func TestSpeedPoolPierced(t *testing.T) {
	p := NewPreset("t")
	p.Base = 10
	r := p.ResolveAttack(AttackInput{AttackDP: 50, Speed: 5, IsPhysical: true})
	if r.Pools.SpeedPoolAfter != 5 {
		t.Fatalf("高速池剩余期望5，得%d", r.Pools.SpeedPoolAfter)
	}
	if r.ActualDP != 45 {
		t.Fatalf("实际DP期望45，得%d", r.ActualDP)
	}
}

// 高速防御10 + 抵高速3、进攻高速5 → 有效高速2 → 池剩8 → 实际DP42
func TestResistThenPierce(t *testing.T) {
	p := NewPreset("t")
	p.Base = 10
	p.ResistSpeed = 3
	r := p.ResolveAttack(AttackInput{AttackDP: 50, Speed: 5, IsPhysical: true})
	if r.EffSpeed != 2 {
		t.Fatalf("有效高速期望2，得%d", r.EffSpeed)
	}
	if r.Pools.SpeedPoolAfter != 8 {
		t.Fatalf("高速池剩余期望8，得%d", r.Pools.SpeedPoolAfter)
	}
	if r.ActualDP != 42 {
		t.Fatalf("实际DP期望42，得%d", r.ActualDP)
	}
}

// 伤害上限：成功数再高也不能超过上限；上限≤0 不生效
func TestDamageLimit(t *testing.T) {
	p := NewPreset("t")
	r := p.ResolveAttack(AttackInput{
		AttackDP: 5, ExplodeOn: 10, BonusSuccess: 20, IsPhysical: true, DamageLimit: 3,
	})
	if r.Miss {
		t.Fatalf("不应未命中: %s", r.MissReason)
	}
	if r.AfterLimit != 3 {
		t.Fatalf("上限后期望3，得 raw=%d afterLimit=%d", r.RawDamage, r.AfterLimit)
	}
	if r.FinalDamage > 3 {
		t.Fatalf("最终伤害不应超过上限，得%d", r.FinalDamage)
	}

	r2 := p.ResolveAttack(AttackInput{
		AttackDP: 5, ExplodeOn: 10, BonusSuccess: 20, IsPhysical: true, DamageLimit: 0,
	})
	if !r2.Miss && r2.AfterLimit != r2.RawDamage {
		t.Fatalf("上限≤0 应不封顶: raw=%d afterLimit=%d", r2.RawDamage, r2.AfterLimit)
	}
}

// 接触攻击：无视破甲池，高速/破魔仍计入
func TestTouchAttackIgnoresArmorPool(t *testing.T) {
	p := NewPreset("t")
	p.Base = 4
	p.Natural = 5
	p.ArmorMelee = 3
	p.Force = 2
	p.TouchAttack = true
	r := p.ResolveAttack(AttackInput{AttackDP: 20, IsPhysical: true})
	// 有效防御 = 高速4 + 破魔2 = 6（破甲池8不计）
	if r.EffectiveDefense != 6 {
		t.Fatalf("接触攻击有效防御期望6，得%d", r.EffectiveDefense)
	}
	if r.ActualDP != 14 {
		t.Fatalf("实际DP期望14，得%d", r.ActualDP)
	}
	// 破甲仍会击破破甲池（展示用）
	r2 := p.ResolveAttack(AttackInput{AttackDP: 20, ArmorPierce: 4, IsPhysical: true})
	if r2.Pools.ArmorPoolAfter != 4 { // 5+3-4=4
		t.Fatalf("破甲应如常击破池: after=%d", r2.Pools.ArmorPoolAfter)
	}
	if r2.EffectiveDefense != 6 {
		t.Fatalf("接触下破甲不应改变有效防御，得%d", r2.EffectiveDefense)
	}
}

func TestTouchAttackTotals(t *testing.T) {
	p := NewPreset("t")
	p.Base = 3
	p.Natural = 4
	p.ArmorMelee = 2
	tot := p.Compute()
	if tot.Melee != 9 {
		t.Fatalf("非接触合计期望9，得%d", tot.Melee)
	}
	p.TouchAttack = true
	tot = p.Compute()
	if tot.Melee != 3 {
		t.Fatalf("接触合计期望3（仅高速），得%d", tot.Melee)
	}
}

// 防御附加成功：命中后从伤害中扣除
func TestDefenseBonusSubtractsDamage(t *testing.T) {
	p := NewPreset("t")
	p.DefenseBonusSuccess = 1
	// 大 DP 保证自然成功 ≥ 1，从而命中并扣除防御附加
	r := p.ResolveAttack(AttackInput{
		AttackDP: 25, ExplodeOn: 10, BonusSuccess: 4, IsPhysical: true,
	})
	if r.Miss {
		t.Fatalf("不应未命中: %s nat=%d", r.MissReason, r.NaturalSuccess)
	}
	if r.RawDamage != r.FinalSuccess-1 {
		t.Fatalf("原始伤害期望 finalSuccess-1=%d，得%d (final=%d)", r.FinalSuccess-1, r.RawDamage, r.FinalSuccess)
	}
}

func TestMixedDamageUsesLowerReduction(t *testing.T) {
	p := NewPreset("t")
	p.DRValue = 5
	p.ERValue = 2
	p.DamageAbsorb = 3
	p.DamageAbsorbType = "physical"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "mixed", DamageLimit: 10})
	if r.Miss {
		t.Fatalf("不应未命中: %s", r.MissReason)
	}
	if r.AfterLimit != 10 {
		t.Fatalf("上限后期望10，得%d", r.AfterLimit)
	}
	if r.AfterDR != 8 {
		t.Fatalf("混合伤害应吃较低减免2，减免后期望8，得%d", r.AfterDR)
	}
	if r.AbsorbApplied != 0 || r.FinalDamage != 8 {
		t.Fatalf("混合伤害下物理吸收不应生效，得 absorb=%d final=%d", r.AbsorbApplied, r.FinalDamage)
	}
}

// 混合伤害 + 全伤害吸收：吸收生效（与物理吸收对混合不生效对照）。
func TestMixedDamageWithAllAbsorb(t *testing.T) {
	p := NewPreset("t")
	p.DRValue = 5
	p.ERValue = 2
	p.DamageAbsorb = 3
	p.DamageAbsorbType = "all"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "mixed", DamageLimit: 10})
	if r.AfterDR != 8 {
		t.Fatalf("减免后期望8，得%d", r.AfterDR)
	}
	if r.AbsorbApplied != 3 || r.FinalDamage != 5 {
		t.Fatalf("全伤害吸收对混合应生效3，得 absorb=%d final=%d", r.AbsorbApplied, r.FinalDamage)
	}
}

// 能量伤害 + 能量吸收：吸收生效。
func TestEnergyDamageWithEnergyAbsorb(t *testing.T) {
	p := NewPreset("t")
	p.DamageAbsorb = 3
	p.DamageAbsorbType = "energy"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "energy", EnergyType: "fire", DamageLimit: 10})
	if r.AbsorbApplied != 3 || r.FinalDamage != 7 {
		t.Fatalf("能量吸收对能量伤害应生效3，得 absorb=%d final=%d", r.AbsorbApplied, r.FinalDamage)
	}
}

// 能量伤害 + 全伤害吸收：吸收生效。
func TestEnergyDamageWithAllAbsorb(t *testing.T) {
	p := NewPreset("t")
	p.DamageAbsorb = 3
	p.DamageAbsorbType = "all"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "energy", EnergyType: "fire", DamageLimit: 10})
	if r.AbsorbApplied != 3 || r.FinalDamage != 7 {
		t.Fatalf("全伤害吸收对能量伤害应生效3，得 absorb=%d final=%d", r.AbsorbApplied, r.FinalDamage)
	}
}

// DR/魔法：攻击带【魔法】特性时穿透，否则生效。
func TestDRMagicPiercedByMagic(t *testing.T) {
	p := NewPreset("t")
	p.DRValue = 5
	p.DRType = "magic"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, IsPhysical: true, Magic: true, DamageLimit: 10})
	if r.DREff != 0 || r.AfterDR != 10 {
		t.Fatalf("魔法攻击应穿透 DR/魔法，得 drEff=%d afterDR=%d", r.DREff, r.AfterDR)
	}
	r2 := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, IsPhysical: true, DamageLimit: 10})
	if r2.DREff != 5 || r2.AfterDR != 5 {
		t.Fatalf("非魔法攻击不应穿透 DR/魔法，得 drEff=%d afterDR=%d", r2.DREff, r2.AfterDR)
	}
}

// DR/神兵：攻击带【神兵】特性时穿透。
func TestDRDivinePiercedByDivine(t *testing.T) {
	p := NewPreset("t")
	p.DRValue = 5
	p.DRType = "divine"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, IsPhysical: true, Divine: true, DamageLimit: 10})
	if r.DREff != 0 || r.AfterDR != 10 {
		t.Fatalf("神兵攻击应穿透 DR/神兵，得 drEff=%d afterDR=%d", r.DREff, r.AfterDR)
	}
}

// DR/穿刺：物理子类型匹配时穿透，不匹配则 DR 生效。
func TestDRPiercingByPhysSubtype(t *testing.T) {
	p := NewPreset("t")
	p.DRValue = 5
	p.DRType = "piercing"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, IsPhysical: true, PhysSubtype: "piercing", DamageLimit: 10})
	if r.DREff != 0 || r.AfterDR != 10 {
		t.Fatalf("穿刺攻击应穿透 DR/穿刺，得 drEff=%d afterDR=%d", r.DREff, r.AfterDR)
	}
	r2 := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, IsPhysical: true, PhysSubtype: "slashing", DamageLimit: 10})
	if r2.DREff != 5 || r2.AfterDR != 5 {
		t.Fatalf("挥砍攻击不应穿透 DR/穿刺，得 drEff=%d afterDR=%d", r2.DREff, r2.AfterDR)
	}
}

// 能量抗力按子类型匹配：火焰抗力对火焰生效、对雷电不生效；全能量抗力对任意生效。
func TestERSubtypeMatching(t *testing.T) {
	p := NewPreset("t")
	p.ERValue = 4
	p.ERType = "fire"
	rFire := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "energy", EnergyType: "fire", DamageLimit: 10})
	if rFire.EREff != 4 || rFire.AfterDR != 6 {
		t.Fatalf("火焰抗力应对火焰能量生效4，得 erEff=%d afterDR=%d", rFire.EREff, rFire.AfterDR)
	}
	rLightning := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "energy", EnergyType: "lightning", DamageLimit: 10})
	if rLightning.EREff != 0 || rLightning.AfterDR != 10 {
		t.Fatalf("火焰抗力不应对雷电能量生效，得 erEff=%d afterDR=%d", rLightning.EREff, rLightning.AfterDR)
	}
	p.ERType = "all"
	rAll := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "energy", EnergyType: "lightning", DamageLimit: 10})
	if rAll.EREff != 4 || rAll.AfterDR != 6 {
		t.Fatalf("全能量抗力应对雷电生效4，得 erEff=%d afterDR=%d", rAll.EREff, rAll.AfterDR)
	}
}

func TestEnergyDamageUsesERAndSkipsPhysicalAbsorb(t *testing.T) {
	p := NewPreset("t")
	p.DRValue = 5
	p.ERValue = 2
	p.DamageAbsorb = 3
	p.DamageAbsorbType = "physical"
	r := p.ResolveAttack(AttackInput{AttackDP: 100, BonusSuccess: 30, DamageKind: "energy", DamageLimit: 10})
	if r.Miss {
		t.Fatalf("不应未命中: %s", r.MissReason)
	}
	if r.AfterDR != 8 {
		t.Fatalf("能量伤害应吃ER2，减免后期望8，得%d", r.AfterDR)
	}
	if r.AbsorbApplied != 0 || r.FinalDamage != 8 {
		t.Fatalf("能量伤害不应触发物理吸收，得 absorb=%d final=%d", r.AbsorbApplied, r.FinalDamage)
	}
}

// 完美防御：不进三池；生效时直接从攻击 DP 扣除；不被击破
func TestPerfectDefense(t *testing.T) {
	p := NewPreset("t")
	p.Base = 5
	p.PerfectDefense = 3
	p.PerfectDefenseActive = true
	r := p.ResolveAttack(AttackInput{AttackDP: 20, Speed: 100, IsPhysical: true})
	// 高速击破后高速池为0，完美仍生效
	if r.Pools.SpeedPoolAfter != 0 {
		t.Fatalf("高速池应被击破到0，得%d", r.Pools.SpeedPoolAfter)
	}
	if r.PerfectDefense != 3 {
		t.Fatalf("完美防御期望3，得%d", r.PerfectDefense)
	}
	// actualDP = 20 - 0 - 3 = 17
	if r.ActualDP != 17 {
		t.Fatalf("实际DP期望17，得%d (effDef=%d)", r.ActualDP, r.EffectiveDefense)
	}
	// 未生效时不扣
	p.PerfectDefenseActive = false
	r2 := p.ResolveAttack(AttackInput{AttackDP: 20, IsPhysical: true})
	if r2.ActualDP != 15 { // 20-5
		t.Fatalf("未生效时期望DP15，得%d", r2.ActualDP)
	}
}
