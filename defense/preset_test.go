package defense

import "testing"

// 场景：基础成分合计（无特殊状态）
// 参数：基础5、闪避2、天生3、盔甲近战4、力场1
// 预期：近战 = 5+2+3+4+1 = 15；远程不含近战盔甲
func TestBasicTotals(t *testing.T) {
	p := NewPreset("测试")
	p.Base = 5
	p.Dodge = 2
	p.Natural = 3
	p.ArmorMelee = 4
	p.ArmorRanged = 1
	p.Force = 1
	tot := p.Compute()
	if tot.Melee != 15 {
		t.Fatalf("近战期望15，得到%d", tot.Melee)
	}
	// 远程: 5+2+3+1+1 = 12
	if tot.Ranged != 12 {
		t.Fatalf("远程期望12，得到%d", tot.Ranged)
	}
}

// 场景：措手不及清零基础与闪避与格挡
// 前置：基础5、闪避3、格挡4且开启格挡
// 操作：FlatFooted=true
// 预期：上述三项生效值均为0
func TestFlatFooted(t *testing.T) {
	p := NewPreset("")
	p.Base = 5
	p.Dodge = 3
	p.Block = 4
	p.Blocking = true
	p.FlatFooted = true
	tot := p.Compute()
	if tot.SpeedPoolMelee != 0 {
		t.Fatalf("措手不及时高速池应为0，得到%d", tot.SpeedPoolMelee)
	}
}

// 场景：全力防御使基础翻倍
// 参数：基础5，FullDefense=true
// 预期：生效基础10
func TestFullDefense(t *testing.T) {
	p := NewPreset("")
	p.Base = 5
	p.FullDefense = true
	snap := p.Snapshot()
	eff := snap["effective"].(map[string]int)
	if eff["base"] != 10 {
		t.Fatalf("全力防御基础期望10，得到%d", eff["base"])
	}
}

// 场景：格挡未开启时格挡防御与盾牌不计入
// 参数：Block=5, ShieldMelee=3, Blocking=false
// 预期：近战不含这8点
func TestBlockRequiresCheckbox(t *testing.T) {
	p := NewPreset("")
	p.Base = 2
	p.Block = 5
	p.ShieldMelee = 3
	p.Blocking = false
	tot := p.Compute()
	if tot.Melee != 2 {
		t.Fatalf("未格挡时期望近战2，得到%d", tot.Melee)
	}
	p.Blocking = true
	tot = p.Compute()
	if tot.Melee != 2+5+3 {
		t.Fatalf("格挡时期望近战10，得到%d", tot.Melee)
	}
}

// 场景：掩蔽开启计入 CoverBonus
// 参数：CoverBonus=4, Cover=true
// 预期：OtherPool 含4
func TestCover(t *testing.T) {
	p := NewPreset("")
	p.Cover = true
	p.CoverBonus = 4
	tot := p.Compute()
	if tot.OtherPool < 4 || tot.Melee < 4 {
		t.Fatalf("掩蔽未计入: %+v", tot)
	}
}

// 场景：抵高速/破甲/破魔字段可保存
func TestResistsStored(t *testing.T) {
	p := NewPreset("")
	p.ResistSpeed = 2
	p.ResistAP = 3
	p.ResistMagic = 1
	snap := p.Snapshot()
	if snap["resistSpeed"] != 2 || snap["resistAP"] != 3 || snap["resistMagic"] != 1 {
		t.Fatalf("抵消值未正确保存: %+v", snap)
	}
}

// 场景：自定义成分（仅近战、破甲类）
func TestCustomMeleeArmor(t *testing.T) {
	p := NewPreset("")
	p.Customs = []CustomComp{{
		Name: "龙鳞外挂", Value: 5, Category: CatArmor, Scope: "melee",
	}}
	tot := p.Compute()
	if tot.Melee != 5 || tot.Ranged != 0 {
		t.Fatalf("自定义近战盔甲类期望近战5远程0，得到 M%d R%d", tot.Melee, tot.Ranged)
	}
	if tot.ArmorPoolMelee != 5 {
		t.Fatalf("破甲池期望5，得到%d", tot.ArmorPoolMelee)
	}
}

// 场景：规则书影牙巨龙式合计
// 基础8+闪避1+洞察1+天生6+偏斜2 = 18
func TestDragonExample(t *testing.T) {
	p := NewPreset("影牙")
	p.Base = 8
	p.Dodge = 1
	p.Insight = 1
	p.Natural = 6
	p.Deflection = 2
	tot := p.Compute()
	if tot.Melee != 18 {
		t.Fatalf("期望18，得到%d", tot.Melee)
	}
}
