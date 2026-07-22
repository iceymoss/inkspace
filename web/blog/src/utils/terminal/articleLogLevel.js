const HASH_LEVELS = ['INFO', 'NOTE', 'WARN']
const NOTE_TAGS = ['观点', '思考', '随笔', 'opinion']
const WARN_TAGS = ['踩坑', '故障', '复盘', 'warning']

function tagValue(tag) {
  if (tag && typeof tag === 'object') return `${tag.id ?? ''}:${tag.name ?? ''}`
  return String(tag ?? '')
}

function stableHash(value) {
  let hash = 2166136261
  for (let index = 0; index < value.length; index += 1) {
    hash ^= value.charCodeAt(index)
    hash = Math.imul(hash, 16777619)
  }
  return hash >>> 0
}

export function getArticleLogLevel(article = {}) {
  if (article?.is_top) return 'FEAT'

  const tags = Array.isArray(article?.tags) ? article.tags.map(tagValue) : []
  const semanticTags = tags.join('|').toLowerCase()
  if (WARN_TAGS.some(tag => semanticTags.includes(tag))) return 'WARN'
  if (NOTE_TAGS.some(tag => semanticTags.includes(tag))) return 'NOTE'

  const category = article?.category
  const categoryId = category && typeof category === 'object' ? category.id : article?.category_id
  const categoryName = category && typeof category === 'object' ? category.name : category
  const hashInput = [article?.id ?? '', categoryId ?? '', categoryName ?? '', ...tags.sort()].join('|')
  return HASH_LEVELS[stableHash(hashInput) % HASH_LEVELS.length]
}
