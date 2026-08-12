import { useEffect, useMemo, useState } from "react";
import { Activity, Dices, HeartPulse, RotateCcw, Save, Shield, Swords, Trash2 } from "lucide-react";
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
const explodeOptions = [10, 9, 8];

const emptyHp: HpSnapshot = {
  max: 20,
  intact: 20,
  b: 0,
  l: 0,
  a: 0,
  total: 0,
  status: "加载中",
  log: [],
};

const numberFields: Array<{ key: keyof DefensePreset; label: string; group: "speed" | "armor" | "magic" | "damage" | "other" }> = [
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
  { key: "perfectDefense", label: "完美防御", group: "damage" },
  { key: "defenseBonusSuccess", label: "防御附加", group: "damage" },
  { key: "damageAbsorb", label: "伤害吸收", group: "damage" },
  { key: "drValue", label: "DR", group: "damage" },
  { key: "erValue", label: "ER", group: "damage" },
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

function StatPill({ label, value, tone = "default" }: { label: string; value: number | string; tone?: "default" | "good" | "warn" | "bad" | "blue" }) {
  const toneClass = {
    default: "text-foreground",
    good: "text-emerald-300",
    warn: "text-amber-300",
    bad: "text-rose-300",
    blue: "text-sky-300",
  }[tone];
  return (
    <div className="rounded-md border border-border bg-background/70 px-3 py-2">
      <div className="text-[11px] text-muted-foreground">{label}</div>
      <div className={`mt-1 text-lg font-semibold tabular-nums ${toneClass}`}>{value}</div>
    </div>
  );
}

function NumberInput({ label, value, onChange, min, allowNegative = false }: { label: string; value: number; onChange: (value: number) => void; min?: number; allowNegative?: boolean }) {
  return (
    <label className="space-y-1.5">
      <Label>{label}</Label>
      <Input
        type="number"
        value={value}
        min={min ?? (allowNegative ? undefined : 0)}
        step={1}
        onChange={(event) => onChange(toInt(event.target.value))}
      />
    </label>
  );
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

function HpPanel({ hp, onHp, onError }: { hp: HpSnapshot; onHp: (hp: HpSnapshot) => void; onError: (error: string) => void }) {
  const [damageAmount, setDamageAmount] = useState(1);
  const [damageType, setDamageType] = useState<(typeof damageTypes)[number]>("严重");
  const [draft, setDraft] = useState({ max: hp.max, b: hp.b, l: hp.l, a: hp.a });

  useEffect(() => {
    setDraft({ max: hp.max, b: hp.b, l: hp.l, a: hp.a });
  }, [hp.max, hp.b, hp.l, hp.a]);

  const run = async (work: () => Promise<HpSnapshot>) => {
    try {
      onHp(await work());
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const saveValues = () => run(() => postJson<HpSnapshot>("/api/setvalues", draft));
  const applyDamage = () => run(() => postJson<HpSnapshot>("/api/damage", { amount: damageAmount, type: damageType }));
  const shortRest = () => run(() => postJson<HpSnapshot>("/api/shortrest"));
  const reset = () => {
    if (window.confirm("重置为满血？")) run(() => postJson<HpSnapshot>("/api/reset"));
  };
  const longRest = async () => {
    try {
      const options = await getJson<LongRestOptions>("/api/longrest/options");
      let mode = "";
      if (options.CanClearL && options.CanConvertA) {
        mode = window.confirm("长休选择：确定=严重恢复为完好，取消=恶性转严重") ? "L" : "A";
      } else if (options.CanClearL) {
        mode = "L";
      } else if (options.CanConvertA) {
        mode = "A";
      } else {
        onError("无严重/恶性伤害可长休处理");
        return;
      }
      onHp(await postJson<HpSnapshot>("/api/longrest", { mode }));
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const max = Math.max(hp.max, 1);
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
        <div className="text-sm font-semibold text-foreground">{hp.status}</div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-5 gap-2">
          <StatPill label="上限" value={hp.max} tone="blue" />
          <StatPill label="完好" value={hp.intact} tone={hp.intact < 0 ? "bad" : "good"} />
          <StatPill label="冲击" value={hp.b} tone="blue" />
          <StatPill label="严重" value={hp.l} tone="warn" />
          <StatPill label="恶性" value={hp.a} tone="bad" />
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

        <div className="grid grid-cols-4 gap-2">
          <NumberInput label="上限" value={draft.max} min={1} onChange={(value) => setDraft((prev) => ({ ...prev, max: value }))} />
          <NumberInput label="冲击" value={draft.b} onChange={(value) => setDraft((prev) => ({ ...prev, b: value }))} />
          <NumberInput label="严重" value={draft.l} onChange={(value) => setDraft((prev) => ({ ...prev, l: value }))} />
          <NumberInput label="恶性" value={draft.a} onChange={(value) => setDraft((prev) => ({ ...prev, a: value }))} />
        </div>
        <Button variant="secondary" className="w-full" onClick={saveValues}><Save className="h-4 w-4" />保存血量</Button>

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

        <div className="grid grid-cols-3 gap-2">
          <Button variant="secondary" onClick={shortRest}>短休</Button>
          <Button variant="secondary" onClick={longRest}>长休</Button>
          <Button variant="destructive" onClick={reset}><RotateCcw className="h-4 w-4" />重置</Button>
        </div>

        <div className="max-h-36 overflow-auto rounded-md border border-border bg-background/70 p-2 text-xs text-muted-foreground">
          {hp.log?.length ? hp.log.slice().reverse().map((line, index) => <div key={`${line}-${index}`} className="border-b border-border/60 py-1 last:border-0">{line}</div>) : <div>暂无记录</div>}
        </div>
      </CardContent>
    </Card>
  );
}

function DefensePanel({ snapshot, onDefense, onError }: { snapshot?: DefenseSnapshot; onDefense: (snapshot: DefenseSnapshot) => void; onError: (error: string) => void }) {
  const [draft, setDraft] = useState<DefensePreset | undefined>(snapshot ? cloneDefense(snapshot) : undefined);
  const [traitText, setTraitText] = useState("");

  useEffect(() => {
    if (snapshot) setDraft(cloneDefense(snapshot));
  }, [snapshot]);

  const update = <K extends keyof DefensePreset>(key: K, value: DefensePreset[K]) => setDraft((prev) => (prev ? { ...prev, [key]: value } : prev));
  const updateNumber = (key: keyof DefensePreset, value: number) => update(key, value as never);

  const save = async () => {
    if (!draft) return;
    const body: DefensePreset = {
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
    try {
      onDefense(await postJson<DefenseSnapshot>("/api/defense/update", body));
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

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
    update("traits", [...(draft.traits ?? []), next]);
    setTraitText("");
  };

  if (!draft || !snapshot) {
    return <Card><CardContent className="text-sm text-muted-foreground">防御数据加载中</CardContent></Card>;
  }

  const totals = snapshot.totals;
  const effective = snapshot.effective;
  const groups = {
    speed: numberFields.filter((field) => field.group === "speed"),
    armor: numberFields.filter((field) => field.group === "armor"),
    magic: numberFields.filter((field) => field.group === "magic"),
    damage: numberFields.filter((field) => field.group === "damage"),
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
        <div className="flex gap-2">
          <Button size="sm" variant="secondary" onClick={save}><Save className="h-4 w-4" />保存</Button>
          <Button size="sm" variant="ghost" onClick={reset}><Trash2 className="h-4 w-4" /></Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-3 gap-2">
          <StatPill label="高速池" value={totals.speedPoolMelee} tone="blue" />
          <StatPill label="破甲池" value={draft.touchAttack ? 0 : totals.armorPoolMelee} tone="warn" />
          <StatPill label="破魔/其他" value={totals.magicPool + totals.otherPool} tone="good" />
        </div>
        <div className="grid grid-cols-3 gap-2 text-xs text-muted-foreground">
          <div>基础 {draft.base}{effective.base !== draft.base ? ` -> ${effective.base}` : ""}</div>
          <div>闪避 {draft.dodge}{effective.dodge !== draft.dodge ? ` -> ${effective.dodge}` : ""}</div>
          <div>格挡 {draft.block}{effective.block !== draft.block ? ` -> ${effective.block}` : ""}</div>
        </div>

        <div className="flex flex-wrap gap-2">
          <Checkbox label="格挡中" checked={draft.blocking} onChange={(event) => update("blocking", event.target.checked)} />
          <Checkbox label="措手不及" checked={draft.flatFooted} onChange={(event) => update("flatFooted", event.target.checked)} />
          <Checkbox label="全力防御" checked={draft.fullDefense} onChange={(event) => update("fullDefense", event.target.checked)} />
          <Checkbox label="掩蔽" checked={draft.cover} onChange={(event) => update("cover", event.target.checked)} />
          <Checkbox label="接触攻击" checked={draft.touchAttack} onChange={(event) => update("touchAttack", event.target.checked)} />
          <Checkbox label="完美防御生效" checked={draft.perfectDefenseActive} onChange={(event) => update("perfectDefenseActive", event.target.checked)} />
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
          <div className="section-label">伤害相关</div>
          {renderFields(groups.damage)}
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
              <button key={`${trait}-${index}`} className="rounded-md border border-border bg-background px-2 py-1 text-xs text-foreground" onClick={() => update("traits", draft.traits.filter((_, i) => i !== index))}>
                {trait} x
              </button>
            ))}
          </div>
        </section>
      </CardContent>
    </Card>
  );
}

function DicePanel({ onDice, onError }: { onDice: (result: DiceResult, source: string) => void; onError: (error: string) => void }) {
  const [dp, setDp] = useState(5);
  const [bonus, setBonus] = useState(0);
  const [explodeOn, setExplodeOn] = useState(10);
  const [result, setResult] = useState<DiceResult>();
  const [saves, setSaves] = useState({ fort: { dp: 5, bonus: 0, out: "-" }, ref: { dp: 5, bonus: 0, out: "-" }, will: { dp: 5, bonus: 0, out: "-" } });

  const roll = async (source = "掷骰", rollDp = dp, rollBonus = bonus, rollExplode = explodeOn) => {
    try {
      const next = await postJson<DiceResult>("/api/dice/roll", { dp: rollDp, explodeOn: rollExplode, bonus: rollBonus });
      setResult(next);
      onDice(next, source);
      return next;
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const rollSave = async (kind: keyof typeof saves, label: string) => {
    const next = await roll(`${label}豁免`, saves[kind].dp, saves[kind].bonus, 10);
    if (!next) return;
    setSaves((prev) => ({ ...prev, [kind]: { ...prev[kind], out: next.criticalFailure ? "大失败" : String(next.finalSuccess) } }));
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><Dices className="h-4 w-4" />掷骰 / 豁免</CardTitle>
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
        <div className="space-y-2">
          <div className="section-label">三豁免</div>
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
        </div>
      </CardContent>
    </Card>
  );
}

function AttackPanel({ onDice, onError }: { onDice: (result: DiceResult, source: string) => void; onError: (error: string) => void }) {
  const [form, setForm] = useState({ attackDP: 8, speed: 0, armorPierce: 0, magicPierce: 0, explodeOn: 10, bonusSuccess: 0, isPhysical: true, damageLimit: 0, damageType: "严重" });
  const [result, setResult] = useState<AttackResult>();

  const update = (key: keyof typeof form, value: number | boolean | string) => setForm((prev) => ({ ...prev, [key]: value }));
  const resolve = async () => {
    try {
      const next = await postJson<AttackResult>("/api/defense/resolve", {
        ranged: false,
        attackDP: form.attackDP,
        speed: form.speed,
        armorPierce: form.armorPierce,
        magicPierce: form.magicPierce,
        explodeOn: form.explodeOn,
        bonusSuccess: form.bonusSuccess,
        isPhysical: form.isPhysical,
        damageLimit: form.damageLimit,
      });
      setResult(next);
      if (next.roll) onDice(next.roll, "伤害");
    } catch (error) {
      onError(error instanceof Error ? error.message : String(error));
    }
  };

  const pools = result?.pools;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><Swords className="h-4 w-4" />攻击结算</CardTitle>
        <Button size="sm" onClick={resolve}>结算</Button>
      </CardHeader>
      <CardContent className="space-y-4">
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
            <Select value={form.explodeOn} onChange={(event) => update("explodeOn", toInt(event.target.value, 10))}>
              {explodeOptions.map((value) => <option key={value} value={value}>{value}</option>)}
            </Select>
          </label>
          <label className="space-y-1.5">
            <Label>伤害类型</Label>
            <Select value={form.damageType} onChange={(event) => update("damageType", event.target.value)}>
              {damageTypes.map((type) => <option key={type}>{type}</option>)}
            </Select>
          </label>
          <div className="flex items-end"><Checkbox label="物理伤害" checked={form.isPhysical} onChange={(event) => update("isPhysical", event.target.checked)} /></div>
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
                <div className="text-xl font-bold text-sky-300">最终伤害 {result.finalDamage}（{form.damageType}）</div>
                <div className="text-muted-foreground">请自行在血量区扣除 {result.finalDamage} 点{form.damageType}伤害</div>
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
      </CardContent>
    </Card>
  );
}

export default function App() {
  const [hp, setHp] = useState<HpSnapshot>(emptyHp);
  const [defense, setDefense] = useState<DefenseSnapshot>();
  const [error, setError] = useState("");
  const [diceLog, setDiceLog] = useState<Array<{ source: string; result: DiceResult }>>([]);

  useEffect(() => {
    void (async () => {
      try {
        const [hpState, defenseState] = await Promise.all([
          getJson<HpSnapshot>("/api/state"),
          getJson<DefenseSnapshot>("/api/defense/state"),
        ]);
        setHp(hpState);
        setDefense(defenseState);
      } catch (loadError) {
        setError(loadError instanceof Error ? loadError.message : String(loadError));
      }
    })();
  }, []);

  const pushDice = (result: DiceResult, source: string) => {
    setDiceLog((prev) => [{ result, source }, ...prev].slice(0, 12));
  };

  const defenseTotal = useMemo(() => {
    if (!defense) return 0;
    return defense.totals.speedPoolMelee + (defense.touchAttack ? 0 : defense.totals.armorPoolMelee) + defense.totals.magicPool + defense.totals.otherPool;
  }, [defense]);

  return (
    <main className="min-h-screen bg-background text-foreground">
      <div className="mx-auto max-w-[1440px] px-4 py-4">
        <header className="mb-4 flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-2xl font-semibold tracking-normal text-foreground">无限规则计算器</h1>
            <p className="mt-1 text-sm text-muted-foreground">血量 · 掷骰 · 防御预设 · 攻击结算</p>
          </div>
          <div className="flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-2 text-sm text-muted-foreground">
            <Activity className="h-4 w-4 text-emerald-300" />
            防御合计 <span className="font-semibold text-sky-300">{defenseTotal}</span>
          </div>
        </header>

        {error && (
          <div className="mb-4 flex items-center justify-between rounded-lg border border-rose-500/40 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
            <span>{error}</span>
            <Button size="sm" variant="ghost" onClick={() => setError("")}>关闭</Button>
          </div>
        )}

        <div className="grid gap-4 xl:grid-cols-[1fr_1.1fr_0.95fr]">
          <div className="space-y-4">
            <HpPanel hp={hp} onHp={setHp} onError={setError} />
            <AttackPanel onDice={pushDice} onError={setError} />
          </div>
          <DefensePanel snapshot={defense} onDefense={setDefense} onError={setError} />
          <div className="space-y-4">
            <DicePanel onDice={pushDice} onError={setError} />
            <Card>
              <CardHeader><CardTitle>骰子记录</CardTitle></CardHeader>
              <CardContent className="space-y-2">
                {diceLog.length ? diceLog.map((entry, index) => (
                  <div key={`${entry.source}-${index}`} className="rounded-md border border-border bg-background/60 p-2">
                    <div className="mb-1 text-xs text-sky-300">[{entry.source}] 成功 {entry.result.criticalFailure ? "大失败" : entry.result.finalSuccess}</div>
                    <div className="break-words font-mono text-xs text-muted-foreground">{formatDiceFaces(entry.result)}</div>
                  </div>
                )) : <div className="text-sm text-muted-foreground">暂无记录</div>}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </main>
  );
}
