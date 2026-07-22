export function formatSwissCode(prefix, id, minLength = 3) {
  const value = Number(id)
  if (!Number.isInteger(value) || value <= 0) return `${prefix}—`
  return `${prefix}${String(value).padStart(minLength, '0')}`
}

const PLATE_SPANS = [6, 3, 3, 3, 3, 6]

export function getSwissPlateSpan(index) {
  return PLATE_SPANS[index % PLATE_SPANS.length]
}
