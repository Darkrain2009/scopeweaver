// Z.ai documentation verified 2026-10-03. Presets change form fields only.
// https://docs.z.ai/guides/llm/glm-5.3
// https://zcode.z.ai/en/docs/configuration
export const LLM_PROFILE_PRESET_IDS = ["zai-general", "zai-coding-reference", "openrouter"] as const;
export type LLMProfilePresetId = (typeof LLM_PROFILE_PRESET_IDS)[number];

// OpenRouter model IDs offered by the new-profile dropdown. Order is the menu order.
export const OPENROUTER_MODEL_IDS = ["z-ai/glm-5.3-flash"] as const;

export interface LLMProfilePresetFields {
  format: "openai";
  model: string;
  base_url: string;
  thinking_type: "enabled";
  reasoning_effort: "max";
  context_window_k: number;
  max_tokens: number;
  max_tokens_field: "";
}

export function getLLMProfilePreset(id: string): Readonly<LLMProfilePresetFields> {
  if (!(LLM_PROFILE_PRESET_IDS as readonly string[]).includes(id)) {
    throw new RangeError(`Unknown LLM profile preset: ${id}`);
  }
  const openrouter = id === "openrouter";
  return {
    format: "openai",
    model: openrouter ? OPENROUTER_MODEL_IDS[0] : "glm-5.3",
    base_url:
      id === "zai-general"
        ? "https://api.z.ai/api/paas/v4"
        : id === "zai-coding-reference"
          ? "https://api.z.ai/api/coding/paas/v4"
          : "https://openrouter.ai/api/v1",
    thinking_type: "enabled",
    reasoning_effort: "max",
    context_window_k: 1000,
    // Preserve the existing default: omit the output limit and use provider defaults.
    max_tokens: 0,
    max_tokens_field: "",
  };
}

// Unknown fields, including credentials and profile names, survive unchanged.
// The original draft is not mutated, and no request or persistence occurs here.
export function applyLLMProfilePreset<T extends object>(draft: T, id: string): T & LLMProfilePresetFields {
  return { ...draft, ...getLLMProfilePreset(id) };
}
