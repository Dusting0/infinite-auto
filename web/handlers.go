package web

import (
	"encoding/json"
	"html/template"
	"net/http"

	"infinite-calc/defense"
	"infinite-calc/dice"
	"infinite-calc/rest"
)

func jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func jsonErr(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// RegisterRoutes 注册所有 HTTP 路由（状态按浏览器 Session 隔离）
func RegisterRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// 首次访问即建立 session cookie
		_ = sessions.get(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, _ := template.New("page").Parse(PageHTML)
		t.Execute(w, nil)
	})

	// 兼容旧链接：防御与血量已合并到同一页
	http.HandleFunc("/defense", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/#defense", http.StatusFound)
	})

	http.HandleFunc("/api/defense/state", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		jsonOK(w, s.Defense.Snapshot())
	})

	http.HandleFunc("/api/defense/update", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var in defense.Preset
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			jsonErr(w, "参数错误: "+err.Error())
			return
		}
		if err := s.Defense.Update(&in); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.Defense.Snapshot())
	})

	http.HandleFunc("/api/defense/reset", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		s.Defense.Reset()
		jsonOK(w, s.Defense.Snapshot())
	})

	http.HandleFunc("/api/dice/roll", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			DP        int `json:"dp"`
			ExplodeOn int `json:"explodeOn"`
			Bonus     int `json:"bonus"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		jsonOK(w, dice.Roll(req.DP, req.ExplodeOn, req.Bonus))
	})

	http.HandleFunc("/api/defense/resolve", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var in defense.AttackInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		jsonOK(w, s.Defense.ResolveAttack(in))
	})

	http.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/setmax", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var req struct {
			Max int `json:"max"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		if err := s.HP.SetMax(req.Max); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/damage", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var req struct {
			Amount int    `json:"amount"`
			Type   string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		if err := s.HP.ApplyDamage(req.Amount, req.Type); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/heal", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var req struct {
			Amount int    `json:"amount"`
			Type   string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		if err := s.HP.HealType(req.Type, req.Amount); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/shortrest", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		rest.ApplyShortRest(s.HP)
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/longrest/options", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		opt := rest.AvailableLongRestOptions(s.HP)
		jsonOK(w, opt)
	})

	http.HandleFunc("/api/longrest", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var req struct {
			Mode string `json:"mode"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if err := rest.ApplyLongRest(s.HP, req.Mode); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/setparts", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var req struct {
			Intact int `json:"intact"`
			B      int `json:"b"`
			L      int `json:"l"`
			A      int `json:"a"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		if err := s.HP.SetParts(req.Intact, req.B, req.L, req.A); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/setdamages", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var req struct {
			B int `json:"b"`
			L int `json:"l"`
			A int `json:"a"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		if err := s.HP.SetDamages(req.B, req.L, req.A); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/setvalues", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		var req struct {
			Max int `json:"max"`
			B   int `json:"b"`
			L   int `json:"l"`
			A   int `json:"a"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErr(w, "参数错误")
			return
		}
		if err := s.HP.SetValues(req.Max, req.B, req.L, req.A); err != nil {
			jsonErr(w, err.Error())
			return
		}
		jsonOK(w, s.HP.Snapshot())
	})

	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		s := sessions.get(w, r)
		s.HP.Reset()
		jsonOK(w, s.HP.Snapshot())
	})
}
