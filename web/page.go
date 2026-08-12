package web

// PageHTML 血量+掷骰（左） / 防御预设（右）
const PageHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>无限规则计算器 v0.1</title>
<style>
  :root {
    --bg:#0f1419; --card:#1a2332; --border:#2d3a4f;
    --text:#e7ecf3; --muted:#8b9bb4; --accent:#3b82f6;
    --danger:#ef4444; --ok:#22c55e; --warn:#f59e0b; --b:#38bdf8;
  }
  *{box-sizing:border-box;margin:0;padding:0}
  body{font-family:"Segoe UI","PingFang SC","Microsoft YaHei",sans-serif;background:var(--bg);color:var(--text);min-height:100vh;padding:12px;line-height:1.35}
  .container{max-width:1100px;margin:0 auto}
  h1{text-align:center;font-size:1.15rem;margin-bottom:10px;background:linear-gradient(90deg,#60a5fa,#a78bfa);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
  .layout{display:grid;grid-template-columns:minmax(0,1fr) minmax(0,1.15fr) minmax(0,0.95fr);gap:12px;align-items:start}
  @media (max-width:1100px){.layout{grid-template-columns:1fr 1fr}}@media (max-width:720px){.layout{grid-template-columns:1fr}}
  .col{display:flex;flex-direction:column;gap:8px;min-width:0}
  .card{background:var(--card);border:1px solid var(--border);border-radius:8px;padding:10px 12px}
  .card h2{font-size:0.7rem;color:var(--muted);margin-bottom:8px;text-transform:uppercase;letter-spacing:0.04em}
  .row{display:flex;flex-wrap:wrap;gap:6px;align-items:center}
  .grid3{display:grid;grid-template-columns:repeat(3,1fr);gap:6px}
  .group{margin-bottom:8px}
  .group:last-child{margin-bottom:0}
  .group-title{font-size:0.72rem;color:#93c5fd;font-weight:600}
  .group-head{display:flex;align-items:center;justify-content:space-between;gap:8px;margin-bottom:5px}
  .group-head .resist{display:flex;align-items:center;gap:4px;font-size:0.68rem;color:var(--muted)}
  .group-head .resist input{width:48px;text-align:center;font:inherit;padding:3px 4px;border-radius:4px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .group-head .resist label{white-space:nowrap}
  label.field{display:flex;flex-direction:column;gap:2px;font-size:0.7rem;color:var(--muted)}
  label.field input,label.field select{font:inherit;padding:6px 8px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text);width:100%;box-sizing:border-box}
  label.check{display:flex;align-items:center;gap:5px;font-size:0.8rem;color:var(--text);cursor:pointer;padding:5px 8px;background:#0f172a;border-radius:5px;border:1px solid var(--border)}
  label.check input{accent-color:var(--accent)}
  button{font:inherit;cursor:pointer;border-radius:6px;border:none;padding:7px 12px;background:var(--accent);color:#fff;font-weight:600;font-size:0.85rem}
  button.secondary{background:#334155}
  button.danger{background:var(--danger)}
  button:hover{opacity:0.92}
  .hint{font-size:0.68rem;color:var(--muted);margin-top:4px}
  .pool{font-size:0.75rem;color:var(--muted);margin-top:4px}
  .eff{color:var(--ok);font-size:0.68rem}
  textarea{width:100%;min-height:44px;font:inherit;padding:6px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .main-row{display:flex;gap:10px;align-items:stretch}
  .stat-list{display:flex;flex-direction:column;gap:4px;width:100px;flex-shrink:0}
  .stat-item{display:flex;align-items:center;justify-content:space-between;gap:4px;padding:4px 7px;border-radius:6px;background:#0f172a;border:1px solid transparent}
  .stat-item.editable{cursor:pointer}
  .stat-item.editable:hover{border-color:var(--border);background:#132033}
  .stat-item .lab{font-size:0.7rem;color:var(--muted)}
  .stat-item .val{font-size:0.9rem;font-weight:700;font-variant-numeric:tabular-nums;min-width:1.4em;text-align:right}
  .stat-item.max .val{color:#a78bfa}
  .stat-item.intact .val{color:var(--ok)}
  .stat-item.intact.neg .val{color:var(--danger)}
  .stat-item.b .val{color:var(--b)}
  .stat-item.l .val{color:var(--warn)}
  .stat-item.a .val{color:var(--danger)}
  .right-col{flex:1;min-width:0;display:flex;flex-direction:column;gap:6px}
  .bar-header{display:flex;align-items:center;justify-content:space-between;gap:6px}
  .status-badge{font-weight:600;font-size:0.8rem}
  .rest-btns{display:flex;gap:4px;flex-wrap:wrap}
  .rest-btns button{padding:5px 8px;font-size:0.75rem}
  .hp-bar{display:flex;height:28px;border-radius:6px;overflow:hidden;border:1px solid var(--border);background:#0f172a;width:100%}
  .hp-seg{display:flex;align-items:center;justify-content:center;font-size:0.65rem;font-weight:700;color:#0f172a;cursor:pointer;min-width:0}
  .hp-seg:hover{filter:brightness(1.12)}
  .hp-seg.intact{background:var(--ok);cursor:default}
  .hp-seg.b{background:var(--b)}
  .hp-seg.l{background:var(--warn)}
  .hp-seg.a{background:var(--danger)}
  .damage-row{display:flex;gap:6px;flex-wrap:wrap;align-items:center}
  .damage-row .label{font-size:0.75rem;color:var(--muted)}
  .damage-row input{width:56px;font:inherit;padding:5px 6px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .damage-row select{min-width:90px;font:inherit;padding:6px 8px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text);box-sizing:border-box}
  .log{max-height:90px;overflow-y:auto;overflow-x:hidden;font-size:0.72rem;color:var(--muted);font-family:ui-monospace,monospace;white-space:pre-wrap;word-break:break-word}
  .log div{padding:2px 0;border-bottom:1px solid #1e293b}
  .dice-form{display:grid;grid-template-columns:1fr 1fr 1fr auto;gap:6px;align-items:end}
  @media (max-width:500px){.dice-form{grid-template-columns:1fr 1fr}}
  .dice-result{margin-top:8px;background:#0f172a;border-radius:6px;padding:8px;min-height:52px}
  .dice-result .final{font-size:1.4rem;font-weight:700;color:#60a5fa}
  .dice-result .final.crit{color:var(--danger)}
  .dice-result .meta{font-size:0.72rem;color:var(--muted);margin-top:2px}
  .die-chips{display:flex;flex-wrap:wrap;gap:4px;margin-top:6px}
  .die{display:inline-flex;align-items:center;justify-content:center;width:26px;height:26px;border-radius:5px;font-size:0.75rem;font-weight:700;background:#1e293b;color:var(--muted)}
  .die.ok{background:#14532d;color:#86efac}
  .die.exp{background:#1e3a5f;color:#93c5fd;border:1px dashed #3b82f6}
  .die.exp.ok{background:#14532d;border-color:#22c55e;color:#86efac}
  .dice-log{max-height:100px;overflow-y:auto;font-size:0.7rem;color:var(--muted);margin-top:6px}
  .dice-log div{padding:2px 0;border-bottom:1px solid #1e293b}
  .totals{display:flex;gap:8px}
  .tot{flex:1;background:#0f172a;border-radius:6px;padding:8px;text-align:center}
  .tot .lab{font-size:0.7rem;color:var(--muted)}
  .tot .val{font-size:1.25rem;font-weight:700}
  .tot.melee .val{color:#60a5fa}
  .tot.ranged .val{color:#a78bfa}
  .slash-pair{display:flex;align-items:center;gap:4px}
  .slash-pair input{width:48px;text-align:center;font:inherit;padding:5px 3px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .slash-pair .slash{color:var(--muted);font-size:0.95rem;font-weight:600}
  .slash-pair .sub{font-size:0.6rem;color:var(--muted);text-align:center}
  .slash-col{display:flex;flex-direction:column;gap:1px;align-items:center}
  .split{display:grid;grid-template-columns:1.4fr 1fr;gap:10px}
  @media (max-width:600px){.split{grid-template-columns:1fr}}
  .split-right{display:flex;flex-direction:column;gap:8px;border-left:1px solid var(--border);padding-left:10px}
  @media (max-width:600px){.split-right{border-left:none;padding-left:0;border-top:1px solid var(--border);padding-top:8px}}
  .r-field{display:flex;flex-direction:column;gap:3px}
  .r-field > span{font-size:0.7rem;color:var(--muted)}
  .r-field input,.r-field select{font:inherit;padding:6px 8px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text);box-sizing:border-box;width:100%}
  .absorb-line{display:flex;gap:6px;align-items:center}
  .absorb-line input{width:56px;text-align:center}
  .absorb-line select{flex:1;min-width:0}
  .dr-line{display:flex;align-items:center;gap:5px}
  .dr-line .dr-lab{font-size:1rem;font-weight:700;color:#93c5fd}
  .dr-line input[type=number]{width:48px;text-align:center}
  .dr-line input[type=text]{flex:1;min-width:0}
  .trait-list{display:flex;flex-direction:column;gap:4px}
  .trait-row{display:flex;align-items:center;gap:4px}
  .trait-row input{flex:1;font:inherit;padding:5px 8px;border-radius:5px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .trait-row button{padding:4px 7px;font-size:0.72rem;min-width:26px}
  .trait-row .ord{color:var(--muted);font-size:0.68rem;width:14px;text-align:center}
  .modal-mask{display:none;position:fixed;inset:0;background:rgba(0,0,0,0.55);z-index:1000;align-items:center;justify-content:center;padding:16px}
  .modal-mask.show{display:flex}
  .modal{background:var(--card);border:1px solid var(--border);border-radius:12px;padding:18px 20px;width:100%;max-width:340px;box-shadow:0 16px 48px rgba(0,0,0,0.45)}
  .modal h3{font-size:1rem;margin-bottom:6px}
  .modal .desc{font-size:0.8rem;color:var(--muted);margin-bottom:12px}
  .modal input{width:100%;margin-bottom:14px;font-size:1.1rem;text-align:center;font:inherit;padding:8px;border-radius:6px;border:1px solid var(--border);background:#0f172a;color:var(--text)}
  .modal-actions{display:flex;gap:8px;justify-content:flex-end}
  .modal-actions button{min-width:72px}
  .modal-choices{display:flex;flex-direction:column;gap:8px;margin-bottom:12px}

  input[type=number]::-webkit-inner-spin-button,
  input[type=number]::-webkit-outer-spin-button{-webkit-appearance:none;margin:0}
  input[type=number]{-moz-appearance:textfield;appearance:textfield}
  .modal-choices button{width:100%;text-align:left;background:#0f172a;border:1px solid var(--border);color:var(--text)}
  .atk-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:6px;margin-top:6px}
  .atk-result{margin-top:8px;background:#0f172a;border-radius:6px;padding:8px;font-size:0.75rem;color:var(--muted);line-height:1.5;min-height:40px}
  .atk-result .big{font-size:1.3rem;font-weight:700;color:#60a5fa}
  .atk-result .miss{color:var(--danger);font-weight:700;font-size:1.1rem}

  .num-tap{
    display:inline-flex;align-items:center;justify-content:center;
    min-width:2.2em;padding:4px 8px;border-radius:5px;
    background:#0f172a;border:1px solid var(--border);
    color:var(--text);font-weight:700;font-variant-numeric:tabular-nums;
    cursor:pointer;font-size:0.9rem;transition:border-color .15s,background .15s;
    user-select:none;
  }
  .num-tap:hover{border-color:var(--accent);background:#132033}
  .num-tap:active{transform:scale(0.97)}
  label.field .num-tap{width:100%;box-sizing:border-box;padding:6px 8px;justify-content:flex-start}
  .group-head .resist .num-tap{min-width:2.4em;padding:3px 6px;font-size:0.8rem}
  .absorb-line .num-tap{min-width:2.6em}
  .dr-line .num-tap{min-width:2.4em}
  .damage-row .num-tap{min-width:2.6em}
  .dice-form .num-tap{width:100%;justify-content:flex-start;padding:6px 8px}
  .atk-grid .num-tap{width:100%;justify-content:flex-start;padding:6px 8px}



  .footer-credit{text-align:center;color:var(--muted);font-size:0.72rem;margin-top:14px;margin-bottom:6px;opacity:0.8;letter-spacing:0.02em}

  .modal-num-row{display:flex;align-items:center;gap:8px;margin-bottom:14px}
  .modal-num-row input{width:100%;margin-bottom:0 !important;font-size:1.1rem;text-align:center}
  .modal-step{
    flex-shrink:0;width:40px;height:40px;padding:0;
    border-radius:8px;font-size:1.25rem;font-weight:700;
    display:inline-flex;align-items:center;justify-content:center;
    background:#334155;color:#fff;border:none;cursor:pointer;
  }
  .modal-step:hover{opacity:0.9}

  label.field select{line-height:1.35;min-height:34px}
  label.field .num-tap{min-height:34px}

  .save-grid{display:flex;flex-direction:column;gap:6px}
  .save-row{display:grid;grid-template-columns:48px 1fr 1fr auto minmax(72px,auto);gap:6px;align-items:end}
  .save-name{font-size:0.8rem;font-weight:600;color:#93c5fd;padding-bottom:8px}
  .save-out{font-size:0.85rem;font-weight:700;color:#60a5fa;padding-bottom:6px;text-align:right;min-width:72px}
  @media (max-width:500px){.save-row{grid-template-columns:40px 1fr 1fr;}.save-out{grid-column:1/-1;text-align:left;padding-bottom:0}}

  .dice-log{max-height:240px;overflow-y:auto;overflow-x:hidden;font-size:0.72rem;color:var(--muted);font-family:ui-monospace,monospace;word-break:break-word;white-space:pre-wrap}
  .dice-log div{padding:3px 0;border-bottom:1px solid #1e293b;line-height:1.4}
  .dice-log .src{color:#93c5fd;margin-right:4px}
  .dice-log .faces{color:#e2e8f0}

  .edit-label{display:inline-flex;align-items:center;gap:3px;max-width:100%;cursor:text;min-width:0}
  .edit-label-text{
    display:inline-block;max-width:5.5em;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;
    font-size:0.72rem;color:var(--muted);line-height:1.2;
  }
  .edit-label .edit-ico{font-size:0.65rem;opacity:0.55;flex-shrink:0}
  .edit-label:hover .edit-ico{opacity:0.9;color:var(--accent)}
  .edit-label input.inline-edit{
    border:none;outline:none;background:transparent;color:var(--text);
    font:inherit;font-size:0.72rem;padding:0;margin:0;width:5.5em;min-width:0;
  }
  label.field > .edit-label{margin-bottom:2px}

  .dice-result{margin-top:8px;padding:10px;border:1px solid var(--border);border-radius:8px;background:#0f172a}
  .dice-result .succ-block{display:flex;align-items:stretch;gap:10px;margin-bottom:8px}
  .dice-result .succ-bar{width:3px;border-radius:2px;background:linear-gradient(180deg,#60a5fa,#2563eb);flex-shrink:0}
  .dice-result .succ-body{min-width:0}
  .dice-result .succ-lab{font-size:0.72rem;color:var(--muted);margin-bottom:2px}
  .dice-result .succ-num{font-size:2.4rem;font-weight:800;line-height:1;color:#60a5fa;font-variant-numeric:tabular-nums}
  .dice-result .succ-num.crit{color:var(--danger)}
  .dice-result .meta{font-size:0.75rem;color:var(--muted);margin-bottom:8px;word-break:break-word}
  .dice-result .faces-lab{font-size:0.7rem;color:var(--muted);margin-bottom:6px}
  .dice-result .dice-chips{display:flex;flex-wrap:wrap;align-items:center;gap:6px}
  .dice-result .die{
    display:inline-flex;align-items:center;justify-content:center;
    min-width:32px;height:32px;padding:0 6px;border-radius:8px;
    background:#1e293b;border:1px solid #334155;color:#e2e8f0;
    font-weight:700;font-size:0.9rem;font-variant-numeric:tabular-nums;
  }
  .dice-result .die.ok{border-color:#34d399;color:#6ee7b7;background:#0f2a22}
  .dice-result .die.exp.ok{border-color:#34d399}
  .dice-result .die-plus{
    display:inline-flex;align-items:center;justify-content:center;
    padding:0 2px;border:none;background:transparent;
    color:#64748b;font-size:1rem;font-weight:600;line-height:1;
    min-width:0;height:auto;
  }
</style>
</head>
<body>
<div class="container">
  <h1>无限规则计算器 <span style="font-size:0.75rem;opacity:0.7;font-weight:500">v0.1</span></h1>
  <div class="layout">
    <div class="col">
      <div class="card">
        <h2>血量</h2>
        <div class="main-row">
          <div class="stat-list">
            <div class="stat-item max editable" id="rowMax"><span class="lab">上限</span><span class="val" id="max">20</span></div>
            <div class="stat-item intact" id="rowIntact"><span class="lab">完好</span><span class="val" id="intact">20</span></div>
            <div class="stat-item b editable" id="rowB"><span class="lab">冲击</span><span class="val" id="b">0</span></div>
            <div class="stat-item l editable" id="rowL"><span class="lab">严重</span><span class="val" id="l">0</span></div>
            <div class="stat-item a editable" id="rowA"><span class="lab">恶性</span><span class="val" id="a">0</span></div>
          </div>
          <div class="right-col">
            <div class="bar-header">
              <span class="status-badge" id="status">正常</span>
              <div class="rest-btns">
                <button type="button" class="secondary" id="btnShort">短休</button>
                <button type="button" class="secondary" id="btnLong">长休</button>
                <button type="button" class="danger" id="btnHpReset">重置</button>
              </div>
            </div>
            <div class="hp-bar" id="hpBar"></div>
            <div class="damage-row">
              <span class="label">受伤</span>
              <span class="num-tap" id="dmgAmount" data-value="1">1</span>
              <select id="dmgType">
                <option value="冲击">冲击</option>
                <option value="严重" selected>严重</option>
                <option value="恶性">恶性</option>
              </select>
              <button type="button" id="btnDamage">应用</button>
            </div>
          </div>
        </div>
        
        <div style="margin-top:10px;padding-top:8px;border-top:1px solid var(--border)">
          <div style="font-size:0.72rem;color:#93c5fd;font-weight:600;margin-bottom:6px">自定义伤害</div>
          <div class="row" style="margin-bottom:6px">
            <label class="check"><input type="checkbox" id="atkPhysical" checked> 物理伤害</label>
          </div>
          <div class="atk-grid">
            <label class="field">攻击DP<span class="num-tap" id="atkDP" data-value="10">10</span></label>
            <label class="field">加骰
              <select id="atkExplode">
                <option value="10" selected>10 加骰</option>
                <option value="9">9 加骰</option>
                <option value="8">8 加骰</option>
              </select>
            </label>
            <label class="field">附加成功<span class="num-tap" id="atkBonus" data-value="0">0</span></label>
            <label class="field">高速<span class="num-tap" id="atkSpeed" data-value="0">0</span></label>
            <label class="field">破甲<span class="num-tap" id="atkAP" data-value="0">0</span></label>
            <label class="field">破魔<span class="num-tap" id="atkMP" data-value="0">0</span></label>
            <label class="field">伤害级别
              <select id="atkDmgType">
                <option value="冲击">冲击</option>
                <option value="严重" selected>严重</option>
                <option value="恶性">恶性</option>
              </select>
            </label>
            <label class="field">伤害上限<span class="num-tap" id="atkDmgLimit" data-value="0" title="≤0 表示不限制">0</span></label>
          </div>
          <div class="row" style="margin-top:6px">
            <button type="button" id="btnAtkResolve">结算伤害</button>
          </div>
          <div class="atk-result" id="atkResult"><div>填写攻击参数后点击结算。结果仅展示，不自动扣血。</div></div>
        </div>

        <div class="log" id="log" style="margin-top:6px"></div>
      </div>

      <div class="card">
        <h2>三豁免</h2>
        <p class="hint" style="margin-bottom:6px">强韧＝耐力+求生 · 反射＝敏捷+运动 · 意志＝决心+感受（自我保护）。默认10加骰。</p>
        <div class="save-grid">
          <div class="save-row">
            <span class="save-name">强韧</span>
            <label class="field">DP<span class="num-tap" id="saveFortDP" data-value="5">5</span></label>
            <label class="field">附加成功<span class="num-tap" id="saveFortBonus" data-value="0">0</span></label>
            <button type="button" class="secondary" id="btnSaveFort">掷</button>
            <span class="save-out" id="saveFortOut">—</span>
          </div>
          <div class="save-row">
            <span class="save-name">反射</span>
            <label class="field">DP<span class="num-tap" id="saveRefDP" data-value="5">5</span></label>
            <label class="field">附加成功<span class="num-tap" id="saveRefBonus" data-value="0">0</span></label>
            <button type="button" class="secondary" id="btnSaveRef">掷</button>
            <span class="save-out" id="saveRefOut">—</span>
          </div>
          <div class="save-row">
            <span class="save-name">意志</span>
            <label class="field">DP<span class="num-tap" id="saveWillDP" data-value="5">5</span></label>
            <label class="field">附加成功<span class="num-tap" id="saveWillBonus" data-value="0">0</span></label>
            <button type="button" class="secondary" id="btnSaveWill">掷</button>
            <span class="save-out" id="saveWillOut">—</span>
          </div>
        </div>
      </div>
    </div>
    <div class="col" id="defense">
      <div class="card">
        <h2>防御合计</h2>
        <div class="totals">
          <div class="tot melee"><div class="lab">高速防御</div><div class="val" id="tSpeed">0</div></div>
          <div class="tot ranged"><div class="lab">破甲防御</div><div class="val" id="tArmor">0</div></div>
          <div class="tot" style="flex:1;background:#0f172a;border-radius:6px;padding:8px;text-align:center"><div class="lab" style="font-size:0.7rem;color:var(--muted)">破魔防御</div><div class="val" id="tMagic" style="font-size:1.25rem;font-weight:700;color:#34d399">0</div></div>
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
          <label class="check"><input type="checkbox" id="touchAttack" title="接触攻击：检定无视天生/盔甲/盾牌"> 接触攻击</label>
        </div>
      </div>
      <div class="card">
        <h2>防御成分 / 伤害相关</h2>
        <div class="split">
          <div>
            <div class="group">
              <div class="group-head">
                <div class="group-title">高速防御</div>
                <div class="resist"><label>抵高速</label><span class="num-tap" id="resistSpeed" data-value="0">0</span></div>
              </div>
              <div class="grid3">
                <label class="field">基础 <span class="num-tap" id="base" data-value="0">0</span><span class="eff" id="effBase"></span></label>
                <label class="field">闪避 <span class="num-tap" id="dodge" data-value="0">0</span><span class="eff" id="effDodge"></span></label>
                <label class="field">格挡 <span class="num-tap" id="block" data-value="0">0</span><span class="eff" id="effBlock"></span></label>
              </div>
            </div>
            <div class="group">
              <div class="group-head">
                <div class="group-title">破甲防御</div>
                <div class="resist"><label>抵破甲</label><span class="num-tap" id="resistAP" data-value="0">0</span></div>
              </div>
              <div class="grid3">
                <label class="field">天生 <span class="num-tap" id="natural" data-value="0">0</span></label>
                <label class="field">盔甲 <span class="num-tap" id="armorMelee" data-value="0">0</span></label>
                <label class="field">盾牌 <span class="eff" id="effSM"></span>
                  <span class="num-tap" id="shieldMelee" data-value="0">0</span>
                </label>
                <input type="hidden" id="armorRanged" value="0">
                <input type="hidden" id="shieldRanged" value="0">
              </div>
            </div>
            <div class="group">
              <div class="group-head">
                <div class="group-title">破魔防御</div>
                <div class="resist"><label>抵破魔</label><span class="num-tap" id="resistMagic" data-value="0">0</span></div>
              </div>
              <div class="grid3" style="margin-bottom:6px">
                <label class="field">力场 <span class="num-tap" id="force" data-value="0">0</span></label>
                <label class="field">偏斜 <span class="num-tap" id="deflection" data-value="0">0</span></label>
                <label class="field">洞察 <span class="num-tap" id="insight" data-value="0">0</span></label>
              </div>
              <div class="grid3">
                <label class="field">掩蔽 <span class="num-tap" id="coverBonus" data-value="4">4</span><span class="eff" id="effCover"></span></label>
                <label class="field"><span class="edit-label" data-name-id="other2Name"><span class="edit-label-text" id="other2Name">其他1</span><span class="edit-ico" title="编辑名称">✎</span></span> <span class="num-tap" id="other2" data-value="0">0</span></label>
                <label class="field"><span class="edit-label" data-name-id="other3Name"><span class="edit-label-text" id="other3Name">其他2</span><span class="edit-ico" title="编辑名称">✎</span></span> <span class="num-tap" id="other3" data-value="0">0</span></label>
              </div>
            </div>
            <div class="def-group">
              <div class="group-head">
                <div class="group-title">完美防御</div>
              </div>
              <div class="def-grid">
                <label class="field">完美防御 <span class="num-tap" id="perfectDefense" data-value="0">0</span></label>
                <label class="check" style="align-self:end;justify-self:start"><input type="checkbox" id="perfectDefenseActive"> 生效</label>
              </div>
            </div>
          </div>
          <div class="split-right">
            <div class="r-field"><span>防御附加成功</span><span class="num-tap" id="defenseBonusSuccess" data-value="0">0</span></div>
            <div class="r-field"><span>伤害吸收</span>
              <div class="absorb-line">
                <span class="num-tap" id="damageAbsorb" data-value="0">0</span>
                <select id="damageAbsorbType"><option value="physical">物理</option><option value="all">全伤害</option></select>
              </div>
            </div>
            <div class="r-field"><span>物理伤害减免</span>
              <div class="dr-line">
                <span class="dr-lab">DR</span>
                <span class="num-tap" id="drValue" data-value="0">0</span>
                <span class="slash">／</span>
                <input type="text" id="drType" placeholder="神兵" value="">
              </div>
            </div>
            <div class="r-field"><span>能量抗力</span>
              <div class="dr-line">
                <span class="dr-lab">ER</span>
                <span class="num-tap" id="erValue" data-value="0">0</span>
                <span class="slash">／</span>
                <input type="text" id="erType" placeholder="能量类型" value="">
              </div>
            </div>
          </div>
        </div>

        <div class="row" style="justify-content:flex-end;margin-top:8px">
          <button type="button" class="danger" id="btnDefReset">清空防御</button>
        </div>
      </div>
      <div class="card">
        <h2>特性</h2>
        <div class="trait-list" id="traitList"></div>
        <p class="hint">回车新增 · ↑↓排序 · ×删除</p>
      </div>

    </div>
    <div class="col" id="diceCol">
      <div class="card">
        <h2>掷骰</h2>
        <div class="dice-form">
          <label class="field">DP<span class="num-tap" id="diceDP" data-value="5">5</span></label>
          <label class="field">加骰
            <select id="diceExplode">
              <option value="10" selected>10 加骰</option>
              <option value="9">9 加骰</option>
              <option value="8">8 加骰</option>
            </select>
          </label>
          <label class="field">附加成功<span class="num-tap" id="diceBonus" data-value="0">0</span></label>
          <button type="button" id="btnRoll">掷骰</button>
        </div>
        <p class="hint">DP&le;0 自动机运骰（仅10成功；首骰1且无成功=大失败）。附加成功仅在自然成功&gt;0时计入，可为负，最终&ge;0。</p>
        
        <div class="dice-result" id="diceResult"><div class="meta">等待掷骰…</div></div>
        <div class="dice-log" id="diceLog"></div>
      </div>
    </div>
  </div>
  <p class="footer-credit">Powered by Dusting</p>
</div>
<div class="modal-mask" id="modalMask">
  <div class="modal" role="dialog">
    <h3 id="modalTitle">标题</h3>
    <div class="desc" id="modalDesc"></div>
    <div id="modalBody"></div>
    <div class="modal-actions" id="modalActions"></div>
  </div>
</div>
<script>
let modalClosed=null;
async function api(path,body){const res=await fetch(path,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body||{})});const data=await res.json();if(data.error){await showAlert('错误',data.error);return null;}return data;}
function hideModal(){document.getElementById('modalMask').classList.remove('show');}
function finishModal(v){hideModal();const fn=modalClosed;modalClosed=null;if(fn)fn(v);}
function showAlert(title,msg){return new Promise(r=>{modalClosed=r;document.getElementById('modalTitle').textContent=title||'';document.getElementById('modalDesc').textContent=msg||'';document.getElementById('modalBody').innerHTML='';const act=document.getElementById('modalActions');act.innerHTML='';const ok=document.createElement('button');ok.type='button';ok.textContent='确定';ok.onclick=()=>finishModal(true);act.appendChild(ok);document.getElementById('modalMask').classList.add('show');ok.focus();});}
function showConfirm(title,msg){return new Promise(r=>{modalClosed=r;document.getElementById('modalTitle').textContent=title||'';document.getElementById('modalDesc').textContent=msg||'';document.getElementById('modalBody').innerHTML='';const act=document.getElementById('modalActions');act.innerHTML='';const c=document.createElement('button');c.type='button';c.textContent='取消';c.className='secondary';c.onclick=()=>finishModal(false);const ok=document.createElement('button');ok.type='button';ok.textContent='确定';ok.onclick=()=>finishModal(true);act.appendChild(c);act.appendChild(ok);document.getElementById('modalMask').classList.add('show');ok.focus();});}
function showNumberPrompt(title,desc,def,allowNeg){return new Promise(r=>{modalClosed=r;document.getElementById('modalTitle').textContent=title||'';document.getElementById('modalDesc').textContent=desc||'';const minAttr=allowNeg?'':' min="0"';const val=(def==null?0:def);
document.getElementById('modalBody').innerHTML='<div class="modal-num-row"><button type="button" class="modal-step" id="modalDec">−</button><input type="number" id="modalInput"'+minAttr+' step="1" value="'+val+'"><button type="button" class="modal-step" id="modalInc">+</button></div>';
const act=document.getElementById('modalActions');act.innerHTML='';
const step=delta=>{const el=document.getElementById('modalInput');let n=parseInt(el.value,10);if(isNaN(n))n=0;n+=delta;if(!allowNeg&&n<0)n=0;el.value=n;el.focus();el.select();};
document.getElementById('modalDec').onclick=()=>step(-1);
document.getElementById('modalInc').onclick=()=>step(1);
const c=document.createElement('button');c.type='button';c.textContent='取消';c.className='secondary';c.onclick=()=>finishModal(null);
const ok=document.createElement('button');ok.type='button';ok.id='modalOk';ok.textContent='确定';ok.onclick=()=>{const n=parseInt(document.getElementById('modalInput').value,10);finishModal(isNaN(n)?null:n);};
act.appendChild(c);act.appendChild(ok);
document.getElementById('modalMask').classList.add('show');
setTimeout(()=>{const el=document.getElementById('modalInput');if(el){el.focus();el.select();}},40);});}
function showChoice(title,desc,choices){return new Promise(r=>{modalClosed=r;document.getElementById('modalTitle').textContent=title||'';document.getElementById('modalDesc').textContent=desc||'';let h='<div class="modal-choices">';choices.forEach((c,i)=>{h+='<button type="button" data-idx="'+i+'">'+c.label+'</button>';});h+='</div>';document.getElementById('modalBody').innerHTML=h;const act=document.getElementById('modalActions');act.innerHTML='';const c=document.createElement('button');c.type='button';c.textContent='取消';c.className='secondary';c.onclick=()=>finishModal(null);act.appendChild(c);document.querySelectorAll('.modal-choices button').forEach(btn=>{btn.onclick=()=>finishModal(choices[parseInt(btn.getAttribute('data-idx'),10)].value);});document.getElementById('modalMask').classList.add('show');});}

function getNum(id){const el=document.getElementById(id);if(!el)return 0;if(el.tagName==='INPUT'){const n=parseInt(el.value,10);return isNaN(n)?0:n;}const v=el.dataset?el.dataset.value:null;if(v!=null&&v!==''){const n=parseInt(v,10);return isNaN(n)?0:n;}const n=parseInt(el.textContent,10);return isNaN(n)?0:n;}
function setNum(id,v){const el=document.getElementById(id);if(!el)return;const n=parseInt(v,10);const val=isNaN(n)?0:n;if(el.tagName==='INPUT'){el.value=String(val);return;}if(el.dataset)el.dataset.value=String(val);el.textContent=String(val);}
async function editNum(id,label,opts){opts=opts||{};const cur=getNum(id);const n=await showNumberPrompt(label||id,opts.desc||('当前 '+cur),cur,!!opts.allowNeg);if(n===null)return;if(opts.min!=null&&n<opts.min){await showAlert('错误','不能小于 '+opts.min);return;}setNum(id,n);if(opts.onChange)await opts.onChange();}

let current={intact:20,b:0,l:0,a:0,max:20,total:0};
function renderBar(s){const bar=document.getElementById('hpBar');bar.innerHTML='';const max=Math.max(s.max,1);const parts=[{key:'intact',val:Math.max(s.intact,0),cls:'intact',label:'完好',ed:false},{key:'b',val:s.b,cls:'b',label:'冲击',ed:true},{key:'l',val:s.l,cls:'l',label:'严重',ed:true},{key:'a',val:s.a,cls:'a',label:'恶性',ed:true}];let any=false;parts.forEach(p=>{if(p.val<=0)return;any=true;const seg=document.createElement('div');seg.className='hp-seg '+p.cls;seg.style.flex=p.val;seg.textContent=(p.val>=2||p.val/max>0.1)?p.val:'';seg.title=p.label+' '+p.val;if(p.ed)seg.addEventListener('click',()=>editDamage(p.key,p.label));bar.appendChild(seg);});if(!any)bar.innerHTML='<div style="flex:1;display:flex;align-items:center;justify-content:center;color:var(--muted);font-size:0.7rem">点击左侧编辑</div>';}
function renderHp(s){current={intact:s.intact,b:s.b,l:s.l,a:s.a,max:s.max,total:s.total};document.getElementById('intact').textContent=s.intact;document.getElementById('b').textContent=s.b;document.getElementById('l').textContent=s.l;document.getElementById('a').textContent=s.a;document.getElementById('max').textContent=s.max;document.getElementById('status').textContent=s.status;const ri=document.getElementById('rowIntact');if(s.intact<0)ri.classList.add('neg');else ri.classList.remove('neg');renderBar(s);const logEl=document.getElementById('log');if(!s.log||!s.log.length)logEl.innerHTML='<div style="color:#64748b">暂无记录</div>';else logEl.innerHTML=s.log.slice().reverse().map(l=>'<div>'+l+'</div>').join('');}
async function editMax(){const n=await showNumberPrompt('修改上限','当前 '+current.max,current.max);if(n===null)return;if(n<=0){await showAlert('错误','上限须为正');return;}const s=await api('/api/setmax',{max:n});if(s)renderHp(s);}
async function editDamage(key,label){const old=current[key]||0;const n=await showNumberPrompt(old===0?('插入「'+label+'」'):('修改「'+label+'」'),'当前 '+old,old);if(n===null)return;if(n<0){await showAlert('错误','不能为负');return;}const d={b:current.b,l:current.l,a:current.a};d[key]=n;const s=await api('/api/setdamages',d);if(s)renderHp(s);}
async function applyDamage(){const s=await api('/api/damage',{amount:getNum('dmgAmount'),type:document.getElementById('dmgType').value});if(s)renderHp(s);}
async function shortRest(){const s=await api('/api/shortrest',{});if(s)renderHp(s);}
async function longRest(){const opt=await(await fetch('/api/longrest/options')).json();let mode=null;if(opt.CanClearL&&opt.CanConvertA){mode=await showChoice('长休','请选择：',[{label:'L · 严重→完好',value:'L'},{label:'A · 恶性→严重',value:'A'}]);if(mode===null)return;}else if(opt.CanClearL)mode='L';else if(opt.CanConvertA)mode='A';else{await showAlert('长休','无严重/恶性');return;}const s=await api('/api/longrest',{mode});if(s)renderHp(s);}
async function resetHp(){if(!(await showConfirm('重置','重置为满血？')))return;const s=await api('/api/reset',{});if(s)renderHp(s);}
document.getElementById('rowMax').onclick=editMax;
document.getElementById('rowB').onclick=()=>editDamage('b','冲击');
document.getElementById('rowL').onclick=()=>editDamage('l','严重');
document.getElementById('rowA').onclick=()=>editDamage('a','恶性');
document.getElementById('dmgAmount').onclick=()=>editNum('dmgAmount','受到伤害',{min:1});
document.getElementById('btnDamage').onclick=applyDamage;
document.getElementById('btnShort').onclick=shortRest;
document.getElementById('btnLong').onclick=longRest;
document.getElementById('btnHpReset').onclick=resetHp;
const diceLogLines=[];
function formatDiceFaces(r){
  if(r.diceLog)return r.diceLog;
  if(r.batches&&r.batches.length){
    return r.batches.map(b=>'['+(b||[]).join(',')+']').join('+');
  }
  const parts=[];
  if(r.dice&&r.dice.length)parts.push('['+r.dice.join(',')+']');
  (r.explosions||[]).forEach(v=>parts.push('['+v+']'));
  return parts.join('+')||'—';
}
function pushDiceLog(r,source){
  const log=document.getElementById('diceLog');
  if(!log)return;
  const faces=formatDiceFaces(r);
  const line=document.createElement('div');
  const src=source?('<span class="src">['+source+']</span> '):'';
  let tail=r.criticalFailure?'大失败':('成功 '+r.finalSuccess);
  line.innerHTML=src+'<span class="faces">'+faces+'</span> · '+tail+(r.summary?' · '+r.summary:'');
  log.insertBefore(line,log.firstChild);
}
function renderDiceResult(r,source){
  const box=document.getElementById('diceResult');
  if(!box)return;
  const crit=!!r.criticalFailure;
  const final=crit?0:(r.finalSuccess!=null?r.finalSuccess:0);
  let chips='';
  const batches=r.batches&&r.batches.length?r.batches:[r.dice||[]];
  batches.forEach((batch,bi)=>{
    if(bi)chips+='<span class="die-plus">+</span>';
    (batch||[]).forEach(v=>{
      const ok=r.chanceDie&&bi===0?(v===10):(v>=8);
      chips+='<span class="die'+(bi?' exp':'')+(ok?' ok':'')+'">'+v+'</span>';
    });
  });
  const src=source?('['+source+'] '):'';
  const summary=r.summary||'';
  box.innerHTML=
    '<div class="succ-block">'+
      '<div class="succ-bar"></div>'+
      '<div class="succ-body">'+
        '<div class="succ-lab">最终成功</div>'+
        '<div class="succ-num'+(crit?' crit':'')+'">'+(crit?'×':final)+'</div>'+
      '</div>'+
    '</div>'+
    '<div class="meta">'+src+(crit?'大失败 · ':'')+summary+'</div>'+
    (chips?('<div class="faces-lab">骰面</div><div class="dice-chips">'+chips+'</div>'):'');
  pushDiceLog(r,source||'掷骰');
}

async function doRoll(){const dp=getNum('diceDP');const explodeOn=parseInt(document.getElementById('diceExplode').value,10)||10;const bonus=getNum('diceBonus');const r=await api('/api/dice/roll',{dp,explodeOn,bonus});if(r)renderDiceResult(r,'掷骰');}
document.getElementById('diceDP').onclick=()=>editNum('diceDP','DP');
document.getElementById('diceBonus').onclick=()=>editNum('diceBonus','附加成功',{allowNeg:true});
document.getElementById('btnRoll').onclick=doRoll;
let traits=[''];
const defFields=['base','dodge','block','natural','armorMelee','armorRanged','shieldMelee','shieldRanged','force','deflection','insight','coverBonus','other2','other3','perfectDefense','resistSpeed','resistAP','resistMagic','defenseBonusSuccess','damageAbsorb','drValue','erValue'];
function renderTraits(focusIdx){const box=document.getElementById('traitList');box.innerHTML='';if(!traits.length)traits=[''];traits.forEach((t,i)=>{const row=document.createElement('div');row.className='trait-row';row.innerHTML='<span class="ord">'+(i+1)+'</span><button type="button" class="secondary" data-act="up" data-i="'+i+'">↑</button><button type="button" class="secondary" data-act="down" data-i="'+i+'">↓</button><input type="text" data-i="'+i+'" value="'+(t||'').replace(/"/g,'&quot;')+'" placeholder="特性，回车新增"><button type="button" class="danger" data-act="del" data-i="'+i+'">×</button>';box.appendChild(row);});box.querySelectorAll('button[data-act]').forEach(btn=>{btn.onclick=()=>{const i=parseInt(btn.getAttribute('data-i'),10);const act=btn.getAttribute('data-act');syncTraitsFromDom();if(act==='del'){traits.splice(i,1);if(!traits.length)traits=[''];renderTraits();saveDef();}else if(act==='up'&&i>0){const x=traits[i];traits[i]=traits[i-1];traits[i-1]=x;renderTraits(i-1);saveDef();}else if(act==='down'&&i<traits.length-1){const x=traits[i];traits[i]=traits[i+1];traits[i+1]=x;renderTraits(i+1);saveDef();}};});box.querySelectorAll('input[data-i]').forEach(inp=>{inp.addEventListener('change',()=>{syncTraitsFromDom();saveDef();});inp.addEventListener('keydown',e=>{if(e.key!=='Enter')return;e.preventDefault();syncTraitsFromDom();const i=parseInt(inp.getAttribute('data-i'),10);traits.splice(i+1,0,'');renderTraits(i+1);});});if(focusIdx!=null){const el=box.querySelector('input[data-i="'+focusIdx+'"]');if(el)el.focus();}}
function syncTraitsFromDom(){const arr=[];document.querySelectorAll('#traitList input[data-i]').forEach(inp=>arr.push(inp.value));traits=arr.length?arr:[''];}
function traitsForSave(){syncTraitsFromDom();return traits.map(t=>t.trim()).filter(t=>t.length>0);}
function readDefForm(){const o={name:'默认预设',flatFooted:document.getElementById('flatFooted').checked,fullDefense:document.getElementById('fullDefense').checked,blocking:document.getElementById('blocking').checked,cover:document.getElementById('cover').checked,touchAttack:document.getElementById('touchAttack').checked,perfectDefenseActive:document.getElementById('perfectDefenseActive').checked,traits:traitsForSave(),notes:'',other2Name:(document.getElementById('other2Name')&&document.getElementById('other2Name').textContent)||'其他1',other3Name:(document.getElementById('other3Name')&&document.getElementById('other3Name').textContent)||'其他2',damageAbsorbType:document.getElementById('damageAbsorbType').value,drType:document.getElementById('drType').value.trim(),erType:document.getElementById('erType').value.trim(),customs:[]};defFields.forEach(id=>{o[id]=getNum(id);});return o;}
function applyDefSnapshot(s){if(!s)return;document.getElementById('flatFooted').checked=!!s.flatFooted;document.getElementById('fullDefense').checked=!!s.fullDefense;document.getElementById('blocking').checked=!!s.blocking;document.getElementById('cover').checked=!!s.cover;document.getElementById('touchAttack').checked=!!s.touchAttack;document.getElementById('perfectDefenseActive').checked=!!s.perfectDefenseActive;defFields.forEach(id=>{if(s[id]!==undefined)setNum(id,s[id]);});const n2=document.getElementById('other2Name');if(n2)n2.textContent=s.other2Name||'其他1';const n3=document.getElementById('other3Name');if(n3)n3.textContent=s.other3Name||'其他2';setNum('armorRanged',getNum('armorMelee'));setNum('shieldRanged',getNum('shieldMelee'));document.getElementById('damageAbsorbType').value=s.damageAbsorbType||'physical';document.getElementById('drType').value=s.drType||'';document.getElementById('erType').value=s.erType||'';traits=(Array.isArray(s.traits)&&s.traits.length)?s.traits.slice():[''];renderTraits();const t=s.totals||{};document.getElementById('tSpeed').textContent=t.speedPoolMelee||0;document.getElementById('tArmor').textContent=s.touchAttack?0:(t.armorPoolMelee||0);document.getElementById('tMagic').textContent=(t.magicPool||0)+(t.otherPool||0);document.getElementById('pools').textContent=  '合计 '+( (t.speedPoolMelee||0) + (s.touchAttack?0:(t.armorPoolMelee||0)) + (t.magicPool||0) + (t.otherPool||0) )  +(s.touchAttack?'（接触攻击，破甲池不计入）':'');const e=s.effective||{};document.getElementById('effBase').textContent=e.base!==s.base?('→'+e.base):'';document.getElementById('effDodge').textContent=e.dodge!==s.dodge?('→'+e.dodge):'';document.getElementById('effBlock').textContent=e.block!==s.block?('→'+e.block):'';document.getElementById('effSM').textContent=(e.shieldMelee!==s.shieldMelee)?('→'+e.shieldMelee):'';document.getElementById('effCover').textContent=e.cover!==undefined?('→'+e.cover):'';setNum('armorRanged',getNum('armorMelee'));setNum('shieldRanged',getNum('shieldMelee'));}
async function saveDef(){setNum('armorRanged',getNum('armorMelee'));setNum('shieldRanged',getNum('shieldMelee'));const s=await api('/api/defense/update',readDefForm());if(s)applyDefSnapshot(s);}

document.getElementById('btnDefReset').onclick=async()=>{if(!(await showConfirm('清空','清空防御预设？')))return;const s=await api('/api/defense/reset',{});if(s)applyDefSnapshot(s);};
['flatFooted','fullDefense','blocking','cover','touchAttack','perfectDefenseActive'].forEach(id=>document.getElementById(id).addEventListener('change',saveDef));
const defLabels={"base": "基础", "dodge": "闪避", "block": "格挡", "natural": "天生", "armorMelee": "盔甲", "armorRanged": "盔甲(远)", "shieldMelee": "盾牌", "shieldRanged": "盾牌(远)", "force": "力场", "deflection": "偏斜", "insight": "洞察", "coverBonus": "掩蔽", "other2": "其他2", "other3": "其他3", "resistSpeed": "抵高速", "resistAP": "抵破甲", "resistMagic": "抵破魔", "defenseBonusSuccess": "防御附加成功", "damageAbsorb": "伤害吸收", "drValue": "物理伤害减免", "erValue": "能量抗力", "perfectDefense": "完美防御"};defFields.forEach(id=>{  const el=document.getElementById(id);  if(!el)return;  el.onclick=()=>editNum(id,defLabels[id]||id,{onChange:saveDef});});
document.getElementById('damageAbsorbType').addEventListener('change',saveDef);
document.getElementById('drType').addEventListener('change',saveDef);document.getElementById('erType').addEventListener('change',saveDef);

document.addEventListener('keydown',e=>{const open=document.getElementById('modalMask').classList.contains('show');if(e.key==='Escape'&&open){e.preventDefault();finishModal(null);return;}if(e.key!=='Enter')return;if(open){e.preventDefault();const ok=document.getElementById('modalOk')||document.querySelector('#modalActions button:not(.secondary)');if(ok)ok.click();return;}if(e.target&&e.target.getAttribute('data-enter')==='applyDamage'){e.preventDefault();applyDamage();}});
document.getElementById('modalMask').addEventListener('click',e=>{if(e.target.id==='modalMask')finishModal(null);});

async function resolveAttack(){
  const body={
    ranged:false,
    attackDP:getNum('atkDP'),
    speed:getNum('atkSpeed'),
    armorPierce:getNum('atkAP'),
    magicPierce:getNum('atkMP'),
    explodeOn:parseInt(document.getElementById('atkExplode').value,10)||10,
    bonusSuccess:getNum('atkBonus'),
    isPhysical:document.getElementById('atkPhysical').checked,
    damageLimit:getNum('atkDmgLimit')
  };
  const r=await api('/api/defense/resolve',body);
  if(!r)return;
  if(r.roll)renderDiceResult(r.roll,'伤害');
  const box=document.getElementById('atkResult');
  const dmgType=document.getElementById('atkDmgType').value;
  const pl=r.pools||{};
  const spd=pl.speedPoolAfter!=null?pl.speedPoolAfter:0;
  const arm=pl.armorPoolAfter!=null?pl.armorPoolAfter:0;
  const mag=pl.magicPoolAfter!=null?pl.magicPoolAfter:0;
  if(r.miss){
    box.innerHTML='<div class="miss">未命中</div>'+
      '<div>'+(r.missReason||'')+'</div>'+
      '<div>有效防御 '+r.effectiveDefense+'（高速 '+spd+' / 破甲 '+arm+' / 破魔 '+mag+(r.perfectDefense?' · 完美'+r.perfectDefense:'')+'）</div>'+
      '<div>实际DP '+r.actualDP+
        (r.naturalSuccess!=null?' · 自然成功 '+r.naturalSuccess:'')+
        (r.defenseBonus?' · 防御附加 '+r.defenseBonus:'')+'</div>';
    return;
  }
  box.innerHTML=
    '<div class="big">最终伤害 '+r.finalDamage+'（'+dmgType+'）</div>'+
    '<div>有效防御 '+r.effectiveDefense+(r.touchAttack?'（接触攻击，破甲不计入）':'')+
      '（高速 '+spd+' / 破甲 '+arm+' / 破魔 '+mag+(r.perfectDefense?' · 完美'+r.perfectDefense:'')+'）</div>'+
    '<div>击破 高速'+r.effSpeed+'（抵'+r.resistSpeed+'）/ 破甲'+r.effArmorPierce+'（抵'+r.resistAP+'）/ 破魔'+r.effMagicPierce+'（抵'+r.resistMagic+'）</div>'+
    '<div>实际DP '+r.actualDP+' · 自然成功 '+r.naturalSuccess+' · 最终成功 '+r.finalSuccess+'</div>'+
    '<div>成功 '+r.finalSuccess+(r.defenseBonus?' − 防御附加'+r.defenseBonus:'')+
      ' → 伤害 '+r.rawDamage+
      (r.damageLimit>0?' → 上限后 '+r.afterLimit:'')+
      ' → 减免后 '+r.afterDR+' → 吸收后 '+r.afterAbsorb+'</div>'+
    '<div style="color:#93c5fd;margin-top:4px">请自行在血量区扣除 '+r.finalDamage+' 点'+dmgType+'伤害</div>';
}

document.getElementById('atkDP').onclick=()=>editNum('atkDP','攻击DP');
document.getElementById('atkSpeed').onclick=()=>editNum('atkSpeed','高速');
document.getElementById('atkAP').onclick=()=>editNum('atkAP','破甲');
document.getElementById('atkMP').onclick=()=>editNum('atkMP','破魔');
document.getElementById('atkBonus').onclick=()=>editNum('atkBonus','附加成功',{allowNeg:true});
document.getElementById('atkDmgLimit').onclick=()=>editNum('atkDmgLimit','伤害上限',{desc:'≤0 表示不限制'});
document.getElementById('btnAtkResolve').onclick=resolveAttack;


async function rollSave(kind){
  const map={fort:['saveFortDP','saveFortBonus','saveFortOut','强韧'],ref:['saveRefDP','saveRefBonus','saveRefOut','反射'],will:['saveWillDP','saveWillBonus','saveWillOut','意志']};
  const [dpId,bonusId,outId,name]=map[kind];
  const dp=getNum(dpId), bonus=getNum(bonusId);
  const r=await api('/api/dice/roll',{dp,explodeOn:10,bonus});
  if(!r)return;
  const el=document.getElementById(outId);
  if(r.criticalFailure){el.textContent='大失败';el.style.color='var(--danger)';}
  else{el.textContent=r.finalSuccess;el.style.color='#60a5fa';}
  el.title=r.summary;
  renderDiceResult(r,name+'豁免');
}
document.getElementById('btnSaveFort').onclick=()=>rollSave('fort');
document.getElementById('btnSaveRef').onclick=()=>rollSave('ref');
document.getElementById('btnSaveWill').onclick=()=>rollSave('will');
document.getElementById('saveFortDP').onclick=()=>editNum('saveFortDP','强韧 DP');
document.getElementById('saveFortBonus').onclick=()=>editNum('saveFortBonus','强韧附加成功',{allowNeg:true});
document.getElementById('saveRefDP').onclick=()=>editNum('saveRefDP','反射 DP');
document.getElementById('saveRefBonus').onclick=()=>editNum('saveRefBonus','反射附加成功',{allowNeg:true});
document.getElementById('saveWillDP').onclick=()=>editNum('saveWillDP','意志 DP');
document.getElementById('saveWillBonus').onclick=()=>editNum('saveWillBonus','意志附加成功',{allowNeg:true});


function bindEditLabels(){
  document.querySelectorAll('.edit-label').forEach(wrap=>{
    if(wrap.dataset.bound)return;
    wrap.dataset.bound='1';
    const startEdit=()=>{
      const textEl=wrap.querySelector('.edit-label-text');
      if(!textEl||wrap.querySelector('input.inline-edit'))return;
      const cur=textEl.textContent||'';
      const inp=document.createElement('input');
      inp.type='text';inp.className='inline-edit';inp.value=cur;
      inp.maxLength=12;
      textEl.style.display='none';
      wrap.insertBefore(inp,textEl);
      inp.focus();inp.select();
      const commit=()=>{
        let v=(inp.value||'').trim();
        if(!v)v=textEl.id==='other3Name'?'其他2':'其他1';
        textEl.textContent=v;
        textEl.style.display='';
        inp.remove();
        if(typeof saveDef==='function')saveDef();
      };
      inp.addEventListener('blur',commit);
      inp.addEventListener('keydown',e=>{
        if(e.key==='Enter'){e.preventDefault();inp.blur();}
        if(e.key==='Escape'){inp.value=cur;inp.blur();}
      });
    };
    wrap.addEventListener('click',e=>{
      e.preventDefault();e.stopPropagation();
      startEdit();
    });
  });
}
bindEditLabels();

(async()=>{const hs=await fetch('/api/state');renderHp(await hs.json());const ds=await fetch('/api/defense/state');applyDefSnapshot(await ds.json());})();
</script>
</body>
</html>
`
