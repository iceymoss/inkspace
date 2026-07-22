export function tokenizeCommand(input) {
  if (typeof input !== 'string') throw new TypeError('Command input must be a string')

  const tokens = []
  let token = ''
  let quote = null
  let escaping = false
  let tokenStarted = false

  for (const character of input) {
    if (escaping) {
      token += character
      tokenStarted = true
      escaping = false
      continue
    }

    if (character === '\\' && quote !== "'") {
      escaping = true
      tokenStarted = true
      continue
    }

    if (quote) {
      if (character === quote) quote = null
      else token += character
      continue
    }

    if (character === '"' || character === "'") {
      quote = character
      tokenStarted = true
      continue
    }

    if (/\s/.test(character)) {
      if (tokenStarted) {
        tokens.push(token)
        token = ''
        tokenStarted = false
      }
      continue
    }

    token += character
    tokenStarted = true
  }

  if (quote) throw new SyntaxError('Unterminated quote')
  if (escaping) throw new SyntaxError('Trailing escape character')
  if (tokenStarted) tokens.push(token)

  return tokens
}

export function parseCommand(input) {
  const [name = '', ...args] = tokenizeCommand(input)
  return { name: name.toLowerCase(), args }
}
