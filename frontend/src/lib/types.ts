export interface HpSnapshot {
  max: number;
  intact: number;
  b: number;
  l: number;
  a: number;
  total: number;
  status: string;
  log: string[];
}

export interface LongRestOptions {
  CanClearL: boolean;
  CanConvertA: boolean;
}

export type DefenseCategory = "speed" | "armor" | "magic" | "other";

export interface CustomComp {
  name: string;
  value: number;
  category: DefenseCategory;
  scope: "both" | "melee" | "ranged" | string;
  onlyWhenBlocking: boolean;
  lostOnFlatFooted: boolean;
}

export interface DefenseTotals {
  melee: number;
  ranged: number;
  speedPoolMelee: number;
  speedPoolRanged: number;
  armorPoolMelee: number;
  armorPoolRanged: number;
  magicPool: number;
  otherPool: number;
}

export interface EffectiveDefense {
  base: number;
  dodge: number;
  block: number;
  shieldMelee: number;
  shieldRanged: number;
  cover: number;
}

export interface DefensePreset {
  name: string;
  flatFooted: boolean;
  fullDefense: boolean;
  blocking: boolean;
  cover: boolean;
  touchAttack: boolean;
  base: number;
  dodge: number;
  block: number;
  natural: number;
  armorMelee: number;
  armorRanged: number;
  shieldMelee: number;
  shieldRanged: number;
  force: number;
  deflection: number;
  insight: number;
  coverBonus: number;
  other2: number;
  other3: number;
  other2Name: string;
  other3Name: string;
  perfectDefense: number;
  perfectDefenseActive: boolean;
  resistSpeed: number;
  resistAP: number;
  resistMagic: number;
  defenseBonusSuccess: number;
  damageAbsorb: number;
  damageAbsorbType: "physical" | "all" | string;
  drValue: number;
  drType: string;
  erValue: number;
  erType: string;
  customs: CustomComp[];
  notes: string;
  traits: string[];
}

export interface DefenseSnapshot extends DefensePreset {
  totals: DefenseTotals;
  effective: EffectiveDefense;
}

export interface DiceResult {
  dp: number;
  explodeOn: number;
  bonusSuccess: number;
  chanceDie: boolean;
  dice: number[];
  explosions: number[];
  batches: number[][];
  naturalSuccess: number;
  finalSuccess: number;
  bonusApplied: boolean;
  criticalFailure: boolean;
  summary: string;
  diceLog: string;
}

export interface AttackInput {
  ranged: boolean;
  attackDP: number;
  speed: number;
  armorPierce: number;
  magicPierce: number;
  explodeOn: number;
  bonusSuccess: number;
  isPhysical: boolean;
  damageKind: "physical" | "energy" | "mixed" | string;
  damageLimit: number;
}

export interface PoolBreakdown {
  blockBefore: number;
  blockAfter: number;
  dodgeBefore: number;
  dodgeAfter: number;
  baseBefore: number;
  baseAfter: number;
  shieldBefore: number;
  shieldAfter: number;
  armorBefore: number;
  armorAfter: number;
  naturalBefore: number;
  naturalAfter: number;
  forceBefore: number;
  forceAfter: number;
  deflectionBefore: number;
  deflectionAfter: number;
  insightBefore: number;
  insightAfter: number;
  coverBefore: number;
  coverAfter: number;
  other2Before: number;
  other2After: number;
  other3Before: number;
  other3After: number;
  speedPoolBefore: number;
  speedPoolAfter: number;
  armorPoolBefore: number;
  armorPoolAfter: number;
  magicPoolBefore: number;
  magicPoolAfter: number;
}

export interface AttackResult {
  miss: boolean;
  missReason?: string;
  effectiveDefense: number;
  actualDP: number;
  effSpeed: number;
  effArmorPierce: number;
  effMagicPierce: number;
  resistSpeed: number;
  resistAP: number;
  resistMagic: number;
  touchAttack: boolean;
  perfectDefense: number;
  perfectDefenseActive: boolean;
  pools: PoolBreakdown;
  roll?: DiceResult;
  naturalSuccess: number;
  finalSuccess: number;
  defenseBonus: number;
  rawDamage: number;
  afterLimit: number;
  damageLimit: number;
  afterDR: number;
  afterAbsorb: number;
  finalDamage: number;
  drValue: number;
  drType: string;
  erValue: number;
  erType: string;
  absorbApplied: number;
  absorbType: string;
  summary: string;
}
