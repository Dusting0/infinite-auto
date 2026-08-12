package hp

import (
	"testing"
)

// ---------------------------------------------------------------------------
// 场景：初始化与基础查询
// ---------------------------------------------------------------------------

// 场景：创建指定上限的生命状态
// 参数：max = 20
// 预期：完好=20，B/L/A 均为 0，状态为正常
func TestNewHPState(t *testing.T) {
	h := NewHPState(20)
	if h.Max != 20 || h.Intact() != 20 || h.B != 0 || h.L != 0 || h.A != 0 {
		t.Fatalf("初始化失败: Max=%d Intact=%d B=%d L=%d A=%d", h.Max, h.Intact(), h.B, h.L, h.A)
	}
	if h.Status() != "✅ 正常" {
		t.Fatalf("期望状态正常，得到 %s", h.Status())
	}
}

// 场景：未设置上限时的状态文案
// 参数：max = 0
// 预期：Status 返回「未设置生命值」
func TestStatusUnset(t *testing.T) {
	h := NewHPState(0)
	if h.Status() != "未设置生命值" {
		t.Fatalf("期望「未设置生命值」，得到 %s", h.Status())
	}
}

// ---------------------------------------------------------------------------
// 场景：正常填入完好（未溢出）
// ---------------------------------------------------------------------------

// 场景：有完好生命值时受到冲击伤害
// 前置：上限 10，当前满血
// 参数：amount=3, type=冲击
// 预期：完好=7，B=3，L=0，A=0，状态正常
func TestApplyDamageFillIntact(t *testing.T) {
	h := NewHPState(10)
	if err := h.ApplyDamage(3, "冲击"); err != nil {
		t.Fatal(err)
	}
	if h.Intact() != 7 || h.B != 3 || h.L != 0 || h.A != 0 {
		t.Fatalf("期望 完好7 B3，实际 Intact=%d B=%d L=%d A=%d", h.Intact(), h.B, h.L, h.A)
	}
}

// 场景：多种伤害依次填入，未超过上限
// 前置：上限 20
// 参数：3冲击 → 4严重 → 5恶性
// 预期：完好=8，B=3，L=4，A=5（与规则书长示例前半段一致）
func TestApplyDamageMixedNoOverflow(t *testing.T) {
	h := NewHPState(20)
	h.ApplyDamage(3, "冲击")
	h.ApplyDamage(4, "严重")
	h.ApplyDamage(5, "恶性")
	if h.Intact() != 8 || h.B != 3 || h.L != 4 || h.A != 5 {
		t.Fatalf("期望 完好8 B3 L4 A5，实际 Intact=%d B=%d L=%d A=%d", h.Intact(), h.B, h.L, h.A)
	}
}

// ---------------------------------------------------------------------------
// 场景：溢出升级（核心规则）
// 规则：先填完好；溢出部分按向上取整，每 2 点溢出触发 1 次升级（优先 B→L，再 L→A）
// ---------------------------------------------------------------------------

// 场景：已满严重伤，再受到 2 点伤害 → 升级 1 次
// 前置：上限 10，已有 10 严重（无完好）
// 参数：amount=2, type=严重
// 预期：完好=0，B=0，L=9，A=1
func TestOverflowUpgrade_2onFullL(t *testing.T) {
	h := NewHPState(10)
	h.ApplyDamage(10, "严重")
	h.ApplyDamage(2, "严重")
	if h.Intact() != 0 || h.B != 0 || h.L != 9 || h.A != 1 {
		t.Fatalf("期望 0完好 0B 9L 1A，实际 Intact=%d B=%d L=%d A=%d", h.Intact(), h.B, h.L, h.A)
	}
}

// 场景：已满严重伤，再受到 1 点恶性伤害 → 向上取整触发 1 次升级
// 前置：上限 2，已有 2 严重
// 参数：amount=1, type=恶性
// 预期：完好=0，L=1，A=1
func TestOverflowUpgrade_1onFullL_Ceil(t *testing.T) {
	h := NewHPState(2)
	h.ApplyDamage(2, "严重")
	if err := h.ApplyDamage(1, "恶性"); err != nil {
		t.Fatal(err)
	}
	if h.Intact() != 0 || h.B != 0 || h.L != 1 || h.A != 1 {
		t.Fatalf("期望 1L 1A，实际 Intact=%d B=%d L=%d A=%d", h.Intact(), h.B, h.L, h.A)
	}
}

// 场景：溢出 3 点 → 向上取整为 2 次升级
// 前置：上限 5，已有 5 严重
// 参数：amount=3
// 预期：L=3，A=2
func TestOverflowUpgrade_3_CeilTo2(t *testing.T) {
	h := NewHPState(5)
	h.ApplyDamage(5, "严重")
	h.ApplyDamage(3, "冲击")
	if h.L != 3 || h.A != 2 {
		t.Fatalf("期望 3L 2A，实际 L=%d A=%d", h.L, h.A)
	}
}

// 场景：有冲击时优先升级冲击
// 前置：上限 10，4冲击 + 6严重（已满）
// 参数：amount=2
// 预期：优先 B→L，得到 3B 7L
func TestOverflowUpgrade_PreferB(t *testing.T) {
	h := NewHPState(10)
	h.ApplyDamage(4, "冲击")
	h.ApplyDamage(6, "严重")
	h.ApplyDamage(2, "冲击")
	if h.B != 3 || h.L != 7 || h.A != 0 {
		t.Fatalf("期望 3B 7L，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：先填完好，再对溢出部分升级
// 前置：上限 10，已有 7 严重（剩 3 完好）
// 参数：amount=5 严重
// 预期：先填 3 完好 → 临时 10L，剩余 2 点溢出 → 升级 1 次 → 9L 1A
func TestFillThenOverflowUpgrade(t *testing.T) {
	h := NewHPState(10)
	h.ApplyDamage(7, "严重")
	h.ApplyDamage(5, "严重")
	if h.Intact() != 0 || h.L != 9 || h.A != 1 {
		t.Fatalf("期望 9L 1A，实际 Intact=%d L=%d A=%d", h.Intact(), h.L, h.A)
	}
}

// ---------------------------------------------------------------------------
// 场景：死亡与昏迷判定
// ---------------------------------------------------------------------------

// 场景：全部转为恶性 → 死亡
// 前置：上限 3
// 操作：连续造成足够恶性/升级使 A 填满
// 预期：Status 为「死亡」
func TestDeathWhenAllMalignant(t *testing.T) {
	h := NewHPState(3)
	h.ApplyDamage(3, "恶性")
	if h.Status() != "💀 死亡" {
		t.Fatalf("期望死亡，得到 %s (A=%d)", h.Status(), h.A)
	}
}

// 场景：无完好但未全恶性 → 昏迷
// 前置：上限 5，受到 5 严重
// 预期：Status 为「昏迷」
func TestUnconsciousWhenNoIntact(t *testing.T) {
	h := NewHPState(5)
	h.ApplyDamage(5, "严重")
	if h.Status() != "😵 昏迷（无完好生命值）" {
		t.Fatalf("期望昏迷，得到 %s", h.Status())
	}
}

// ---------------------------------------------------------------------------
// 场景：恢复与重置
// ---------------------------------------------------------------------------

// 场景：短休清除全部冲击
// 前置：有 4 冲击 + 2 严重
// 操作：HealB()
// 预期：B=0，L 不变
func TestHealB(t *testing.T) {
	h := NewHPState(10)
	h.ApplyDamage(4, "冲击")
	h.ApplyDamage(2, "严重")
	h.HealB()
	if h.B != 0 || h.L != 2 {
		t.Fatalf("期望 B=0 L=2，实际 B=%d L=%d", h.B, h.L)
	}
}

// 场景：按类型恢复指定数量
// 前置：3冲击 3严重 2恶性
// 参数：恢复 2 点严重
// 预期：L=1，其余不变
func TestHealType(t *testing.T) {
	h := NewHPState(20)
	h.ApplyDamage(3, "冲击")
	h.ApplyDamage(3, "严重")
	h.ApplyDamage(2, "恶性")
	if err := h.HealType("严重", 2); err != nil {
		t.Fatal(err)
	}
	if h.B != 3 || h.L != 1 || h.A != 2 {
		t.Fatalf("期望 B3 L1 A2，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：重置为满血
// 操作：Reset()
// 预期：B=L=A=0，完好=上限
func TestReset(t *testing.T) {
	h := NewHPState(15)
	h.ApplyDamage(5, "严重")
	h.ApplyDamage(3, "恶性")
	h.Reset()
	if h.B != 0 || h.L != 0 || h.A != 0 || h.Intact() != 15 {
		t.Fatalf("重置失败: B=%d L=%d A=%d Intact=%d", h.B, h.L, h.A, h.Intact())
	}
}

// ---------------------------------------------------------------------------
// 场景：修改上限
// ---------------------------------------------------------------------------

// 场景：下调上限后伤害可超过上限，完好为负
// 前置：上限 10，已有 8 严重
// 操作：SetMax(5)
// 预期：Max=5，L=8，Intact=-3
func TestSetMaxDown(t *testing.T) {
	h := NewHPState(10)
	h.ApplyDamage(8, "严重")
	if err := h.SetMax(5); err != nil {
		t.Fatal(err)
	}
	if h.Max != 5 {
		t.Fatalf("Max 未更新: %d", h.Max)
	}
	if h.L != 8 || h.Intact() != -3 {
		t.Fatalf("期望 L=8 Intact=-3，实际 L=%d Intact=%d", h.L, h.Intact())
	}
}

// 场景：非法参数
// 参数：amount=0 或 负数；未设置上限就造成伤害
// 预期：返回 error
func TestApplyDamageInvalid(t *testing.T) {
	h := NewHPState(10)
	if err := h.ApplyDamage(0, "冲击"); err == nil {
		t.Fatal("amount=0 应返回错误")
	}
	h2 := NewHPState(0)
	if err := h2.ApplyDamage(1, "冲击"); err == nil {
		t.Fatal("未设置上限应返回错误")
	}
}

// ---------------------------------------------------------------------------
// 场景：短休 / 长休（规则书自然恢复）
// ---------------------------------------------------------------------------

// 场景：短休清除全部冲击
// 前置：上限10，4冲击 + 3严重 + 1恶性
// 操作：ShortRest()
// 预期：B=0，L=3，A=1 不变
func TestShortRest(t *testing.T) {
	h := NewHPState(10)
	h.ApplyDamage(4, "冲击")
	h.ApplyDamage(3, "严重")
	h.ApplyDamage(1, "恶性")
	h.ShortRest()
	if h.B != 0 || h.L != 3 || h.A != 1 {
		t.Fatalf("短休后期望 B0 L3 A1，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：长休模式 L — 全部严重恢复为完好
// 前置：2冲击 + 4严重 + 2恶性
// 操作：LongRest("L")
// 预期：B=2，L=0，A=2
func TestLongRestClearL(t *testing.T) {
	h := NewHPState(20)
	h.ApplyDamage(2, "冲击")
	h.ApplyDamage(4, "严重")
	h.ApplyDamage(2, "恶性")
	if err := h.LongRest("L"); err != nil {
		t.Fatal(err)
	}
	if h.B != 2 || h.L != 0 || h.A != 2 {
		t.Fatalf("期望 B2 L0 A2，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：长休模式 A — 全部恶性转化为严重
// 前置：1冲击 + 2严重 + 3恶性
// 操作：LongRest("A")
// 预期：B=1，L=5，A=0
func TestLongRestConvertA(t *testing.T) {
	h := NewHPState(20)
	h.ApplyDamage(1, "冲击")
	h.ApplyDamage(2, "严重")
	h.ApplyDamage(3, "恶性")
	if err := h.LongRest("A"); err != nil {
		t.Fatal(err)
	}
	if h.B != 1 || h.L != 5 || h.A != 0 {
		t.Fatalf("期望 B1 L5 A0，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：长休非法模式
// 参数：mode = "X"
// 预期：返回 error
func TestLongRestInvalidMode(t *testing.T) {
	h := NewHPState(10)
	if err := h.LongRest("X"); err == nil {
		t.Fatal("非法模式应返回错误")
	}
}

// ---------------------------------------------------------------------------
// 场景：自定义血量（二级面板）
// ---------------------------------------------------------------------------

// 场景：合法自定义各分项
// 参数：max=12, B=2, L=3, A=1
// 预期：上限12，完好=6，B2 L3 A1
func TestSetValuesOK(t *testing.T) {
	h := NewHPState(10)
	if err := h.SetValues(12, 2, 3, 1); err != nil {
		t.Fatal(err)
	}
	if h.Max != 12 || h.Intact() != 6 || h.B != 2 || h.L != 3 || h.A != 1 {
		t.Fatalf("期望 Max12 完好6 B2 L3 A1，实际 Max=%d Intact=%d B=%d L=%d A=%d",
			h.Max, h.Intact(), h.B, h.L, h.A)
	}
}

// 场景：分项之和超过上限
// 参数：max=5, B=2, L=2, A=2（和=6>5）
// 预期：返回 error，状态不变
func TestSetValuesExceedMax(t *testing.T) {
	h := NewHPState(10)
	h.ApplyDamage(1, "冲击")
	err := h.SetValues(5, 2, 2, 2)
	if err == nil {
		t.Fatal("之和超过上限应返回错误")
	}
	// 原状态应保持
	if h.Max != 10 || h.B != 1 {
		t.Fatalf("失败后状态被修改: Max=%d B=%d", h.Max, h.B)
	}
}

// 场景：负数分项非法
// 参数：B=-1
// 预期：返回 error
func TestSetValuesNegative(t *testing.T) {
	h := NewHPState(10)
	if err := h.SetValues(10, -1, 0, 0); err == nil {
		t.Fatal("负数应返回错误")
	}
}

// ---------------------------------------------------------------------------
// 场景：SetParts — 按分项设定，上限自动 = 完好+B+L+A
// ---------------------------------------------------------------------------

// 场景：合法设定各分项，上限自动计算
// 参数：完好=5, B=2, L=3, A=1
// 预期：Max=11, Intact=5, B2 L3 A1
func TestSetPartsOK(t *testing.T) {
	h := NewHPState(10)
	if err := h.SetParts(5, 2, 3, 1); err != nil {
		t.Fatal(err)
	}
	if h.Max != 11 || h.Intact() != 5 || h.B != 2 || h.L != 3 || h.A != 1 {
		t.Fatalf("期望 Max11 完好5 B2 L3 A1，实际 Max=%d Intact=%d B=%d L=%d A=%d",
			h.Max, h.Intact(), h.B, h.L, h.A)
	}
}

// 场景：完好为 0 也可以（已满伤）
// 参数：完好=0, B=0, L=4, A=2 → Max=6
// 预期：成功，Intact=0
func TestSetPartsZeroIntact(t *testing.T) {
	h := NewHPState(10)
	if err := h.SetParts(0, 0, 4, 2); err != nil {
		t.Fatal(err)
	}
	if h.Max != 6 || h.Intact() != 0 || h.L != 4 || h.A != 2 {
		t.Fatalf("期望 Max6 L4 A2，实际 Max=%d Intact=%d L=%d A=%d", h.Max, h.Intact(), h.L, h.A)
	}
}

// 场景：负数非法
// 参数：完好=-1
// 预期：返回 error
func TestSetPartsNegative(t *testing.T) {
	h := NewHPState(10)
	if err := h.SetParts(-1, 0, 0, 0); err == nil {
		t.Fatal("负数应返回错误")
	}
}

// ---------------------------------------------------------------------------
// 场景：SetDamages — 上限固定，完好自动计算
// ---------------------------------------------------------------------------

// 场景：设定伤害后完好 = 上限 - 伤害总和
// 参数：上限10，B=2 L=3 A=1
// 预期：完好=4
func TestSetDamagesOK(t *testing.T) {
	h := NewHPState(10)
	if err := h.SetDamages(2, 3, 1); err != nil {
		t.Fatal(err)
	}
	if h.Max != 10 || h.Intact() != 4 || h.B != 2 || h.L != 3 || h.A != 1 {
		t.Fatalf("期望 Max10 完好4 B2 L3 A1，实际 Max=%d Intact=%d B=%d L=%d A=%d",
			h.Max, h.Intact(), h.B, h.L, h.A)
	}
}

// 场景：伤害之和超过上限时完好为负（提示输入有误，不自动转化）
// 参数：上限5，设定 L=8
// 预期：L=8 保持，完好 = 5-8 = -3
func TestSetDamagesOverflowNegativeIntact(t *testing.T) {
	h := NewHPState(5)
	if err := h.SetDamages(0, 8, 0); err != nil {
		t.Fatal(err)
	}
	if h.L != 8 || h.Intact() != -3 {
		t.Fatalf("期望 L=8 Intact=-3，实际 L=%d Intact=%d", h.L, h.Intact())
	}
	if h.Status() != "⚠️ 伤害超过上限" {
		t.Fatalf("期望状态提示超过上限，得到 %s", h.Status())
	}
}
