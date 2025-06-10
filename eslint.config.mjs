import globals from "globals"
import pluginJs from "@eslint/js"
import tseslint from "typescript-eslint"
import pluginVue from "eslint-plugin-vue"
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended"

/** @type {import('eslint').Linter.Config[]} */
export default [
  {
    ignores: [
      "**/dist/**",
      "**/node_modules/**",
      "**/coverage/**",
      "server/**",
      "infra/**",
      ".github/**",
      ".git/**",
    ],
  },
  { files: ["**/*.{js,mjs,cjs,ts,vue}"] },
  { languageOptions: { globals: globals.browser } },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  ...pluginVue.configs["flat/essential"],
  eslintPluginPrettierRecommended,
  {
    files: ["**/*.vue"],
    languageOptions: { parserOptions: { parser: tseslint.parser } },
  },
  { rules: { "vue/no-multiple-template-root": "off" } },
  { rules: { "vue/no-v-model-argument": "off" } },
  { rules: { "prettier/prettier": "error" } },
  { rules: { "@typescript-eslint/ban-ts-comment": "off" } },
  { rules: { "no-empty-function": "off" } },
  { rules: { "@typescript-eslint/no-empty-function": "off" } },
  { rules: { "vue/no-multiple-template-root": "off" } },
]
