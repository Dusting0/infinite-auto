package hp

import (
	"fmt"
	"strings"
	"sync"
)

// HPState 记录当前生命状态（无限流TRPG · 溢出升级规则）
// 完好 / 冲击(B) / 严重(L) / 恶性(A)
type HPState struct {
	mu  sync.Mutex
	Max int      `json:"max"`
	B   int      `json:"b"` // 冲击伤害
	L   int      `json:"l"` // 严重伤害
	A   int      `json:"a"` // 恶性伤害
	Log []string `json:"log"`
}

func NewHPState(max int) *HPState {
	return &HPState{Max: max, Log: []string{}}
}

func (h *HPState) Intact() int {
	// 完好 = 上限 - (冲击+严重+恶性)；可为负，用于提示用户输入超过上限
	return h.Max - (h.B + h.L + h.A)
}

func (h *HPState) TotalDamage() int {
	return h.B + h.L + h.A
}

func (h *HPState) Status() string {
	if h.Max <= 0 {
		return "未设置生命值"
	}
	if h.Intact() < 0 {
		return "⚠️ 伤害超过上限"
	}
	if h.A >= h.Max && h.B == 0 && h.L == 0 {
		return "💀 死亡"
	}
	if h.Intact() == 0 {
		return "😵 昏迷（无完好生命值）"
	}
	return "✅ 正常"
}

func (h *HPState) Snapshot() map[string]interface{} {
	h.mu.Lock()
	defer h.mu.Unlock()
	logCopy := make([]string, len(h.Log))
	copy(logCopy, h.Log)
	return map[string]interface{}{
		"max":    h.Max,
		"intact": h.Intact(),
		"b":      h.B,
		"l":      h.L,
		"a":      h.A,
		"total":  h.TotalDamage(),
		"status": h.Status(),
		"log":    logCopy,
	}
}

// upgradeOne 升级 1 点最低伤势：B→L 或 L→A
func (h *HPState) upgradeOne() bool {
	if h.B > 0 {
		h.B--
		h.L++
		return true
	}
	if h.L > 0 {
		h.L--
		h.A++
		return true
	}
	return false
}

// convert 兜底：保证总数不超过上限
func (h *HPState) convert() {
	for h.B+h.L+h.A > h.Max {
		if !h.upgradeOne() {
			if h.A > h.Max {
				h.A = h.Max
			}
			break
		}
	}
}

// ApplyDamage 应用伤害（溢出升级版）
// 1. 先填满完好生命值
// 2. 溢出部分每 2 点升级 1 点已有最低伤势（B→L 优先，然后 L→A）
func (h *HPState) ApplyDamage(amount int, dmgType string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Max <= 0 {
		return fmt.Errorf("请先设置生命值上限")
	}
	if amount <= 0 {
		return fmt.Errorf("伤害数量必须大于0")
	}

	before := fmt.Sprintf("完好%d B%d L%d A%d", h.Intact(), h.B, h.L, h.A)

	var typeName string
	switch strings.ToUpper(strings.TrimSpace(dmgType)) {
	case "B", "冲击", "冲击伤害":
		typeName = "冲击"
	case "L", "严重", "严重伤害":
		typeName = "严重"
	case "A", "恶性", "恶性伤害":
		typeName = "恶性"
	default:
		return fmt.Errorf("未知伤害类型")
	}

	// 1. 先用伤害填满完好生命值
	intact := h.Intact()
	fill := amount
	if fill > intact {
		fill = intact
	}
	if fill > 0 {
		switch typeName {
		case "冲击":
			h.B += fill
		case "严重":
			h.L += fill
		case "恶性":
			h.A += fill
		}
	}

	// 2. 溢出伤害用于升级（向上取整：每 2 点溢出 → 升级 1 次，奇数也算 1 次）
	// 例：2严重 + 收到1点伤害 → 1严重 1恶性
	excess := amount - fill
	upgrades := (excess + 1) / 2 // 向上取整
	for i := 0; i < upgrades; i++ {
		if !h.upgradeOne() {
			break
		}
	}

	h.convert()

	after := fmt.Sprintf("完好%d B%d L%d A%d", h.Intact(), h.B, h.L, h.A)
	h.Log = append(h.Log, fmt.Sprintf("受到 %d 点%s伤害 | %s → %s", amount, typeName, before, after))
	if len(h.Log) > 30 {
		h.Log = h.Log[len(h.Log)-30:]
	}
	return nil
}

// ShortRest 短休（规则：≥1小时）
// 效果：将全部冲击伤害(B)恢复为完好
func (h *HPState) ShortRest() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.B > 0 {
		h.Log = append(h.Log, fmt.Sprintf("短休：清除全部 %d 点冲击伤害", h.B))
		h.B = 0
	} else {
		h.Log = append(h.Log, "短休：无冲击伤害需要恢复")
	}
}

// HealB 保留旧名，内部调用 ShortRest，兼容已有调用
func (h *HPState) HealB() {
	h.ShortRest()
}

// LongRest 长休（规则：≥8小时）
// mode:
//   "L" — 将全部严重伤害(L)恢复为完好
//   "A" — 将全部恶性伤害(A)转化为严重伤害(L)
// 规则原文：「可以将全部状态为L的生命值恢复为完好，或者将全部状态为A的生命值转化为L」
func (h *HPState) LongRest(mode string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	switch strings.ToUpper(strings.TrimSpace(mode)) {
	case "L", "严重", "CLEAR_L":
		if h.L > 0 {
			h.Log = append(h.Log, fmt.Sprintf("长休：全部 %d 点严重伤害恢复为完好", h.L))
			h.L = 0
		} else {
			h.Log = append(h.Log, "长休：无严重伤害需要恢复")
		}
	case "A", "恶性", "CONVERT_A":
		if h.A > 0 {
			h.Log = append(h.Log, fmt.Sprintf("长休：全部 %d 点恶性伤害转化为严重伤害", h.A))
			h.L += h.A
			h.A = 0
		} else {
			h.Log = append(h.Log, "长休：无恶性伤害需要转化")
		}
	default:
		return fmt.Errorf("长休模式无效，请使用 L（清严重）或 A（恶性转严重）")
	}
	return nil
}

func (h *HPState) HealType(dmgType string, amount int) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if amount <= 0 {
		return fmt.Errorf("恢复数量必须大于0")
	}
	var name string
	switch strings.ToUpper(strings.TrimSpace(dmgType)) {
	case "B", "冲击":
		if amount > h.B {
			amount = h.B
		}
		h.B -= amount
		name = "冲击"
	case "L", "严重":
		if amount > h.L {
			amount = h.L
		}
		h.L -= amount
		name = "严重"
	case "A", "恶性":
		if amount > h.A {
			amount = h.A
		}
		h.A -= amount
		name = "恶性"
	default:
		return fmt.Errorf("未知类型")
	}
	h.Log = append(h.Log, fmt.Sprintf("恢复 %d 点%s伤害", amount, name))
	return nil
}

func (h *HPState) SetMax(v int) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if v <= 0 {
		return fmt.Errorf("生命值上限必须为正整数")
	}
	old := h.Max
	h.Max = v
	// 不自动转化；若伤害超过新上限，完好为负以提示
	h.Log = append(h.Log, fmt.Sprintf("生命值上限 %d → %d", old, v))
	return nil
}

func (h *HPState) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.B, h.L, h.A = 0, 0, 0
	h.Log = append(h.Log, "重置为满血")
}

// SetValues 自定义设定生命值各分项（保留：上限显式给定）
// 约束：max > 0，b/l/a >= 0，且 b+l+a <= max
// 完好 = max - (b+l+a)
func (h *HPState) SetValues(max, b, l, a int) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if max <= 0 {
		return fmt.Errorf("生命值上限必须为正整数")
	}
	if b < 0 || l < 0 || a < 0 {
		return fmt.Errorf("各类型伤害不能为负数")
	}
	if b+l+a > max {
		return fmt.Errorf("冲击+严重+恶性 之和不能超过生命值上限（%d > %d）", b+l+a, max)
	}

	before := fmt.Sprintf("上限%d 完好%d B%d L%d A%d", h.Max, h.Max-(h.B+h.L+h.A), h.B, h.L, h.A)
	h.Max, h.B, h.L, h.A = max, b, l, a
	after := fmt.Sprintf("上限%d 完好%d B%d L%d A%d", h.Max, h.Intact(), h.B, h.L, h.A)
	h.Log = append(h.Log, fmt.Sprintf("自定义血量 | %s → %s", before, after))
	if len(h.Log) > 30 {
		h.Log = h.Log[len(h.Log)-30:]
	}
	return nil
}

// SetParts 按各分项设定血量，上限自动计算为 完好+冲击+严重+恶性
// 不限制用户输入（只要各项 >= 0）。完好由用户直接填入。
func (h *HPState) SetParts(intact, b, l, a int) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if intact < 0 || b < 0 || l < 0 || a < 0 {
		return fmt.Errorf("各类型生命值不能为负数")
	}
	max := intact + b + l + a
	if max <= 0 {
		return fmt.Errorf("生命值总和必须大于 0")
	}

	before := fmt.Sprintf("上限%d 完好%d B%d L%d A%d", h.Max, h.Max-(h.B+h.L+h.A), h.B, h.L, h.A)
	h.Max, h.B, h.L, h.A = max, b, l, a
	after := fmt.Sprintf("上限%d 完好%d B%d L%d A%d", h.Max, h.Intact(), h.B, h.L, h.A)
	h.Log = append(h.Log, fmt.Sprintf("自定义血量 | %s → %s", before, after))
	if len(h.Log) > 30 {
		h.Log = h.Log[len(h.Log)-30:]
	}
	return nil
}

// SetDamages 在保持生命值上限不变的前提下设定冲击/严重/恶性。
// 完好 = 上限 - (B+L+A)，由系统自动计算，不可直接设定。
// 若 B+L+A > 上限，完好为负值，用于提示输入有误（不自动转化）。
func (h *HPState) SetDamages(b, l, a int) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Max <= 0 {
		return fmt.Errorf("请先设置生命值上限")
	}
	if b < 0 || l < 0 || a < 0 {
		return fmt.Errorf("各类型伤害不能为负数")
	}

	before := fmt.Sprintf("上限%d 完好%d B%d L%d A%d", h.Max, h.Max-(h.B+h.L+h.A), h.B, h.L, h.A)
	h.B, h.L, h.A = b, l, a
	after := fmt.Sprintf("上限%d 完好%d B%d L%d A%d", h.Max, h.Intact(), h.B, h.L, h.A)
	h.Log = append(h.Log, fmt.Sprintf("设定伤害 | %s → %s", before, after))
	if len(h.Log) > 30 {
		h.Log = h.Log[len(h.Log)-30:]
	}
	return nil
}
