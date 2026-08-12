package web

// DefenseHTML 防御预设页面
const DefenseHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>无限规则计算器 · 防御预设</title>
<style>
  :root {
    --bg:#0f1419; --card:#1a2332; --border:#2d3a4f; --text:#e7ecf3;
    --muted:#8b9bb4; --accent:#3b82f6; --ok:#22c55e; --warn:#f59e0b; --danger:#ef4444;
  }
  *{box-sizing:border-box;margin:0;padding:0}
  body{font-family:"Segoe UI","PingFang SC","Microsoft YaHei",sans-serif;background:var(--bg);color:var(--text);min-height:100vh;padding:12px;line-height:1.35}
  .container{max-width:880px;margin:0 auto}
  h1{text-align:center;font-size:1.15rem;margin-bottom:2px;background:linear-gradient(90deg,#60a5fa,#a78bfa);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
  .nav{display:flex;gap:6px;justify-content:center;margin-bottom:10px}
  .nav a{color:var(--muted);text-decoration:none;padding:5px 10px;border-radius:6px;border:1px solid var(--border);font-size:0.8rem}
  .nav a.active,.nav a:hover{color:#fff;background:#334155;border-color:#475569}
  .card{background:var(--card);border:1px solid var(--border);border-radius:8px;padding:10px 12px;margin-bottom:8px}
  .card h2{font-size:0.7rem;color:var(--muted);margin-bottom:8px;text-transform:uppercase;letter-spacing:0.04em}
  .row{display:flex;flex-wrap:wrap;gap:6px;align-items:center}
  .grid3{display:grid;grid-template-columns:repeat(3,1fr);gap:6px}
  .group{margin-bottom:8px}
  .group:last-child{margin-bottom:0}
  .group-title{font-size:0.72rem;color:#93c5fd;margin-bottom:5px;font-weight:600}
  label.field{display:flex;flex-direction:column;gap:2px;font-size:0.7rem;color:var(--muted)}
  label.field input,label.field select{font:inherit;padding:5px 6px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text);width:100%}
  label.check{display:flex;align-items:center;gap:5px;font-size:0.8rem;color:var(--text);cursor:pointer;padding:5px 8px;background:#0f172a;border-radius:5px;border:1px solid var(--border)}
  label.check input{accent-color:var(--accent)}
  button{font:inherit;cursor:pointer;border-radius:6px;border:none;padding:7px 12px;background:var(--accent);color:#fff;font-weight:600;font-size:0.85rem}
  button.secondary{background:#334155}
  button.danger{background:var(--danger)}
  .totals{display:flex;gap:8px}
  .tot{flex:1;background:#0f172a;border-radius:6px;padding:8px;text-align:center}
  .tot .lab{font-size:0.7rem;color:var(--muted)}
  .tot .val{font-size:1.3rem;font-weight:700}
  .tot.melee .val{color:#60a5fa}
  .tot.ranged .val{color:#a78bfa}
  .hint{font-size:0.68rem;color:var(--muted);margin-top:4px}
  .pool{font-size:0.75rem;color:var(--muted);margin-top:4px}
  .eff{color:var(--ok);font-size:0.68rem}
  textarea{width:100%;min-height:48px;font:inherit;padding:6px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  /* 近战／远程斜线 */
  .slash-pair{display:flex;align-items:center;gap:4px}
  .slash-pair input{width:52px;text-align:center;font:inherit;padding:5px 4px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .slash-pair .slash{color:var(--muted);font-size:1rem;font-weight:600;user-select:none}
  .slash-pair .sub{font-size:0.62rem;color:var(--muted);text-align:center}
  .slash-col{display:flex;flex-direction:column;gap:1px;align-items:center}
  /* 左右分栏 */
  .split{display:grid;grid-template-columns:1.55fr 1fr;gap:12px}
  @media (max-width:700px){.split{grid-template-columns:1fr}}
  .split-right{display:flex;flex-direction:column;gap:8px;border-left:1px solid var(--border);padding-left:12px}
  @media (max-width:700px){.split-right{border-left:none;padding-left:0;border-top:1px solid var(--border);padding-top:8px}}
  .r-field{display:flex;flex-direction:column;gap:3px}
  .r-field > span{font-size:0.7rem;color:var(--muted)}
  .r-field input,.r-field select{font:inherit;padding:5px 6px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .absorb-line{display:flex;gap:6px;align-items:center}
  .absorb-line input{width:64px;text-align:center}
  .absorb-line select{flex:1;min-width:0}
  .dr-line{display:flex;align-items:center;gap:6px}
  .dr-line .dr-lab{font-size:1.05rem;font-weight:700;color:#93c5fd}
  .dr-line input[type=number]{width:56px;text-align:center}
  .dr-line input[type=text]{flex:1;min-width:0}
</style>
</head>
<body>
<div class="container">
  <h1>无限规则计算器 · 防御预设</h1>
  <div class="nav">
    <a href="/">血量</a>
    <a href="/defense" class="active">防御预设</a>
  </div>

  <div class="card">
    <h2>合计</h2>
    <div class="totals">
      <div class="tot melee"><div class="lab">近战防御</div><div class="val" id="tMelee">0</div></div>
      <div class="tot ranged"><div class="lab">远程防御</div><div class="val" id="tRanged">0</div></div>
    </div>
    <p class="pool" id="pools"></p>
  </div>

  <div class="card">
    <h2>防御状态</h2>
    <div class="row">
      <label class="check"><input type="checkbox" id="blocking"> 格挡中</label>
      <label class="check"><input type="checkbox" id="flatFooted"> 措手不及</label>
      <label class="check"><input type="checkbox" id="fullDefense"> 全力防御</label>
      <label class="check"><input type="checkbox" id="cover"> 掩蔽</label>
    </div>
  </div>

  <div class="card">
    <h2>防御成分 / 伤害相关</h2>
    <div class="split">
      <!-- 左侧：防御成分 -->
      <div>
        <div class="group">
          <div class="group-title">抵高速</div>
          <div class="grid3">
            <label class="field">基础防御 <input type="number" id="base" value="0"><span class="eff" id="effBase"></span></label>
            <label class="field">闪避防御 <input type="number" id="dodge" value="0"><span class="eff" id="effDodge"></span></label>
            <label class="field">格挡防御 <input type="number" id="block" value="0"><span class="eff" id="effBlock"></span></label>
          </div>
        </div>
        <div class="group">
          <div class="group-title">抵破甲</div>
          <div class="grid3">
            <label class="field">天生防御 <input type="number" id="natural" value="0"></label>
            <label class="field">盔甲
              <div class="slash-pair">
                <div class="slash-col"><input type="number" id="armorMelee" value="0"><span class="sub">近</span></div>
                <span class="slash">／</span>
                <div class="slash-col"><input type="number" id="armorRanged" value="0"><span class="sub">远</span></div>
              </div>
            </label>
            <label class="field">盾牌 <span class="eff" id="effSM"></span><span class="eff" id="effSR"></span>
              <div class="slash-pair">
                <div class="slash-col"><input type="number" id="shieldMelee" value="0"><span class="sub">近</span></div>
                <span class="slash">／</span>
                <div class="slash-col"><input type="number" id="shieldRanged" value="0"><span class="sub">远</span></div>
              </div>
            </label>
          </div>
        </div>
        <div class="group">
          <div class="group-title">抵破魔</div>
          <div class="grid3" style="margin-bottom:6px">
            <label class="field">力场 <input type="number" id="force" value="0"></label>
            <label class="field">偏斜 <input type="number" id="deflection" value="0"></label>
            <label class="field">洞察 <input type="number" id="insight" value="0"></label>
          </div>
          <div class="grid3">
            <label class="field">掩蔽加值 <input type="number" id="coverBonus" value="4"><span class="eff" id="effCover"></span></label>
            <label class="field">其他2 <input type="number" id="other2" value="0"></label>
            <label class="field">其他3 <input type="number" id="other3" value="0"></label>
          </div>
        </div>
      </div>
      <!-- 右侧：伤害相关 -->
      <div class="split-right">
        <div class="r-field">
          <span>防御附加成功</span>
          <input type="number" id="defenseBonusSuccess" value="0">
        </div>
        <div class="r-field">
          <span>伤害吸收</span>
          <div class="absorb-line">
            <input type="number" id="damageAbsorb" value="0">
            <select id="damageAbsorbType">
              <option value="physical">物理伤害吸收</option>
              <option value="all">全伤害吸收</option>
            </select>
          </div>
        </div>
        <div class="r-field">
          <span>伤害减免</span>
          <div class="dr-line">
            <span class="dr-lab">DR</span>
            <input type="number" id="drValue" value="0" min="0">
            <span class="slash">／</span>
            <input type="text" id="drType" placeholder="神兵" value="">
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="card">
    <h2>特性 / 备注</h2>
    <label class="field">特性（防弹、防能量武器等）
      <input type="text" id="traits" placeholder="例如：防弹、防能量武器">
    </label>
    <label class="field" style="margin-top:6px">备注
      <textarea id="notes"></textarea>
    </label>
  </div>

  <div class="row" style="justify-content:flex-end;margin-top:4px">
    <button type="button" class="secondary" id="btnSave">保存预设</button>
    <button type="button" class="danger" id="btnReset">清空</button>
  </div>
</div>

<script>
const fields = [
  'base','dodge','block','natural','armorMelee','armorRanged',
  'shieldMelee','shieldRanged','force','deflection','insight',
  'coverBonus','other2','other3',
  'defenseBonusSuccess','damageAbsorb','drValue'
];

async function api(path, body) {
  const res = await fetch(path, {
    method: 'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(body||{})
  });
  return res.json();
}

function readForm() {
  const o = {
    name: '默认预设',
    flatFooted: document.getElementById('flatFooted').checked,
    fullDefense: document.getElementById('fullDefense').checked,
    blocking: document.getElementById('blocking').checked,
    cover: document.getElementById('cover').checked,
    traits: document.getElementById('traits').value,
    notes: document.getElementById('notes').value,
    damageAbsorbType: document.getElementById('damageAbsorbType').value,
    drType: document.getElementById('drType').value.trim(),
    customs: []
  };
  fields.forEach(id => {
    o[id] = parseInt(document.getElementById(id).value, 10) || 0;
  });
  return o;
}

function applySnapshot(s) {
  if (!s) return;
  document.getElementById('flatFooted').checked = !!s.flatFooted;
  document.getElementById('fullDefense').checked = !!s.fullDefense;
  document.getElementById('blocking').checked = !!s.blocking;
  document.getElementById('cover').checked = !!s.cover;
  fields.forEach(id => {
    if (s[id] !== undefined) document.getElementById(id).value = s[id];
  });
  document.getElementById('traits').value = s.traits || '';
  document.getElementById('notes').value = s.notes || '';
  document.getElementById('damageAbsorbType').value = s.damageAbsorbType || 'physical';
  document.getElementById('drType').value = s.drType || '';

  document.getElementById('tMelee').textContent = s.totals.melee;
  document.getElementById('tRanged').textContent = s.totals.ranged;
  const t = s.totals;
  document.getElementById('pools').textContent =
    '高速池 近'+t.speedPoolMelee+'/远'+t.speedPoolRanged +
    ' · 破甲池 近'+t.armorPoolMelee+'/远'+t.armorPoolRanged +
    ' · 破魔池 '+t.magicPool +
    ' · 其他 '+t.otherPool;

  const e = s.effective || {};
  document.getElementById('effBase').textContent = e.base !== s.base ? ('→'+e.base) : '';
  document.getElementById('effDodge').textContent = e.dodge !== s.dodge ? ('→'+e.dodge) : '';
  document.getElementById('effBlock').textContent = e.block !== s.block ? ('→'+e.block) : '';
  document.getElementById('effSM').textContent = (e.shieldMelee !== s.shieldMelee) ? ('近→'+e.shieldMelee+' ') : '';
  document.getElementById('effSR').textContent = (e.shieldRanged !== s.shieldRanged) ? ('远→'+e.shieldRanged) : '';
  document.getElementById('effCover').textContent = e.cover !== undefined ? ('→'+e.cover) : '';
}

async function save() {
  const s = await api('/api/defense/update', readForm());
  if (s.error) { alert(s.error); return; }
  applySnapshot(s);
}

async function refresh() {
  const res = await fetch('/api/defense/state');
  applySnapshot(await res.json());
}

document.getElementById('btnSave').onclick = save;
document.getElementById('btnReset').onclick = async () => {
  if (!confirm('清空防御预设？')) return;
  const s = await api('/api/defense/reset', {});
  applySnapshot(s);
};

['flatFooted','fullDefense','blocking','cover'].forEach(id => {
  document.getElementById(id).addEventListener('change', save);
});
fields.forEach(id => {
  const el = document.getElementById(id);
  el.addEventListener('change', save);
  el.addEventListener('keydown', e => { if (e.key === 'Enter') { e.preventDefault(); save(); }});
});
document.getElementById('damageAbsorbType').addEventListener('change', save);
document.getElementById('drType').addEventListener('change', save);
document.getElementById('drType').addEventListener('keydown', e => {
  if (e.key === 'Enter') { e.preventDefault(); save(); }
});

refresh();
</script>
</body>
</html>
`
