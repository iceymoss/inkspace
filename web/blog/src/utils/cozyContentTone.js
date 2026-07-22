const TONES = ['moss', 'sky', 'blush']

const SEMANTIC_TONES = [
  { tone: 'moss', pattern: /设计|产品|体验|交互|视觉|\b(?:ui|ux|design|product)\b/i },
  { tone: 'sky', pattern: /代码|开发|前端|后端|工程|编程|\b(?:code|dev|frontend|backend|engineering|vue|go)\b/i },
  { tone: 'blush', pattern: /摄影|照片|生活|旅行|随笔|\b(?:photo|photography|life|travel|essay)\b/i }
]

function resourceText(resource) {
  const tags = Array.isArray(resource?.tags) ? resource.tags.map(tag => tag?.name || tag) : []
  return [
    resource?.category?.name,
    resource?.type,
    resource?.tech_stack,
    resource?.title,
    ...tags
  ].filter(Boolean).join(' ')
}

function stableNumber(value) {
  const id = Number(value)
  if (Number.isSafeInteger(id) && id >= 0) return id
  return String(value ?? '').split('').reduce((hash, character) => ((hash * 31) + character.charCodeAt(0)) >>> 0, 0)
}

export function getCozyContentTone(resource) {
  const text = resourceText(resource)
  const semantic = SEMANTIC_TONES.find(item => item.pattern.test(text))
  return semantic?.tone || TONES[stableNumber(resource?.id) % TONES.length]
}

export function getCozyRotation(resource, index = 0) {
  const rotations = [-2, 1.6, -1.2, 2.2]
  return rotations[(stableNumber(resource?.id) + index) % rotations.length]
}
