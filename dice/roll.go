// Package dice 无限规则掷骰检定。
// 规则要点：
//   - DP 枚 D10；8/9/10 为成功
//   - 加骰：默认 10 加骰（掷出 ≥ 加骰阈值 再投，可连锁），最低 8
//   - 附加成功：仅当自然成功数 > 0 时计入；可为负；最终成功数不低于 0
//   - 机运骰：DP≤0 时投 1 枚，仅 10 算成功；首骰自然 1 且全程无成功 → 大失败
//   - 机运骰的加骰骰子按普通规则（8/9/10）计成功
package dice

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

// 可替换的掷骰函数，便于单测注入固定序列
var rollD10 = secureD10

// Result 一次检定的完整结果
type Result struct {
	DP              int     `json:"dp"`
	ExplodeOn       int     `json:"explodeOn"`
	BonusSuccess    int     `json:"bonusSuccess"`
	ChanceDie       bool    `json:"chanceDie"`
	Dice            []int   `json:"dice"`       // 首轮骰面（主池 / 机运首骰）
	Explosions      []int   `json:"explosions"` // 所有加骰扁平列表（兼容）
	Batches         [][]int `json:"batches"`    // 按加骰波次分组，如 [[1,2,10],[3,10],[4]]
	NaturalSuccess  int     `json:"naturalSuccess"`
	FinalSuccess    int     `json:"finalSuccess"`
	BonusApplied    bool    `json:"bonusApplied"`
	CriticalFailure bool    `json:"criticalFailure"`
	Summary         string  `json:"summary"`
	DiceLog         string  `json:"diceLog"` // 如 [1,2,10]+[3,10]+[4]
}

// Roll 执行一次检定。
// dp: 骰池；explodeOn: 加骰阈值（8/9/10，非法则回退 10）；bonus: 附加成功（可负）。
func Roll(dp, explodeOn, bonus int) Result {
	if explodeOn < 8 || explodeOn > 10 {
		explodeOn = 10
	}
	r := Result{
		DP:           dp,
		ExplodeOn:    explodeOn,
		BonusSuccess: bonus,
	}

	if dp <= 0 {
		r.ChanceDie = true
		rollChanceDie(&r)
	} else {
		rollNormal(&r, dp)
	}

	// 附加成功仅在自然成功 > 0 时计入；最终 ≥ 0
	if r.NaturalSuccess > 0 {
		r.BonusApplied = true
		r.FinalSuccess = r.NaturalSuccess + bonus
		if r.FinalSuccess < 0 {
			r.FinalSuccess = 0
		}
	} else {
		r.BonusApplied = false
		r.FinalSuccess = 0
	}

	r.Summary = buildSummary(r)
	r.DiceLog = formatBatches(r.Batches)
	return r
}

func formatBatches(batches [][]int) string {
	if len(batches) == 0 {
		return ""
	}
	parts := make([]string, 0, len(batches))
	for _, b := range batches {
		s := "["
		for i, v := range b {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprintf("%d", v)
		}
		s += "]"
		parts = append(parts, s)
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += "+" + parts[i]
	}
	return out
}

func rollNormal(r *Result, dp int) {
	wave := make([]int, dp)
	for i := 0; i < dp; i++ {
		wave[i] = rollD10()
	}
	appendWave(r, wave, false)
	for {
		n := 0
		for _, v := range wave {
			if v >= r.ExplodeOn {
				n++
			}
		}
		if n == 0 {
			break
		}
		next := make([]int, n)
		for i := 0; i < n; i++ {
			next[i] = rollD10()
		}
		appendWave(r, next, true)
		wave = next
	}
}

func rollChanceDie(r *Result) {
	v := rollD10()
	wave := []int{v}
	// 机运主骰仅 10 算成功；加骰波次按普通 ≥8
	if v == 10 {
		r.NaturalSuccess++
	}
	r.Dice = []int{v}
	r.Batches = [][]int{{v}}
	for {
		n := 0
		for _, x := range wave {
			if x >= r.ExplodeOn {
				n++
			}
		}
		if n == 0 {
			break
		}
		next := make([]int, n)
		for i := 0; i < n; i++ {
			next[i] = rollD10()
			r.Explosions = append(r.Explosions, next[i])
			if next[i] >= 8 {
				r.NaturalSuccess++
			}
		}
		r.Batches = append(r.Batches, next)
		wave = next
	}
	if r.Dice[0] == 1 && r.NaturalSuccess == 0 {
		r.CriticalFailure = true
	}
}

// appendWave 记录一轮骰面；isExplode 表示是否为加骰波次
func appendWave(r *Result, wave []int, isExplode bool) {
	cp := append([]int(nil), wave...)
	r.Batches = append(r.Batches, cp)
	if !isExplode {
		r.Dice = cp
		for _, v := range wave {
			if v >= 8 {
				r.NaturalSuccess++
			}
		}
		return
	}
	for _, v := range wave {
		r.Explosions = append(r.Explosions, v)
		if v >= 8 {
			r.NaturalSuccess++
		}
	}
}

func secureD10() int {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return int(binary.LittleEndian.Uint64(b[:])%10) + 1
}

func buildSummary(r Result) string {
	if r.ChanceDie {
		if r.CriticalFailure {
			return "机运骰 · 大失败（首骰1，无成功）"
		}
		if r.BonusApplied {
			return fmt.Sprintf("机运骰 · 自然%d + 附加%d → 最终%d", r.NaturalSuccess, r.BonusSuccess, r.FinalSuccess)
		}
		return fmt.Sprintf("机运骰 · 自然%d → 最终%d", r.NaturalSuccess, r.FinalSuccess)
	}
	if r.BonusApplied {
		return fmt.Sprintf("DP%d · 自然%d + 附加%d → 最终%d", r.DP, r.NaturalSuccess, r.BonusSuccess, r.FinalSuccess)
	}
	return fmt.Sprintf("DP%d · 自然%d → 最终%d（附加未计入）", r.DP, r.NaturalSuccess, r.FinalSuccess)
}
