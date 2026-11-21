package defense

import "testing"

// 场景：伤害详情中的“完美”可独立切换生效，且支持正、负与零加值。
// 参数 / 前置状态：攻击 DP8、基础防御3、完美防御2且生效；完美加值和开关由各用例给定。
// 预期结果：实际 DP = 8 + 生效的完美加值 - 3 - 2；未生效时保留输入但不参与计算。
// 边界：负加值使实际 DP≤0 时保持直接未命中规则，不走机运骰。
func TestPerfectBonusActivation(t *testing.T) {
	cases := []struct {
		name   string
		bonus  int
		active bool
		wantDP int
	}{
		{"inactive positive", 4, false, 3},
		{"inactive negative", -4, false, 3},
		{"active positive", 4, true, 7},
		{"active zero", 0, true, 3},
		{"active negative", -2, true, 1},
		{"zero total", -3, true, 0},
		{"negative total", -5, true, -2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPreset("perfect bonus")
			p.Base = 3
			p.PerfectDefense = 2
			p.PerfectDefenseActive = true
			r := p.ResolveAttack(AttackInput{
				AttackDP: 8, PerfectBonus: tc.bonus, PerfectBonusActive: tc.active,
				ExplodeOn: 10, IsPhysical: true,
			})
			if r.ActualDP != tc.wantDP {
				t.Fatalf("actual DP = %d, want %d", r.ActualDP, tc.wantDP)
			}
			if tc.wantDP <= 0 {
				if !r.Miss || r.MissReason != "实际DP≤0" || r.Roll != nil {
					t.Fatalf("nonpositive DP must miss without rolling: %+v", r)
				}
			} else if r.Roll == nil || r.Roll.DP != tc.wantDP {
				t.Fatalf("roll must use actual DP %d, got %+v", tc.wantDP, r.Roll)
			}
		})
	}
}
