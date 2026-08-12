package dice

import "testing"

// 固定序列掷骰：按顺序返回 values，用尽后循环最后一个
func withSeq(values []int, fn func()) {
	i := 0
	old := rollD10
	rollD10 = func() int {
		if i < len(values) {
			v := values[i]
			i++
			return v
		}
		return values[len(values)-1]
	}
	defer func() { rollD10 = old }()
	fn()
}

// 场景：普通 DP3，无加骰、无附加
// 骰：8, 3, 10 → 成功 2（8 与 10）；10 触发加骰再投 2（失败）
// 预期：自然 2，最终 2
func TestBasicSuccess(t *testing.T) {
	withSeq([]int{8, 3, 10, 2}, func() {
		r := Roll(3, 10, 0)
		if r.NaturalSuccess != 2 {
			t.Fatalf("自然成功期望2，得%d 骰%v 加%v", r.NaturalSuccess, r.Dice, r.Explosions)
		}
		if r.FinalSuccess != 2 {
			t.Fatalf("最终期望2，得%d", r.FinalSuccess)
		}
		if r.ChanceDie {
			t.Fatal("不应为机运骰")
		}
	})
}

// 场景：附加成功 +2，自然成功 > 0
// 骰：9 → 自然1 + 附加2 → 最终3
func TestBonusApplied(t *testing.T) {
	withSeq([]int{9}, func() {
		r := Roll(1, 10, 2)
		if !r.BonusApplied || r.FinalSuccess != 3 {
			t.Fatalf("期望最终3且计入附加，得 %+v", r)
		}
	})
}

// 场景：自然成功 0，附加成功再多也不计入
// 骰：1,2,3 → 自然0 → 最终0
func TestBonusNotAppliedWhenZeroNatural(t *testing.T) {
	withSeq([]int{1, 2, 3}, func() {
		r := Roll(3, 10, 5)
		if r.NaturalSuccess != 0 || r.FinalSuccess != 0 || r.BonusApplied {
			t.Fatalf("期望最终0且不计入附加，得 %+v", r)
		}
	})
}

// 场景：附加成功为负，自然成功足够
// 骰：8,9 → 自然2 + 附加(-3) → 最终0（封底）
func TestNegativeBonusClamped(t *testing.T) {
	withSeq([]int{8, 9}, func() {
		r := Roll(2, 10, -3)
		if r.NaturalSuccess != 2 {
			t.Fatalf("自然期望2，得%d", r.NaturalSuccess)
		}
		if r.FinalSuccess != 0 {
			t.Fatalf("最终期望封底0，得%d", r.FinalSuccess)
		}
		if !r.BonusApplied {
			t.Fatal("应计入附加（自然>0）")
		}
	})
}

// 场景：负附加但自然仍有余
// 骰：8,9,10,1 → 自然3（10加骰1失败）+ 附加(-1) → 最终2
func TestNegativeBonusPartial(t *testing.T) {
	withSeq([]int{8, 9, 10, 1}, func() {
		r := Roll(3, 10, -1)
		if r.NaturalSuccess != 3 || r.FinalSuccess != 2 {
			t.Fatalf("期望自然3最终2，得自然%d最终%d 加%v", r.NaturalSuccess, r.FinalSuccess, r.Explosions)
		}
	})
}

// 场景：9 加骰 —— 掷出 9 也再投
// 骰：9, 然后加骰 8 → 自然2
func TestExplodeOn9(t *testing.T) {
	withSeq([]int{9, 8}, func() {
		r := Roll(1, 9, 0)
		if len(r.Explosions) != 1 || r.Explosions[0] != 8 {
			t.Fatalf("期望一次加骰8，得%v", r.Explosions)
		}
		if r.NaturalSuccess != 2 {
			t.Fatalf("期望自然2，得%d", r.NaturalSuccess)
		}
	})
}

// 场景：机运骰成功
// 首骰 10，加骰 3 → 自然1
func TestChanceDieSuccess(t *testing.T) {
	withSeq([]int{10, 3}, func() {
		r := Roll(0, 10, 0)
		if !r.ChanceDie {
			t.Fatal("应为机运骰")
		}
		if r.NaturalSuccess != 1 || r.FinalSuccess != 1 {
			t.Fatalf("期望成功1，得 %+v", r)
		}
		if r.CriticalFailure {
			t.Fatal("不应大失败")
		}
	})
}

// 场景：机运骰大失败 —— 首骰1且无成功
func TestChanceDieCritFail(t *testing.T) {
	withSeq([]int{1}, func() {
		r := Roll(-2, 10, 0)
		if !r.CriticalFailure || r.FinalSuccess != 0 {
			t.Fatalf("期望大失败，得 %+v", r)
		}
	})
}

// 场景：机运骰首骰1但加骰（若 explodeOn 允许）出成功则非大失败
// 10 加骰时 1 不会触发加骰，故用 8 加骰：首骰若为8会触发——改用首骰10测加骰成功
// 本测：首骰1，不触发加骰 → 大失败（已覆盖）
// 另：首骰10，加骰出1 → 有成功，非大失败
func TestChanceDieExplodeNoCritOnLaterOne(t *testing.T) {
	withSeq([]int{10, 1}, func() {
		r := Roll(0, 10, 0)
		if r.CriticalFailure {
			t.Fatalf("已有成功时加骰1不应大失败: %+v", r)
		}
		if r.NaturalSuccess != 1 {
			t.Fatalf("期望自然1，得%d", r.NaturalSuccess)
		}
	})
}

// 场景：加骰连锁
// 10 → 10 → 7：主1成功 + 加骰1成功 → 自然2
func TestExplodeChain(t *testing.T) {
	withSeq([]int{10, 10, 7}, func() {
		r := Roll(1, 10, 0)
		if len(r.Explosions) != 2 {
			t.Fatalf("期望2次加骰，得%v", r.Explosions)
		}
		if r.NaturalSuccess != 2 {
			t.Fatalf("期望自然2，得%d", r.NaturalSuccess)
		}
	})
}
