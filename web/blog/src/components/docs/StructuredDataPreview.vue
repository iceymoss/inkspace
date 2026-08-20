<template>
  <section class="structured-preview">
    <header>
      <div>
        <strong>{{ title }}</strong>
        <span>{{ summary }}</span>
      </div>
      <el-segmented
        v-model="mode"
        :options="['结构', '源码']"
        size="small"
      />
    </header>
    <div
      v-if="mode === '结构' && error"
      class="parse-error"
    >
      无法解析结构：{{ error }}
    </div>
    <div
      v-else-if="mode === '结构' && language === 'csv'"
      class="table-scroll"
    >
      <table>
        <tbody>
          <tr
            v-for="(row, rowIndex) in tableRows"
            :key="rowIndex"
          >
            <component
              :is="rowIndex === 0 ? 'th' : 'td'"
              v-for="(cell, cellIndex) in row"
              :key="cellIndex"
            >
              {{ cell }}
            </component>
          </tr>
        </tbody>
      </table>
    </div>
    <VueJsonPretty
      v-else-if="mode === '结构'"
      class="tree-view"
      :data="parsedValue"
      :deep="3"
      :show-length="true"
      :show-line="true"
      virtual
    />
    <CodeTextEditor
      v-else
      :model-value="content"
      :language="language"
      read-only
    />
  </section>
</template>

<script setup>
import { computed, defineAsyncComponent, ref, watch } from 'vue'
import CodeTextEditor from './CodeTextEditor.vue'

const VueJsonPretty = defineAsyncComponent(async () => {
  await import('vue-json-pretty/lib/styles.css')
  return import('vue-json-pretty')
})

const props = defineProps({
  content: { type: String, default: '' },
  language: { type: String, default: '' }
})
const mode = ref('结构')
const parsedValue = ref(null)
const error = ref('')

watch(() => [props.content, props.language], async () => {
  error.value = ''
  try {
    if (props.language === 'json') parsedValue.value = JSON.parse(props.content)
    else if (props.language === 'jsonl') parsedValue.value = props.content.split(/\r?\n/).filter(Boolean).slice(0, 500).map(line => JSON.parse(line))
    else if (props.language === 'csv') {
      const { default: Papa } = await import('papaparse')
      const result = Papa.parse(props.content, { preview: 500, skipEmptyLines: true })
      if (result.errors.length) throw new Error(result.errors[0].message)
      parsedValue.value = result.data.map(row => row.slice(0, 100))
    } else if (props.language === 'yaml') {
      const { parseDocument } = await import('yaml')
      const document = parseDocument(props.content)
      if (document.errors.length) throw document.errors[0]
      parsedValue.value = document.toJS({ maxAliasCount: 50 })
    }
  } catch (parseError) {
    parsedValue.value = null
    error.value = parseError.message
  }
}, { immediate: true })

const tableRows = computed(() => parsedValue.value || [])
const title = computed(() => ({ json: 'JSON 数据', jsonl: 'JSON Lines 数据', yaml: 'YAML 数据', csv: 'CSV 表格' }[props.language] || '结构化数据'))
const summary = computed(() => props.language === 'csv' ? `最多展示 500 行 · ${tableRows.value.length} 行` : '只读结构视图')
</script>

<style scoped>
.structured-preview { border: 1px solid var(--theme-border-light); background: var(--theme-bg-card); }
.structured-preview header { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 12px 14px; border-bottom: 1px solid var(--theme-border-light); background: var(--theme-bg-secondary); }
.structured-preview header > div { display: flex; flex-direction: column; }
.structured-preview header strong { font-size: 13px; }
.structured-preview header span { color: var(--theme-text-tertiary); font-size: 10px; }
.tree-view { max-height: 70vh; min-height: 240px; margin: 0; padding: 20px; overflow: auto; color: var(--theme-text-primary); font: 13px/1.65 'SFMono-Regular', Consolas, monospace; }
.table-scroll { max-height: 70vh; overflow: auto; }
.table-scroll table { min-width: 100%; border-collapse: collapse; font-size: 12px; white-space: nowrap; }
.table-scroll th, .table-scroll td { max-width: 360px; padding: 8px 11px; overflow: hidden; border-right: 1px solid var(--theme-border-light); border-bottom: 1px solid var(--theme-border-light); text-align: left; text-overflow: ellipsis; }
.table-scroll th { position: sticky; z-index: 1; top: 0; background: var(--theme-bg-secondary); color: var(--theme-primary); font-weight: 700; }
.parse-error { padding: 30px; color: var(--el-color-danger); }
@media (max-width: 600px) { .structured-preview { margin-right: -16px; margin-left: -16px; border-right: 0; border-left: 0; } }
</style>
