// Package rest 统一处理短休 / 长休。
// 当前仅作用于血量模块；后续可在此扩展意志力、不良状态、属性伤害等效果。
package rest

import (
	"fmt"

	"infinite-calc/hp"
)

// ApplyShortRest 短休（规则：≥1 小时）
// 当前对血量的影响：全部冲击伤害(B) → 完好
func ApplyShortRest(h *hp.HPState) {
	if h == nil {
		return
	}
	h.ShortRest()
	// TODO: 恢复 1 点意志力、一种不良状态豁免 等
}

// 长休对血量的两种互斥选项（规则书「或者」）
const (
	ModeClearL   = "L" // 全部严重(L) → 完好
	ModeConvertA = "A" // 全部恶性(A) → 严重(L)
)

// Options 描述当前角色长休时可选的血量效果
type Options struct {
	CanClearL   bool // 是否有严重伤害可清除
	CanConvertA bool // 是否有恶性伤害可转化
}

// AvailableLongRestOptions 根据当前血量判断长休可选效果。
// 仅当 CanClearL 与 CanConvertA 均为 true 时，前端才需要弹窗让用户二选一。
func AvailableLongRestOptions(h *hp.HPState) Options {
	if h == nil {
		return Options{}
	}
	snap := h.Snapshot()
	// Snapshot 返回 map，数值为 int
	l, _ := toInt(snap["l"])
	a, _ := toInt(snap["a"])
	return Options{
		CanClearL:   l > 0,
		CanConvertA: a > 0,
	}
}

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case float64:
		return int(n), true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

// ApplyLongRest 长休（规则：≥8 小时）
// mode 必须为 ModeClearL 或 ModeConvertA，二者互斥（规则书「或者」）。
func ApplyLongRest(h *hp.HPState, mode string) error {
	if h == nil {
		return fmt.Errorf("血量状态为空")
	}
	switch mode {
	case ModeClearL, "严重", "CLEAR_L":
		return h.LongRest("L")
	case ModeConvertA, "恶性", "CONVERT_A":
		return h.LongRest("A")
	default:
		return fmt.Errorf("长休必须选择：L（清全部严重）或 A（恶性转严重）")
	}
	// TODO: 意志力恢复、属性伤害、不良状态豁免 等
}
