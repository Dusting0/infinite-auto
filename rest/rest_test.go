package rest

import (
	"testing"

	"infinite-calc/hp"
)

// 场景：短休只清除冲击，不影响严重与恶性
// 前置：4B + 3L + 1A
// 操作：ApplyShortRest
// 预期：B=0，L=3，A=1
func TestApplyShortRest(t *testing.T) {
	h := hp.NewHPState(20)
	h.ApplyDamage(4, "冲击")
	h.ApplyDamage(3, "严重")
	h.ApplyDamage(1, "恶性")
	ApplyShortRest(h)
	if h.B != 0 || h.L != 3 || h.A != 1 {
		t.Fatalf("期望 B0 L3 A1，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：长休选项探测 — 同时有严重与恶性时两个选项都可用
// 前置：2L + 2A
// 预期：CanClearL=true, CanConvertA=true
func TestAvailableOptionsBoth(t *testing.T) {
	h := hp.NewHPState(10)
	h.ApplyDamage(2, "严重")
	h.ApplyDamage(2, "恶性")
	opt := AvailableLongRestOptions(h)
	if !opt.CanClearL || !opt.CanConvertA {
		t.Fatalf("期望两个选项都可用，得到 %+v", opt)
	}
}

// 场景：仅有严重时，只能清严重
// 前置：3L
// 预期：CanClearL=true, CanConvertA=false
func TestAvailableOptionsOnlyL(t *testing.T) {
	h := hp.NewHPState(10)
	h.ApplyDamage(3, "严重")
	opt := AvailableLongRestOptions(h)
	if !opt.CanClearL || opt.CanConvertA {
		t.Fatalf("期望仅 CanClearL，得到 %+v", opt)
	}
}

// 场景：仅有恶性时，只能转化恶性
// 前置：2A
// 预期：CanClearL=false, CanConvertA=true
func TestAvailableOptionsOnlyA(t *testing.T) {
	h := hp.NewHPState(10)
	h.ApplyDamage(2, "恶性")
	opt := AvailableLongRestOptions(h)
	if opt.CanClearL || !opt.CanConvertA {
		t.Fatalf("期望仅 CanConvertA，得到 %+v", opt)
	}
}

// 场景：长休选择清严重
// 前置：2B + 4L + 3A
// 操作：ApplyLongRest(ModeClearL)
// 预期：B=2，L=0，A=3（恶性不变）
func TestApplyLongRestClearL(t *testing.T) {
	h := hp.NewHPState(20)
	h.ApplyDamage(2, "冲击")
	h.ApplyDamage(4, "严重")
	h.ApplyDamage(3, "恶性")
	if err := ApplyLongRest(h, ModeClearL); err != nil {
		t.Fatal(err)
	}
	if h.B != 2 || h.L != 0 || h.A != 3 {
		t.Fatalf("期望 B2 L0 A3，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：长休选择恶性转严重
// 前置：1B + 2L + 3A
// 操作：ApplyLongRest(ModeConvertA)
// 预期：B=1，L=5，A=0
func TestApplyLongRestConvertA(t *testing.T) {
	h := hp.NewHPState(20)
	h.ApplyDamage(1, "冲击")
	h.ApplyDamage(2, "严重")
	h.ApplyDamage(3, "恶性")
	if err := ApplyLongRest(h, ModeConvertA); err != nil {
		t.Fatal(err)
	}
	if h.B != 1 || h.L != 5 || h.A != 0 {
		t.Fatalf("期望 B1 L5 A0，实际 B=%d L=%d A=%d", h.B, h.L, h.A)
	}
}

// 场景：长休非法 mode
// 参数：mode = "X"
// 预期：返回 error
func TestApplyLongRestInvalid(t *testing.T) {
	h := hp.NewHPState(10)
	if err := ApplyLongRest(h, "X"); err == nil {
		t.Fatal("非法 mode 应返回错误")
	}
}
