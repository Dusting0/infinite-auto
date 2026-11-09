import type { DefensePreset } from "./types";

// 攻击表单（自定义伤害模块的单套攻击参数）。
// 之前定义在 App.tsx，因攻击预设需要被持久化，移到此处与存储层同处。
export type CombatDamageKind = "physical" | "energy" | "mixed";
export type WoundType = "冲击" | "严重" | "恶性";

export type AttackForm = {
  attackDP: number;
  speed: number;
  armorPierce: number;
  magicPierce: number;
  explodeOn: number;
  bonusSuccess: number;
  damageKind: CombatDamageKind;
  damageLimit: number;
  woundType: WoundType;
  magic: boolean;
  divine: boolean;
  physSubtype: string;
  energyType: string;
};

export const defaultAttackForm: AttackForm = {
  attackDP: 8,
  speed: 0,
  armorPierce: 0,
  magicPierce: 0,
  explodeOn: 10,
  bonusSuccess: 0,
  damageKind: "physical",
  damageLimit: 0,
  woundType: "严重",
  magic: false,
  divine: false,
  physSubtype: "",
  energyType: "pure",
};

// 一套攻击预设：名字 + 完整攻击表单快照。
export interface AttackPreset {
  id: string;
  name: string;
  form: AttackForm;
}

// 三豁免每行：DP、附加成功；out 是上次检定的展示值（顺带随预设记忆）。
export interface SaveRow {
  dp: number;
  bonus: number;
  out: string;
}

export type SavesState = { fort: SaveRow; ref: SaveRow; will: SaveRow };

export const defaultSaves: SavesState = {
  fort: { dp: 5, bonus: 0, out: "-" },
  ref: { dp: 5, bonus: 0, out: "-" },
  will: { dp: 5, bonus: 0, out: "-" },
};

export interface DiceState {
  dp: number;
  bonus: number;
  explodeOn: number;
}

export const defaultDice: DiceState = { dp: 5, bonus: 0, explodeOn: 10 };

// 全局预设：整套秘典页面的完整快照。
// 血量取四项原始值（上限/冲击/严重/恶性，完好由服务端推导）；
// 防御存 normalize 后的草稿，切换时经 /api/defense/update 推回服务端。
export interface GlobalPreset {
  id: string;
  name: string;
  hp: { max: number; b: number; l: number; a: number };
  defense?: DefensePreset;
  saves: SavesState;
  saveTemp: number;
  dice: DiceState;
  attacks: AttackPreset[];
  activeAttackId: string;
}

export interface PresetStore {
  list: GlobalPreset[];
  activeId: string;
}

const STORAGE_KEY = "bloodthorn.presets.v1";

export function uid() {
  return Math.random().toString(36).slice(2, 10) + Date.now().toString(36).slice(-4);
}

export function loadStore(): PresetStore | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as PresetStore;
    if (!parsed?.list?.length) return null;
    if (!parsed.list.some((preset) => preset.id === parsed.activeId)) parsed.activeId = parsed.list[0].id;
    return parsed;
  } catch {
    return null;
  }
}

export function saveStore(store: PresetStore) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(store));
  } catch {
    // localStorage 不可用（隐私模式等）时静默降级为仅内存态
  }
}
