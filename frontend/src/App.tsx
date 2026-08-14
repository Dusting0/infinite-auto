import { useEffect, useState, type ReactNode } from "react";
import { AlertTriangle, CheckCircle2, ChevronDown, CircleSlash, Coffee, Dices, EyeOff, HeartPulse, Minus, Moon, Plus, RotateCcw, Shield, Skull, Trash2, X } from "lucide-react";
import { getJson, postJson } from "./lib/api";
import { toInt } from "./lib/utils";
import type {
  AttackResult,
  DefensePreset,
  DefenseSnapshot,
  DiceResult,
  HpSnapshot,
  LongRestOptions,
} from "./lib/types";
import { Button } from "./components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "./components/ui/card";
import { Checkbox } from "./components/ui/checkbox";
import { Input } from "./components/ui/input";
import { Label } from "./components/ui/label";
import { Select } from "./components/ui/select";

const damageTypes = ["冲击", "严重", "恶性"] as const;
const combatDamageKinds = [
  { value: "physical", label: "物理伤害" },
  { value: "energy", label: "能量伤害" },
  { value: "mixed", label: "混合伤害" },
] as const;
const explodeOptions = [10, 9, 8];

const hpStatusMeta: Record<string, { icon: typeof HeartPulse; tone: string }> = {
  ok: { icon: CheckCircle2, tone: "text-emerald-300" },
  dazed: { icon: EyeOff, tone: "text-amber-300" },
  dead: { icon: Skull, tone: "text-rose-300" },
  over: { icon: AlertTriangle, tone: "text-amber-300" },
  unset: { icon: CircleSlash, tone: "text-muted-foreground" },
};

type CombatDamageKind = (typeof combatDamageKinds)[number]["value"];
type AttackForm = {
  attackDP: number;
  speed: number;
  armorPierce: number;
  magicPierce: number;
  explodeOn: number;
  bonusSuccess: number;
  damageKind: CombatDamageKind;
  damageLimit: number;
  woundType: (typeof damageTypes)[number];
};

const defaultAttackForm: AttackForm = {
  attackDP: 8,
  speed: 0,
  armorPierce: 0,
  magicPierce: 0,
  explodeOn: 10,
  bonusSuccess: 0,
  damageKind: "physical",
  damageLimit: 0,
  woundType: "严重",
};

const emptyHp: HpSnapshot = {
  max: 20,
  intact: 20,
  b: 0,
  l: 0,
  a: 0,
  total: 0,
  status: "加载中",
  statusKey: "unset",
  log: [],
};

const numberFields: Array<{ key: keyof DefensePreset; label: string; group: "speed" | "armor" | "magic" | "perfect" | "absorb" | "other" }> = [
  { key: "base", label: "基础", group: "speed" },
  { key: "dodge", label: "闪避", group: "speed" },
  { key: "block", label: "格挡", group: "speed" },
  { key: "natural", label: "天生", group: "armor" },
  { key: "armorMelee", label: "盔甲", group: "armor" },
  { key: "shieldMelee", label: "盾牌", group: "armor" },
  { key: "force", label: "力场", group: "magic" },
  { key: "deflection", label: "偏斜", group: "magic" },
  { key: "insight", label: "洞察", group: "magic" },
  { key: "coverBonus", label: "掩蔽", group: "magic" },
  { key: "other2", label: "其他1", group: "other" },
  { key: "other3", label: "其他2", group: "other" },
  { key: "perfectDefense", label: "完美防御", group: "perfect" },
  { key: "defenseBonusSuccess", label: "防御附加", group: "absorb" },
  { key: "damageAbsorb", label: "伤害吸收", group: "absorb" },
  { key: "drValue", label: "DR", group: "absorb" },
  { key: "erValue", label: "ER", group: "absorb" },
];

function asNumber(value: unknown) {
  return typeof value === "number" && Number.isFinite(value) ? value : 0;
}

function formatDiceFaces(result: DiceResult) {
  if (result.diceLog) return result.diceLog;
  if (result.batches?.length) return result.batches.map((batch) => `[${batch.join(",")}]`).join("+");
  if (result.dice?.length) return `[${result.dice.join(",")}]`;
  return "-";
}

function cloneDefense(snapshot: DefenseSnapshot): DefensePreset {
  const { totals: _totals, effective: _effective, ...preset } = snapshot;
  return {
    ...preset,
    armorRanged: preset.armorMelee,
    shieldRanged: preset.shieldMelee,
    traits: preset.traits?.length ? [...preset.traits] : [],
    customs: preset.customs ?? [],
    notes: preset.notes ?? "",
  };
}

function normalizeDefenseDraft(draft: DefensePreset): DefensePreset {
  return {
    ...draft,
    armorRanged: asNumber(draft.armorMelee),
    shieldRanged: asNumber(draft.shieldMelee),
    traits: (draft.traits ?? []).map((trait) => trait.trim()).filter(Boolean),
    customs: draft.customs ?? [],
    notes: draft.notes ?? "",
    damageAbsorbType: draft.damageAbsorbType || "physical",
    other2Name: draft.other2Name || "其他1",
    other3Name: draft.other3Name || "其他2",
  };
}

function computeDefensePreview(draft: DefensePreset) {
  const base = draft.flatFooted ? 0 : asNumber(draft.base) * (draft.fullDefense ? 2 : 1);
  const dodge = draft.flatFooted ? 0 : asNumber(draft.dodge);
  const block = draft.flatFooted || !draft.blocking ? 0 : asNumber(draft.block);
  const shield = draft.blocking ? asNumber(draft.shieldMelee) : 0;
  const cover = draft.cover ? asNumber(draft.coverBonus) : 0;
  let speed = base + dodge + block;
  let armor = shield + asNumber(draft.armorMelee) + asNumber(draft.natural);
  let magic = asNumber(draft.force) + asNumber(draft.deflection);
  let other = asNumber(draft.insight) + cover + asNumber(draft.other2) + asNumber(draft.other3);

  for (const custom of draft.customs ?? []) {
    let value = asNumber(custom.value);
    if (custom.lostOnFlatFooted && draft.flatFooted) value = 0;
    if (custom.onlyWhenBlocking && !draft.blocking) value = 0;
    if (custom.scope === "ranged") continue;
    if (custom.category === "speed") speed += value;
    else if (custom.category === "armor") armor += value;
    else if (custom.category === "magic") magic += value;
    else other += value;
  }

  return {
    speed,
    armor: draft.touchAttack ? 0 : armor,
    magic: magic + other,
    effective: { base, dodge, block },
  };
}

type NumberEditorProps = {
  label: string;
  value: number;
  onChange: (value: number) => void;
  min?: number;
  allowNegative?: boolean;
  triggerClassName?: string;
  compact?: boolean;
};

function NumberEditor({ label, value, onChange, min, allowNegative = false, triggerClassName, compact = false }: NumberEditorProps) {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState(String(value));
  const lowerBound = min ?? (allowNegative ? undefined : 0);

  useEffect(() => {
    if (open) setDraft(String(value));
  }, [open, value]);

  const adjust = (amount: number) => {
    const next = toInt(draft, value) + amount;
    setDraft(String(lowerBound === undefined ? next : Math.max(lowerBound, next)));
  };
  const commit = () => {
    const next = toInt(draft, value);
    onChange(lowerBound === undefined ? next : Math.max(lowerBound, next));
    setOpen(false);
  };

  return (
    <>
      <button type="button" className={triggerClassName} onClick={() => setOpen(true)} aria-haspopup="dialog" aria-expanded={open} aria-label={`编辑${label}`}>
        {compact ? <><span>{label}</span><strong>{value}</strong></> : <><span className="text-[11px] text-muted-foreground">{label}</span><strong className="mt-1 text-lg font-semibold tabular-nums">{value}</strong></>}
      </button>
      {open && (
        <div className="number-editor-backdrop" role="presentation" onMouseDown={() => setOpen(false)}>
          <section className="number-editor" role="dialog" aria-modal="true" aria-labelledby="number-editor-title" onMouseDown={(event) => event.stopPropagation()}>
            <div className="flex items-start justify-between gap-4">
              <div>
                <div className="eyebrow">数值调整</div>
                <h3 id="number-editor-title" className="font-display mt-1 text-2xl font-semibold text-foreground">{label}</h3>
              </div>
              <button type="button" className="icon-dismiss" onClick={() => setOpen(false)} aria-label="关闭"><X className="h-4 w-4" /></button>
            </div>
            <div className="number-stepper">
              <button type="button" onClick={() => adjust(-1)} aria-label={`${label}减一`}><Minus className="h-5 w-5" /></button>
              <input
                autoFocus
                inputMode="numeric"
                value={draft}
                onFocus={(event) => event.currentTarget.select()}
                onChange={(event) => setDraft(event.target.value)}
                onKeyDown={(event) => {
                  if (event.key === "Enter") commit();
                  if (event.key === "Escape") setOpen(false);
                }}
                aria-label={label}
              />
              <button type="button" onClick={() => adjust(1)} aria-label={`${label}加一`}><Plus className="h-5 w-5" /></button>
            </div>
            <div className="mt-6 grid grid-cols-2 gap-2">
              <Button variant="secondary" onClick={() => setOpen(false)}>取消</Button>
              <Button onClick={commit}>确认数值</Button>
            </div>
          </section>
        </div>
      )}
    </>
  );
}

function StatPill({ label, value, tone = "default", onChange, min, allowNegative = false }: { label: string; value: number | string; tone?: "default" | "good" | "warn" | "bad" | "blue"; onChange?: (value: number) => void; min?: number; allowNegative?: boolean }) {
  const toneClass = {
    default: "text-foreground",
    good: "text-emerald-300",
    warn: "text-amber-300",
    bad: "text-rose-300",
    blue: "text-sky-300",
  }[tone];
  const className = `stat-pill ${onChange ? "stat-pill-editable" : ""} ${toneClass}`;
  if (onChange && typeof value === "number") return <NumberEditor label={label} value={value} onChange={onChange} min={min} allowNegative={allowNegative} triggerClassName={className} />;
  return <div className={className}><span className="text-[11px] text-muted-foreground">{label}</span><strong className="mt-1 text-lg font-semibold tabular-nums">{value}</strong></div>;
}

function NumberInput({ label, value, onChange, min, allowNegative = false }: { label: string; value: number; onChange: (value: number) => void; min?: number; allowNegative?: boolean }) {
  return <NumberEditor label={label} value={value} onChange={onChange} min={min} allowNegative={allowNegative} triggerClassName="number-field-trigger" />;
}

function DiceFaceGrid({ result }: { result: DiceResult }) {
  const batches = result.batches?.length ? result.batches : [result.dice ?? []];
  return (
    <div className="mt-3 grid grid-cols-10 gap-1.5">
      {batches.flatMap((batch, batchIndex) =>
        batch.map((face, index) => {
          const success = result.chanceDie && batchIndex === 0 ? face === 10 : face >= 8;
          return (
            <span
              key={`${batchIndex}-${index}`}
              className={`flex aspect-square items-center justify-center rounded-md border text-xs font-semibold tabular-nums ${
                success ? "border-emerald-400/70 bg-emerald-500/15 text-emerald-200" : "border-border bg-secondary text-muted-foreground"
              } ${batchIndex > 0 ? "border-dashed" : ""}`}
            >
              {face}
            </span>
          );
        }),
      )}
    </div>
  );
}

function DiceResultView({ result, compact = false }: { result?: DiceResult; compact?: boolean }) {
  if (!result) {
    return <div className="rounded-md border border-dashed border-border bg-background/50 px-3 py-4 text-sm text-muted-foreground">暂无结果</div>;
  }
  return (
    <div className="rounded-md border border-border bg-background/70 p-3">
      <div className="flex items-start gap-3">
        <div className="w-1 self-stretch rounded-full bg-primary" />
        <div className="min-w-0 flex-1">
          <div className="text-xs text-muted-foreground">最终成功</div>
          <div className={`text-3xl font-bold tabular-nums ${result.criticalFailure ? "text-rose-300" : "text-sky-300"}`}>
            {result.criticalFailure ? "x" : result.finalSuccess}
          </div>
          <div className="mt-1 break-words text-xs text-muted-foreground">{result.criticalFailure ? "大失败 · " : ""}{result.summary}</div>
        </div>
      </div>
      {!compact && <DiceFaceGrid result={result} />}
      <div className="mt-2 break-words font-mono text-xs text-muted-foreground">{formatDiceFaces(result)}</div>
    </div>
  );
}

function Collapsible({ title, defaultOpen = false, actions, children }: { title: string; defaultOpen?: boolean; actions?: ReactNode; children: ReactNode }) {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <div className="rounded-md border border-border bg-background/40">
      <div className="flex items-center justify-between gap-2 px-3 py-2">
        <button type="button" className="flex flex-1 items-center gap-2 text-sm font-semibold text-foreground" onClick={() => setOpen((value) => !value)}>
          <span>{title}</span>
          <ChevronDown className={`h-4 w-4 text-muted-foreground transition-transform ${open ? "rotate-180" : ""}`} />
        </button>
        {actions && <div className="flex items-center gap-2">{actions}</div>}
      </div>
      {open && <div className="space-y-3 px-3 pb-3 pt-1">{children}</div>}
    </div>
  );
}

function HpPanel({ hp, onHp, onError, defenseDraft, onDefenseDraftChange, onDefense, onDice }: { hp: HpSnapshot; onHp: (hp: HpSnapshot) => void; onError: (error: string) => void; defenseDraft?: DefensePreset; onDefenseDraftChange: (draft: DefensePreset) => void; onDefense: (snapshot: DefenseSnapshot) => void; onDice: (result: DiceResult, source: string) => void }) {
  const [damageAmount, setDamageAmount] = useState(1);
  const [damageType, setDamageType] = useState<(typeof damageTypes)[number]>("严重");
  const [resetOpen, setResetOpen] = useState(false);
  const [longRestOpen, setLongRestOpen] = useState(false);

  const run = async (work: () => Promise<HpSnapshot>) => {
    try {
      onHp(await work());
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const setHpValue = (key: "max" | "b" | "l" | "a", value: number) => run(() => postJson<HpSnapshot>("/api/setvalues", {
    max: key === "max" ? value : hp.max,
    b: key === "b" ? value : hp.b,
    l: key === "l" ? value : hp.l,
    a: key === "a" ? value : hp.a,
  }));
  const applyDamage = () => run(() => postJson<HpSnapshot>("/api/damage", { amount: damageAmount, type: damageType }));
  const shortRest = () => run(() => postJson<HpSnapshot>("/api/shortrest"));
  const reset = () => { setResetOpen(false); run(() => postJson<HpSnapshot>("/api/reset")); };
  const applyLongRest = (mode: "L" | "A") => {
    setLongRestOpen(false);
    run(() => postJson<HpSnapshot>("/api/longrest", { mode }));
  };
  const longRest = async () => {
    try {
      const options = await getJson<LongRestOptions>("/api/longrest/options");
      if (options.CanClearL && options.CanConvertA) {
        setLongRestOpen(true);
      } else if (options.CanClearL) {
        applyLongRest("L");
      } else if (options.CanConvertA) {
        applyLongRest("A");
      } else {
        onError("无严重/恶性伤害可长休处理");
      }
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const max = Math.max(hp.max, 1);
  const statusMeta = hpStatusMeta[hp.statusKey] ?? hpStatusMeta.unset;
  const StatusIcon = statusMeta.icon;
  const parts = [
    { label: "完好", value: Math.max(hp.intact, 0), cls: "bg-emerald-400", tone: "text-emerald-300" },
    { label: "冲击", value: hp.b, cls: "bg-sky-400", tone: "text-sky-300" },
    { label: "严重", value: hp.l, cls: "bg-amber-400", tone: "text-amber-300" },
    { label: "恶性", value: hp.a, cls: "bg-rose-400", tone: "text-rose-300" },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><HeartPulse className="h-4 w-4" />血量</CardTitle>
        <div className={`flex items-center gap-1.5 text-sm font-semibold ${statusMeta.tone}`}><StatusIcon className="h-4 w-4" />{hp.status}</div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-5 gap-2">
          <StatPill label="上限" value={hp.max} tone="blue" min={1} onChange={(value) => setHpValue("max", value)} />
          <StatPill label="完好" value={hp.intact} tone={hp.intact < 0 ? "bad" : "good"} />
          <StatPill label="冲击" value={hp.b} tone="blue" onChange={(value) => setHpValue("b", value)} />
          <StatPill label="严重" value={hp.l} tone="warn" onChange={(value) => setHpValue("l", value)} />
          <StatPill label="恶性" value={hp.a} tone="bad" onChange={(value) => setHpValue("a", value)} />
        </div>

        <div className="flex h-8 overflow-hidden rounded-md border border-border bg-background">
          {parts.map((part) =>
            part.value > 0 ? (
              <div key={part.label} className={`${part.cls} flex items-center justify-center text-xs font-bold text-slate-950`} style={{ flex: part.value / max }} title={`${part.label} ${part.value}`}>
                {part.value >= 2 || part.value / max > 0.1 ? part.value : ""}
              </div>
            ) : null,
          )}
          {parts.every((part) => part.value <= 0) && <div className="flex flex-1 items-center justify-center text-xs text-muted-foreground">暂无生命分段</div>}
        </div>

        <div className="grid grid-cols-[1fr_1fr_auto] gap-2">
          <NumberInput label="受伤" value={damageAmount} min={1} onChange={setDamageAmount} />
          <label className="space-y-1.5">
            <Label>类型</Label>
            <Select value={damageType} onChange={(event) => setDamageType(event.target.value as (typeof damageTypes)[number])}>
              {damageTypes.map((type) => <option key={type}>{type}</option>)}
            </Select>
          </label>
          <div className="flex items-end"><Button onClick={applyDamage}>应用</Button></div>
        </div>

        <CustomDamageSection defenseDraft={defenseDraft} onDefenseDraftChange={onDefenseDraftChange} onDefense={onDefense} onDice={onDice} onError={onError} onResolved={(amount, type) => { setDamageAmount(amount); setDamageType(type); }} />

        <div className="grid grid-cols-3 gap-2">
          <Button variant="secondary" onClick={shortRest}><Coffee className="h-4 w-4" />短休</Button>
          <Button variant="secondary" onClick={longRest}><Moon className="h-4 w-4" />长休</Button>
          <Button variant="destructive" onClick={() => setResetOpen(true)}><RotateCcw className="h-4 w-4" />重置</Button>
        </div>

      </CardContent>
      {longRestOpen && (
        <div className="number-editor-backdrop" role="presentation" onMouseDown={() => setLongRestOpen(false)}>
          <section className="long-rest-dialog" role="dialog" aria-modal="true" aria-labelledby="long-rest-title" onMouseDown={(event) => event.stopPropagation()}>
            <div className="long-rest-dialog-icon"><Moon className="h-6 w-6" /></div>
            <div className="eyebrow">长休选择</div>
            <h3 id="long-rest-title" className="font-display mt-1 text-2xl font-semibold">选择恢复方式</h3>
            <div className="mt-5 grid grid-cols-2 gap-3">
              <button type="button" className="long-rest-option" onClick={() => applyLongRest("L")}>
                <span>恢复严重伤害</span>
                <strong>严重 → 完好</strong>
              </button>
              <button type="button" className="long-rest-option" onClick={() => applyLongRest("A")}>
                <span>压制恶性伤害</span>
                <strong>恶性 → 严重</strong>
              </button>
            </div>
            <div className="mt-5"><Button variant="secondary" className="w-full" onClick={() => setLongRestOpen(false)}>取消</Button></div>
          </section>
        </div>
      )}
      {resetOpen && (
        <div className="number-editor-backdrop" role="presentation" onMouseDown={() => setResetOpen(false)}>
          <section className="reset-dialog" role="dialog" aria-modal="true" aria-labelledby="reset-hp-title" onMouseDown={(event) => event.stopPropagation()}>
            <div className="reset-dialog-icon"><AlertTriangle className="h-6 w-6" /></div>
            <div className="eyebrow text-rose-300">危险操作</div>
            <h3 id="reset-hp-title" className="font-display mt-1 text-2xl font-semibold">重置血量？</h3>
            <p className="mt-2 text-sm leading-6 text-muted-foreground">这会清除冲击、严重和恶性伤害，并将生命状态恢复为满血。</p>
            <div className="mt-6 grid grid-cols-2 gap-2">
              <Button variant="secondary" onClick={() => setResetOpen(false)}>保留当前状态</Button>
              <Button variant="destructive" onClick={reset}><RotateCcw className="h-4 w-4" />确认重置</Button>
            </div>
          </section>
        </div>
      )}
    </Card>
  );
}

function DefensePanel({ snapshot, draft, onDraftChange, onDefense, onError }: { snapshot?: DefenseSnapshot; draft?: DefensePreset; onDraftChange: (draft: DefensePreset) => void; onDefense: (snapshot: DefenseSnapshot) => void; onError: (error: string) => void }) {
  const [traitText, setTraitText] = useState("");

  const saveDraft = async (nextDraft: DefensePreset) => {
    try {
      onDefense(await postJson<DefenseSnapshot>("/api/defense/update", normalizeDefenseDraft(nextDraft)));
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };
  const update = <K extends keyof DefensePreset>(key: K, value: DefensePreset[K]) => {
    if (!draft) return;
    const nextDraft = { ...draft, [key]: value };
    onDraftChange(nextDraft);
    void saveDraft(nextDraft);
  };
  const updateNumber = (key: keyof DefensePreset, value: number) => update(key, value as never);

  const reset = async () => {
    if (!window.confirm("清空防御预设？")) return;
    try {
      onDefense(await postJson<DefenseSnapshot>("/api/defense/reset"));
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const addTrait = () => {
    const next = traitText.trim();
    if (!next || !draft) return;
    const nextDraft = { ...draft, traits: [...(draft.traits ?? []), next] };
    onDraftChange(nextDraft);
    void saveDraft(nextDraft);
    setTraitText("");
  };

  if (!draft || !snapshot) {
    return <Card><CardContent className="text-sm text-muted-foreground">防御数据加载中</CardContent></Card>;
  }

  const preview = computeDefensePreview(draft);
  const groups = {
    speed: numberFields.filter((field) => field.group === "speed"),
    armor: numberFields.filter((field) => field.group === "armor"),
    magic: numberFields.filter((field) => field.group === "magic"),
    perfect: numberFields.filter((field) => field.group === "perfect"),
    absorb: numberFields.filter((field) => field.group === "absorb"),
    other: numberFields.filter((field) => field.group === "other"),
  };

  const renderFields = (fields: typeof numberFields) => (
    <div className="grid grid-cols-3 gap-2">
      {fields.map((field) => (
        <NumberInput key={String(field.key)} label={field.label} value={asNumber(draft[field.key])} allowNegative={field.key === "defenseBonusSuccess"} onChange={(value) => updateNumber(field.key, value)} />
      ))}
    </div>
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><Shield className="h-4 w-4" />防御预设</CardTitle>
        <Button size="sm" variant="ghost" onClick={reset}><Trash2 className="h-4 w-4" /></Button>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-3 gap-2">
          <StatPill label="高速防御" value={preview.speed} tone="blue" />
          <StatPill label="破甲防御" value={preview.armor} tone="warn" />
          <StatPill label="破魔防御" value={preview.magic} tone="good" />
        </div>
        <div className="grid grid-cols-3 gap-2 text-xs text-muted-foreground">
          <div>基础 {draft.base}{preview.effective.base !== draft.base ? ` -> ${preview.effective.base}` : ""}</div>
          <div>闪避 {draft.dodge}{preview.effective.dodge !== draft.dodge ? ` -> ${preview.effective.dodge}` : ""}</div>
          <div>格挡 {draft.block}{preview.effective.block !== draft.block ? ` -> ${preview.effective.block}` : ""}</div>
        </div>

        <div className="flex flex-wrap gap-2">
          <Checkbox label="格挡中" checked={draft.blocking} onChange={(event) => update("blocking", event.target.checked)} />
          <Checkbox label="全力防御" checked={draft.fullDefense} onChange={(event) => update("fullDefense", event.target.checked)} />
          <Checkbox label="掩蔽" checked={draft.cover} onChange={(event) => update("cover", event.target.checked)} />
        </div>

        <section className="space-y-2">
          <div className="section-label">抵高速</div>
          {renderFields(groups.speed)}
        </section>
        <section className="space-y-2">
          <div className="section-label">抵破甲</div>
          {renderFields(groups.armor)}
        </section>
        <section className="space-y-2">
          <div className="section-label">抵破魔 / 其他</div>
          {renderFields([...groups.magic, ...groups.other])}
        </section>
        <section className="space-y-2">
          <div className="section-label">完美防御</div>
          <div className="grid grid-cols-3 gap-2 items-end">
            <NumberInput label="完美防御" value={asNumber(draft.perfectDefense)} onChange={(value) => updateNumber("perfectDefense", value)} />
            <Checkbox label="生效" checked={draft.perfectDefenseActive} onChange={(event) => update("perfectDefenseActive", event.target.checked)} />
          </div>
        </section>
        <section className="space-y-2">
          <div className="section-label">伤害吸收 / 减免</div>
          {renderFields(groups.absorb)}
          <div className="grid grid-cols-3 gap-2">
            <label className="space-y-1.5">
              <Label>吸收类型</Label>
              <Select value={draft.damageAbsorbType || "physical"} onChange={(event) => update("damageAbsorbType", event.target.value)}>
                <option value="physical">物理伤害吸收</option>
                <option value="all">全伤害吸收</option>
              </Select>
            </label>
            <label className="space-y-1.5">
              <Label>DR 备注</Label>
              <Input value={draft.drType ?? ""} onChange={(event) => update("drType", event.target.value)} placeholder="神兵" />
            </label>
            <label className="space-y-1.5">
              <Label>ER 备注</Label>
              <Input value={draft.erType ?? ""} onChange={(event) => update("erType", event.target.value)} placeholder="能量类型" />
            </label>
          </div>
        </section>

        <section className="space-y-2">
          <div className="section-label">特性</div>
          <div className="flex gap-2">
            <Input value={traitText} onChange={(event) => setTraitText(event.target.value)} onKeyDown={(event) => { if (event.key === "Enter") addTrait(); }} placeholder="输入特性" />
            <Button variant="secondary" onClick={addTrait}>添加</Button>
          </div>
          <div className="flex flex-wrap gap-2">
            {(draft.traits ?? []).map((trait, index) => (
              <button key={`${trait}-${index}`} className="rounded-md border border-border bg-background px-2 py-1 text-xs text-foreground" onClick={() => {
                const nextDraft = { ...draft, traits: (draft.traits ?? []).filter((_, i) => i !== index) };
                onDraftChange(nextDraft);
                void saveDraft(nextDraft);
              }}>
                {trait} x
              </button>
            ))}
          </div>
        </section>
      </CardContent>
    </Card>
  );
}

function DicePanel({ result, onDice, onError }: { result?: DiceResult; onDice: (result: DiceResult, source: string) => void; onError: (error: string) => void }) {
  const [dp, setDp] = useState(5);
  const [bonus, setBonus] = useState(0);
  const [explodeOn, setExplodeOn] = useState(10);

  const roll = async (source = "掷骰", rollDp = dp, rollBonus = bonus, rollExplode = explodeOn) => {
    try {
      const next = await postJson<DiceResult>("/api/dice/roll", { dp: rollDp, explodeOn: rollExplode, bonus: rollBonus });
      onDice(next, source);
      return next;
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><Dices className="h-4 w-4" />掷骰</CardTitle>
        <Button size="sm" onClick={() => roll()}>投掷</Button>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-[1fr_1fr_96px] gap-2">
          <NumberInput label="DP" value={dp} allowNegative onChange={setDp} />
          <NumberInput label="附加成功" value={bonus} allowNegative onChange={setBonus} />
          <label className="space-y-1.5">
            <Label>加骰</Label>
            <Select value={explodeOn} onChange={(event) => setExplodeOn(toInt(event.target.value, 10))}>
              {explodeOptions.map((value) => <option key={value} value={value}>{value}</option>)}
            </Select>
          </label>
        </div>
        <DiceResultView result={result} />
      </CardContent>
    </Card>
  );
}

function SavesPanel({ onDice, onError }: { onDice: (result: DiceResult, source: string) => void; onError: (error: string) => void }) {
  const [saves, setSaves] = useState({ fort: { dp: 5, bonus: 0, out: "-" }, ref: { dp: 5, bonus: 0, out: "-" }, will: { dp: 5, bonus: 0, out: "-" } });

  const rollSave = async (kind: keyof typeof saves, label: string) => {
    try {
      const next = await postJson<DiceResult>("/api/dice/roll", { dp: saves[kind].dp, explodeOn: 10, bonus: saves[kind].bonus });
      onDice(next, `${label}豁免`);
      setSaves((prev) => ({ ...prev, [kind]: { ...prev[kind], out: next.criticalFailure ? "大失败" : String(next.finalSuccess) } }));
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><Shield className="h-4 w-4" />三豁免</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {([
          ["fort", "强韧"],
          ["ref", "反射"],
          ["will", "意志"],
        ] as const).map(([kind, label]) => (
          <div key={kind} className="grid grid-cols-[48px_1fr_1fr_auto_64px] items-end gap-2">
            <div className="pb-2 text-sm font-semibold text-sky-300">{label}</div>
            <NumberInput label="DP" value={saves[kind].dp} allowNegative onChange={(value) => setSaves((prev) => ({ ...prev, [kind]: { ...prev[kind], dp: value } }))} />
            <NumberInput label="附加" value={saves[kind].bonus} allowNegative onChange={(value) => setSaves((prev) => ({ ...prev, [kind]: { ...prev[kind], bonus: value } }))} />
            <Button variant="secondary" onClick={() => rollSave(kind, label)}>检定</Button>
            <div className="pb-2 text-right text-sm font-bold text-sky-300">{saves[kind].out}</div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}

function CustomDamageSection({ defenseDraft, onDefenseDraftChange, onDefense, onDice, onError, onResolved }: { defenseDraft?: DefensePreset; onDefenseDraftChange: (draft: DefensePreset) => void; onDefense: (snapshot: DefenseSnapshot) => void; onDice: (result: DiceResult, source: string) => void; onError: (error: string) => void; onResolved?: (amount: number, type: (typeof damageTypes)[number]) => void }) {
  const [form, setForm] = useState<AttackForm>(defaultAttackForm);
  const [result, setResult] = useState<AttackResult>();

  const update = <K extends keyof AttackForm>(key: K, value: AttackForm[K]) => setForm((prev) => ({ ...prev, [key]: value }));
  const resetAttack = () => {
    setForm(defaultAttackForm);
    setResult(undefined);
  };
  const saveDefenseDraft = async (nextDraft = defenseDraft) => {
    if (!nextDraft) return undefined;
    const normalized = normalizeDefenseDraft(nextDraft);
    onDefenseDraftChange(normalized);
    const next = await postJson<DefenseSnapshot>("/api/defense/update", normalized);
    onDefense(next);
    return next;
  };
  const updateDefenseFlag = async (key: "flatFooted" | "touchAttack", value: boolean) => {
    if (!defenseDraft) return;
    const nextDraft = { ...defenseDraft, [key]: value };
    onDefenseDraftChange(nextDraft);
    try {
      await saveDefenseDraft(nextDraft);
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };
  const resolve = async () => {
    try {
      await saveDefenseDraft();
      const next = await postJson<AttackResult>("/api/defense/resolve", {
        ranged: false,
        attackDP: form.attackDP,
        speed: form.speed,
        armorPierce: form.armorPierce,
        magicPierce: form.magicPierce,
        explodeOn: form.explodeOn,
        bonusSuccess: form.bonusSuccess,
        damageKind: form.damageKind,
        isPhysical: form.damageKind !== "energy",
        damageLimit: form.damageLimit,
      });
      setResult(next);
      if (next.roll) onDice(next.roll, "伤害");
      if (!next.miss && onResolved) onResolved(next.finalDamage, form.woundType);
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const pools = result?.pools;

  return (
    <Collapsible title="自定义伤害" actions={<><Button size="sm" variant="secondary" onClick={resetAttack}><RotateCcw className="h-4 w-4" />重置</Button><Button size="sm" onClick={resolve}>结算</Button></>}>
        <div className="grid grid-cols-3 gap-2">
          <NumberInput label="攻击 DP" value={form.attackDP} allowNegative onChange={(value) => update("attackDP", value)} />
          <NumberInput label="高速" value={form.speed} onChange={(value) => update("speed", value)} />
          <NumberInput label="破甲" value={form.armorPierce} onChange={(value) => update("armorPierce", value)} />
          <NumberInput label="破魔" value={form.magicPierce} onChange={(value) => update("magicPierce", value)} />
          <NumberInput label="附加成功" value={form.bonusSuccess} allowNegative onChange={(value) => update("bonusSuccess", value)} />
          <NumberInput label="伤害上限" value={form.damageLimit} onChange={(value) => update("damageLimit", value)} />
        </div>
        <div className="grid grid-cols-3 gap-2">
          <label className="space-y-1.5">
            <Label>加骰</Label>
            <div className="relative">
              <Select className="combat-select" value={form.explodeOn} onChange={(event) => update("explodeOn", toInt(event.target.value, 10))}>
                {explodeOptions.map((value) => <option key={value} value={value}>{value}</option>)}
              </Select>
              <ChevronDown className="combat-select-icon" />
            </div>
          </label>
          <label className="space-y-1.5">
            <Label>血量伤害</Label>
            <div className="relative">
              <Select className="combat-select combat-select-danger" value={form.woundType} onChange={(event) => update("woundType", event.target.value as AttackForm["woundType"])}>
                {damageTypes.map((type) => <option key={type}>{type}</option>)}
              </Select>
              <ChevronDown className="combat-select-icon" />
            </div>
          </label>
          <label className="space-y-1.5">
            <Label>伤害属性</Label>
            <div className="relative">
              <Select className="combat-select" value={form.damageKind} onChange={(event) => update("damageKind", event.target.value as CombatDamageKind)}>
                {combatDamageKinds.map((kind) => <option key={kind.value} value={kind.value}>{kind.label}</option>)}
              </Select>
              <ChevronDown className="combat-select-icon" />
            </div>
          </label>
        </div>
        <div className="flex flex-wrap gap-2">
          <Checkbox label="措手不及" checked={!!defenseDraft?.flatFooted} disabled={!defenseDraft} onChange={(event) => void updateDefenseFlag("flatFooted", event.target.checked)} />
          <Checkbox label="接触攻击" checked={!!defenseDraft?.touchAttack} disabled={!defenseDraft} onChange={(event) => void updateDefenseFlag("touchAttack", event.target.checked)} />
        </div>
        {result ? (
          <div className="rounded-md border border-border bg-background/70 p-3 text-sm">
            {result.miss ? (
              <div className="space-y-1">
                <div className="text-lg font-bold text-rose-300">未命中</div>
                <div className="text-muted-foreground">{result.missReason}</div>
              </div>
            ) : (
              <div className="space-y-1">
                <div className="text-xl font-bold text-sky-300">最终伤害 {result.finalDamage}（{form.woundType}）</div>
                <div className="text-muted-foreground">请自行在血量区扣除 {result.finalDamage} 点{form.woundType}伤害</div>
              </div>
            )}
            <div className="mt-3 grid grid-cols-3 gap-2 text-xs text-muted-foreground">
              <div>有效防御 {result.effectiveDefense}</div>
              <div>实际 DP {result.actualDP}</div>
              <div>防御附加 {result.defenseBonus}</div>
              <div>高速 {pools?.speedPoolAfter ?? 0}</div>
              <div>破甲 {pools?.armorPoolAfter ?? 0}</div>
              <div>破魔 {pools?.magicPoolAfter ?? 0}</div>
            </div>
            {!result.miss && (
              <div className="mt-2 text-xs text-muted-foreground">
                {`成功 ${result.finalSuccess} -> 伤害 ${result.rawDamage}${result.damageLimit > 0 ? ` -> 上限后 ${result.afterLimit}` : ""} -> 减免后 ${result.afterDR} -> 吸收后 ${result.afterAbsorb}`}
              </div>
            )}
          </div>
        ) : <div className="rounded-md border border-dashed border-border bg-background/50 px-3 py-4 text-sm text-muted-foreground">暂无结算</div>}
    </Collapsible>
  );
}

export default function App() {
  const [hp, setHp] = useState<HpSnapshot>(emptyHp);
  const [defense, setDefense] = useState<DefenseSnapshot>();
  const [defenseDraft, setDefenseDraft] = useState<DefensePreset>();
  const [error, setError] = useState("");
  const [diceLog, setDiceLog] = useState<Array<{ source: string; result: DiceResult }>>([]);
  const [latestDiceResult, setLatestDiceResult] = useState<DiceResult>();

  useEffect(() => {
    void (async () => {
      try {
        const [hpState, defenseState] = await Promise.all([
          getJson<HpSnapshot>("/api/state"),
          getJson<DefenseSnapshot>("/api/defense/state"),
        ]);
        setHp(hpState);
        setDefense(defenseState);
        setDefenseDraft(cloneDefense(defenseState));
      } catch (loadError) {
        setError(loadError instanceof Error ? loadError.message : String(loadError));
      }
    })();
  }, []);

  const updateDefenseSnapshot = (snapshot: DefenseSnapshot) => {
    setDefense(snapshot);
    setDefenseDraft(cloneDefense(snapshot));
  };

  const pushDice = (result: DiceResult, source: string) => {
    setLatestDiceResult(result);
    setDiceLog((prev) => [{ result, source }, ...prev].slice(0, 30));
  };

  return (
    <main className="min-h-screen bg-background text-foreground">
      <div className="mx-auto max-w-[1560px] px-4 py-4 lg:px-6 lg:py-6">
        <div className="min-w-0 content-wrap">
        {error && (
          <div className="mb-4 flex items-center justify-between rounded-lg border border-rose-500/40 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
            <span>{error}</span>
            <Button size="sm" variant="ghost" onClick={() => setError("")}>关闭</Button>
          </div>
        )}

        <div className="main-grid">
          <div className="space-y-4">
            <section id="vitality" className="scroll-mt-5"><HpPanel hp={hp} onHp={setHp} onError={setError} defenseDraft={defenseDraft} onDefenseDraftChange={setDefenseDraft} onDefense={updateDefenseSnapshot} onDice={pushDice} /></section>
            <section className="scroll-mt-5"><SavesPanel onDice={pushDice} onError={setError} /></section>
          </div>
          <section id="defense" className="scroll-mt-5"><DefensePanel snapshot={defense} draft={defenseDraft} onDraftChange={setDefenseDraft} onDefense={updateDefenseSnapshot} onError={setError} /></section>
          <div className="space-y-4">
            <section id="dice" className="scroll-mt-5"><DicePanel result={latestDiceResult} onDice={pushDice} onError={setError} /></section>
            <Card>
              <CardHeader><CardTitle>骰子记录</CardTitle></CardHeader>
              <CardContent>
                <div className="max-h-[420px] space-y-2 overflow-y-auto pr-1">
                {diceLog.length ? diceLog.map((entry, index) => (
                  <div key={`${entry.source}-${index}`} className="rounded-md border border-border bg-background/60 p-2">
                    <div className="mb-1 text-xs text-sky-300">[{entry.source}] 成功 {entry.result.criticalFailure ? "大失败" : entry.result.finalSuccess}</div>
                    <div className="break-words font-mono text-xs text-muted-foreground">{formatDiceFaces(entry.result)}</div>
                  </div>
                )) : <div className="text-sm text-muted-foreground">暂无记录</div>}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
        </div>
        <footer className="mt-6 text-center text-xs text-muted-foreground">
          Powered by Dusting
        </footer>
      </div>
    </main>
  );
}
