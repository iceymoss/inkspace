<template>
  <div class="appearance-page">
    <header class="appearance-intro">
      <div>
        <p class="eyebrow">
          APPEARANCE / VOL. 01
        </p>
        <h1>选择你的阅读气质</h1>
        <p class="intro-copy">
          风格决定内容如何被编排，明暗决定此刻如何阅读。预览只在本次会话生效，应用后才会同步到你的账号。
        </p>
      </div>
      <div
        class="issue-mark"
        aria-hidden="true"
      >
        <span>外观</span>
        <i />
        <span>主题</span>
      </div>
    </header>

    <section
      class="editorial-section"
      aria-labelledby="theme-heading"
    >
      <div class="section-heading">
        <div>
          <span class="section-number">01</span>
          <h2 id="theme-heading">
            版式风格
          </h2>
        </div>
        <span class="section-note">THEME COLLECTION</span>
      </div>

      <div class="theme-grid">
        <article
          v-for="(theme, index) in themeRegistry"
          :key="theme.id"
          class="theme-card"
          :class="{
            selected: selectedTheme === theme.id,
            unavailable: theme.status !== 'available'
          }"
        >
          <button
            class="theme-card-button"
            type="button"
            :disabled="theme.status !== 'available'"
            :aria-pressed="selectedTheme === theme.id"
            @click="selectedTheme = theme.id"
          >
            <div
              class="theme-proof"
              :class="`proof-${theme.id}`"
              aria-hidden="true"
            >
              <template v-if="theme.id === 'magazine'">
                <div class="proof-masthead">
                  ISLAND JOURNAL
                </div>
                <div class="proof-rule" />
                <strong>在文字里<br>留一座岛</strong>
                <div class="proof-columns">
                  <i /><i /><i />
                </div>
              </template>
              <template v-else-if="theme.id === 'terminal'">
                <span>$ inkspace --theme</span>
                <strong>inkspace.log_</strong>
                <i>&gt; building stories...</i>
              </template>
              <template v-else-if="theme.id === 'cozy'">
                <i class="tape" />
                <strong>小屿的<br>角落</strong>
                <span>进来坐坐吧</span>
              </template>
              <template v-else>
                <span>04 / GRID</span>
                <strong>CHEN<br>YU®</strong>
                <i />
              </template>
            </div>

            <div class="theme-card-copy">
              <span class="card-index">{{ String(index + 1).padStart(2, '0') }}</span>
              <div>
                <h3>{{ theme.name }}</h3>
                <p class="theme-subtitle">
                  {{ theme.subtitle }}
                </p>
              </div>
              <span
                v-if="theme.status === 'available'"
                class="availability"
              >可用</span>
              <span
                v-else
                class="availability muted"
              >即将推出</span>
            </div>
            <p class="theme-description">
              {{ theme.description }}
            </p>
          </button>
        </article>
      </div>
    </section>

    <section
      class="editorial-section scheme-section"
      aria-labelledby="scheme-heading"
    >
      <div class="section-heading">
        <div>
          <span class="section-number">02</span>
          <h2 id="scheme-heading">
            明暗模式
          </h2>
        </div>
        <span class="section-note">READING LIGHT</span>
      </div>

      <div
        class="scheme-options"
        role="radiogroup"
        aria-label="明暗模式"
      >
        <label
          v-for="option in schemeOptions"
          :key="option.value"
          class="scheme-option"
        >
          <input
            v-model="selectedScheme"
            type="radio"
            name="color-scheme"
            :value="option.value"
          >
          <span
            class="scheme-icon"
            aria-hidden="true"
          >
            <Monitor v-if="option.value === 'system'" />
            <Sunny v-else-if="option.value === 'light'" />
            <Moon v-else />
          </span>
          <span>
            <strong>{{ option.label }}</strong>
            <small>{{ option.description }}</small>
          </span>
          <i class="radio-mark" />
        </label>
      </div>
    </section>

    <section
      class="live-proof"
      aria-labelledby="preview-heading"
    >
      <div class="proof-copy">
        <p class="eyebrow">
          LIVE PROOF
        </p>
        <h2 id="preview-heading">
          {{ selectedTheme === 'terminal' ? 'inkspace.log runtime' : '此刻的屿刊' }}
        </h2>
        <p>{{ selectedTheme === 'terminal' ? '机器语言负责状态与路径，中文正文保持清晰可读。' : '纸张会随你的选择改变光线，内容的秩序保持不变。' }}</p>
        <dl>
          <div><dt>风格</dt><dd>{{ selectedThemeMeta.name }}</dd></div>
          <div><dt>模式</dt><dd>{{ selectedSchemeLabel }}</dd></div>
          <div><dt>实际显示</dt><dd>{{ candidateResolvedScheme === 'dark' ? '深色' : '浅色' }}</dd></div>
        </dl>
      </div>
      <div
        class="mini-issue"
        :class="{
          'mini-terminal': selectedTheme === 'terminal',
          'mini-terminal-light': selectedTheme === 'terminal' && candidateResolvedScheme === 'light'
        }"
        aria-hidden="true"
      >
        <template v-if="selectedTheme === 'terminal'">
          <div class="mini-terminal-bar"><i /><i /><i /><span>visitor@inkspace.log: ~</span></div>
          <div class="mini-terminal-body">
            <p><b>❯</b> whoami</p><span>reader · creator</span>
            <p><b>❯</b> cat appearance.json</p><span>{ theme: "terminal", mode: "{{ candidateResolvedScheme }}" }</span>
            <p><b>❯</b> tail -f stories<span class="mini-caret" /></p>
          </div>
        </template>
        <template v-else>
          <div class="mini-top">
            <span>屿刊</span><small>VOL. 24</small>
          </div>
          <div class="mini-line" />
          <strong>记录那些<br><em>值得停留</em>的事</strong>
          <p>观察日常，整理知识，也把尚未完成的想法留在这里。</p>
          <div class="mini-footer">
            <span>FEATURED STORIES</span><i />
          </div>
        </template>
      </div>
    </section>

    <footer class="action-bar">
      <div
        class="status-copy"
        aria-live="polite"
      >
        <span :class="{ active: hasChanges || appearance.isPreviewing }" />
        {{ appearance.isPreviewing ? '正在预览，尚未保存' : hasChanges ? '有待预览的更改' : '已与账号偏好同步' }}
      </div>
      <div class="actions">
        <button
          v-if="appearance.isPreviewing"
          type="button"
          class="text-action"
          @click="cancelPreview"
        >
          取消预览
        </button>
        <button
          type="button"
          class="outline-action"
          :disabled="!hasChanges"
          @click="previewSelection"
        >
          预览效果
        </button>
        <button
          type="button"
          class="primary-action"
          :disabled="!canSave"
          @click="saveSelection"
        >
          {{ appearance.saving ? '正在应用…' : '应用主题' }}
        </button>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Monitor, Moon, Sunny } from '@element-plus/icons-vue'
import { useAppearanceStore } from '@/stores/appearance'
import { getTheme, themeRegistry } from '@/themes/registry'

const appearance = useAppearanceStore()
const selectedTheme = ref(appearance.savedPreference.ui_theme)
const selectedScheme = ref(appearance.savedPreference.color_scheme)

const schemeOptions = [
  { value: 'system', label: '跟随系统', description: '随设备设置自动切换' },
  { value: 'light', label: '浅色', description: '保持明亮纸张' },
  { value: 'dark', label: '深色', description: '保持低照度阅读' }
]

const candidate = computed(() => ({
  ui_theme: selectedTheme.value,
  color_scheme: selectedScheme.value
}))
const selectedThemeMeta = computed(() => getTheme(selectedTheme.value))
const selectedSchemeLabel = computed(() => schemeOptions.find((item) => item.value === selectedScheme.value)?.label)
const candidateResolvedScheme = computed(() => selectedScheme.value === 'system'
  ? (window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
  : selectedScheme.value)
const hasChanges = computed(() => (
  candidate.value.ui_theme !== appearance.savedPreference.ui_theme ||
  candidate.value.color_scheme !== appearance.savedPreference.color_scheme
))
const canSave = computed(() => !appearance.saving && (hasChanges.value || appearance.isPreviewing))

watch(
  () => appearance.savedPreference,
  (preference) => {
    if (appearance.isPreviewing) return
    selectedTheme.value = preference.ui_theme
    selectedScheme.value = preference.color_scheme
  },
  { deep: true }
)

function previewSelection() {
  appearance.preview(candidate.value)
}

function cancelPreview() {
  appearance.cancelPreview()
  selectedTheme.value = appearance.savedPreference.ui_theme
  selectedScheme.value = appearance.savedPreference.color_scheme
}

async function saveSelection() {
  try {
    await appearance.save(candidate.value)
    ElMessage.success('外观偏好已同步')
  } catch (error) {
    selectedTheme.value = appearance.savedPreference.ui_theme
    selectedScheme.value = appearance.savedPreference.color_scheme
    ElMessage.error('外观偏好保存失败，已恢复原设置')
    console.error('Failed to save appearance preference:', error)
  }
}

onBeforeRouteLeave(() => {
  if (appearance.isPreviewing) appearance.cancelPreview()
})
</script>

<style scoped>
.appearance-page {
  width: min(1160px, 100%);
  margin: 0 auto;
  padding: 34px clamp(20px, 4vw, 58px) 110px;
  color: var(--ink);
  background: var(--bg);
  border: 1px solid var(--hairline);
}

.appearance-intro {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 56px;
  min-height: 225px;
  padding: 28px 0 52px;
}

.eyebrow,
.section-note {
  margin: 0 0 22px;
  color: var(--accent);
  font-size: 11px;
  letter-spacing: .26em;
}

.appearance-intro h1 {
  margin: 0;
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: clamp(36px, 5.2vw, 62px);
  font-weight: 500;
  line-height: 1.2;
  letter-spacing: .04em;
}

.intro-copy {
  max-width: 640px;
  margin: 25px 0 0;
  color: var(--sub);
  font-size: 15px;
  line-height: 1.9;
}

.issue-mark {
  display: flex;
  align-items: center;
  gap: 18px;
  writing-mode: vertical-rl;
  color: var(--sub);
  font-family: Georgia, 'Songti SC', serif;
  font-size: 13px;
  letter-spacing: .5em;
}

.issue-mark i {
  width: 1px;
  height: 58px;
  background: var(--hairline);
}

.editorial-section {
  padding: 50px 0 58px;
  border-top: 1px solid var(--hairline);
}

.section-heading,
.section-heading > div {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 18px;
}

.section-heading > div { justify-content: flex-start; }
.section-number { color: var(--accent); font: 18px Georgia, serif; }
.section-heading h2 { margin: 0; font-size: 26px; letter-spacing: .07em; }
.section-note { margin: 0; color: var(--sub); }

.theme-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-top: 38px;
}

.theme-card {
  position: relative;
  border: 1px solid var(--hairline);
  background: var(--surface);
  transition: border-color .25s ease, transform .25s ease;
}

.theme-card.selected { border-color: var(--accent); }
.theme-card:not(.unavailable):hover { transform: translateY(-4px); border-color: var(--accent); }

.theme-card-button {
  width: 100%;
  padding: 0;
  color: inherit;
  text-align: left;
  border: 0;
  background: transparent;
  cursor: pointer;
}

.theme-card-button:disabled { cursor: default; }

.theme-proof {
  position: relative;
  height: 178px;
  padding: 20px;
  overflow: hidden;
}

.proof-magazine { color: #1b1a18; background: #f6f3eb; }
.proof-masthead { font: 8px Georgia, serif; letter-spacing: .25em; }
.proof-rule { height: 1px; margin: 9px 0 23px; background: #cfcbc0; }
.proof-magazine strong { font: 500 25px/1.35 Georgia, 'Songti SC', serif; }
.proof-columns { display: flex; gap: 5px; margin-top: 18px; }
.proof-columns i { width: 28%; height: 2px; background: #33544a; }

.proof-terminal { color: #8ca1bd; background: #0d1320; font: 9px Consolas, monospace; }
.proof-terminal strong { display: block; margin: 31px 0 18px; color: #edf3fb; font-size: 25px; }
.proof-terminal i { color: #6fcf8e; font-style: normal; }

.proof-cozy { color: #43362a; background: #fbf6ed; transform: rotate(-1deg); }
.proof-cozy .tape { position: absolute; top: 8px; left: 38%; width: 55px; height: 15px; background: rgba(217, 164, 65, .55); transform: rotate(3deg); }
.proof-cozy strong { display: block; margin: 29px 0 16px; font-size: 25px; line-height: 1.35; }
.proof-cozy span { padding-bottom: 3px; border-bottom: 2px wavy #75845c; font-size: 10px; }

.proof-swiss { color: #111; background: #fff; border-bottom: 1px solid #111; font: 8px Helvetica, sans-serif; }
.proof-swiss::after { content: ''; position: absolute; inset: 0 34% 0 auto; width: 1px; background: #ddd; }
.proof-swiss strong { display: block; margin-top: 20px; font: 800 29px/1 Helvetica, sans-serif; }
.proof-swiss i { position: absolute; right: 18px; bottom: 18px; width: 26px; height: 26px; background: #002fa7; }

.theme-card-copy {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 12px;
  align-items: start;
  padding: 20px 18px 9px;
  border-top: 1px solid var(--hairline);
}

.card-index { color: var(--accent); font: 12px Georgia, serif; }
.theme-card h3 { margin: -4px 0 2px; font-size: 18px; letter-spacing: .04em; }
.theme-subtitle { margin: 0; color: var(--sub); font-size: 11px; }
.availability { color: var(--accent); font-size: 10px; letter-spacing: .1em; }
.availability.muted { color: var(--sub); }
.theme-description { min-height: 80px; margin: 0; padding: 10px 18px 20px 45px; color: var(--sub); font-size: 12px; line-height: 1.7; }
.theme-card.unavailable .theme-proof { filter: saturate(.35); opacity: .72; }

.scheme-options {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  margin-top: 38px;
  border: 1px solid var(--hairline);
}

.scheme-option {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 16px;
  align-items: center;
  min-height: 96px;
  padding: 20px 24px;
  cursor: pointer;
  border-right: 1px solid var(--hairline);
  transition: background-color .25s ease;
}

.scheme-option:last-child { border-right: 0; }
.scheme-option:has(input:checked) { background: var(--accent-soft); }
.scheme-option input { position: absolute; opacity: 0; pointer-events: none; }
.scheme-icon { width: 25px; height: 25px; color: var(--accent); }
.scheme-icon svg { width: 100%; height: 100%; }
.scheme-option strong { display: block; margin-bottom: 6px; font-family: Georgia, 'Songti SC', serif; font-weight: 500; }
.scheme-option small { color: var(--sub); }
.radio-mark { width: 13px; height: 13px; border: 1px solid var(--hairline); border-radius: 50%; }
.scheme-option input:checked ~ .radio-mark { border: 4px solid var(--accent); }

.live-proof {
  display: grid;
  grid-template-columns: 1fr 1.15fr;
  gap: clamp(42px, 8vw, 100px);
  align-items: center;
  padding: 66px clamp(18px, 4vw, 54px);
  background: var(--bg-soft);
  border-top: 1px solid var(--hairline);
  border-bottom: 1px solid var(--hairline);
}

.proof-copy h2 { margin: 0 0 18px; font-size: 30px; letter-spacing: .06em; }
.proof-copy > p:not(.eyebrow) { color: var(--sub); line-height: 1.8; }
.proof-copy dl { margin: 30px 0 0; border-top: 1px solid var(--hairline); }
.proof-copy dl div { display: flex; justify-content: space-between; padding: 11px 0; border-bottom: 1px solid var(--hairline); font-size: 12px; }
.proof-copy dt { color: var(--sub); }
.proof-copy dd { margin: 0; color: var(--ink); }

.mini-issue {
  min-height: 310px;
  padding: 30px 34px;
  color: var(--ink);
  background: var(--surface);
  border: 1px solid var(--hairline);
}
.mini-top { display: flex; justify-content: space-between; align-items: baseline; }
.mini-top span { color: var(--accent); font: 20px Georgia, 'Songti SC', serif; }
.mini-top small { color: var(--sub); font-size: 9px; letter-spacing: .2em; }
.mini-line { height: 1px; margin: 15px 0 28px; background: var(--hairline); }
.mini-issue > strong { font: 500 clamp(26px, 3vw, 38px)/1.35 Georgia, 'Songti SC', serif; letter-spacing: .05em; }
.mini-issue em { color: var(--accent); font-style: normal; }
.mini-issue > p { max-width: 360px; margin: 18px 0; color: var(--sub); font-size: 12px; line-height: 1.8; }
.mini-footer { display: flex; align-items: center; gap: 15px; margin-top: 28px; color: var(--sub); font-size: 8px; letter-spacing: .2em; }
.mini-footer i { flex: 1; height: 1px; background: var(--hairline); }
.mini-terminal { padding: 0; overflow: hidden; border-radius: 14px; background: var(--panel, #131b2c); font-family: var(--terminal-mono, Consolas, monospace); }
.mini-terminal-light { --panel: #fff; --line: #dde4ef; --sub: #7c88a1; --bright: #0e1730; --green: #1f9d57; --accent: #1668c7; }
.mini-terminal-bar { display: flex; align-items: center; gap: 7px; padding: 12px 15px; border-bottom: 1px solid var(--line, #26344e); }
.mini-terminal-bar i { width: 10px; height: 10px; border-radius: 50%; background: #f26d6d; }
.mini-terminal-bar i:nth-child(2) { background: #f2c46d; }
.mini-terminal-bar i:nth-child(3) { background: #6dd087; }
.mini-terminal-bar span { margin: 0 auto; color: var(--sub); font-size: 10px; }
.mini-terminal-body { padding: 22px 25px; font-size: 12px; line-height: 2; }
.mini-terminal-body p, .mini-terminal-body span { margin: 0; }
.mini-terminal-body p { color: var(--bright, var(--ink)); }
.mini-terminal-body p b { color: var(--green, #6fcf8e); }
.mini-terminal-body > span { color: var(--sub); }
.mini-caret { display: inline-block; width: 7px; height: 13px; margin-left: 4px; background: var(--accent); vertical-align: -2px; animation: mini-blink 1.1s steps(1) infinite; }
@keyframes mini-blink { 50% { opacity: 0; } }

.action-bar {
  position: sticky;
  z-index: 10;
  bottom: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin: 58px -1px -111px;
  padding: 17px clamp(20px, 4vw, 58px);
  background: color-mix(in srgb, var(--bg) 92%, transparent);
  border: 1px solid var(--hairline);
  backdrop-filter: blur(12px);
}

.status-copy { color: var(--sub); font-size: 12px; }
.status-copy span { display: inline-block; width: 7px; height: 7px; margin-right: 8px; background: var(--hairline); border-radius: 50%; }
.status-copy span.active { background: var(--accent); }
.actions { display: flex; align-items: center; gap: 10px; }
.actions button { min-height: 40px; padding: 0 20px; font: inherit; cursor: pointer; }
.actions button:disabled { cursor: not-allowed; opacity: .45; }
.text-action { color: var(--sub); border: 0; background: transparent; }
.outline-action { color: var(--ink); border: 1px solid var(--ink); background: transparent; }
.primary-action { color: var(--bg); border: 1px solid var(--ink); background: var(--ink); }
.primary-action:not(:disabled):hover { border-color: var(--accent); background: var(--accent); }

@media (max-width: 900px) {
  .theme-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .scheme-options { grid-template-columns: 1fr; }
  .scheme-option { border-right: 0; border-bottom: 1px solid var(--hairline); }
  .scheme-option:last-child { border-bottom: 0; }
  .live-proof { grid-template-columns: 1fr; }
}

@media (max-width: 560px) {
  .appearance-page { padding: 20px 16px 124px; }
  .appearance-intro { grid-template-columns: 1fr; min-height: auto; padding-bottom: 38px; }
  .issue-mark, .section-note { display: none; }
  .theme-grid { grid-template-columns: 1fr; }
  .theme-description { min-height: 0; }
  .live-proof { margin: 0 -16px; padding: 44px 16px; }
  .mini-issue { padding: 25px 22px; }
  .action-bar { align-items: stretch; flex-direction: column; margin: 40px -17px -125px; padding: 13px 16px; }
  .status-copy { text-align: center; }
  .actions { display: grid; grid-template-columns: 1fr 1fr; }
  .actions .text-action { grid-column: 1 / -1; min-height: 28px; }
  .actions button { padding: 0 12px; }
}

@media (prefers-reduced-motion: reduce) { .mini-caret { animation: none; } }
</style>
