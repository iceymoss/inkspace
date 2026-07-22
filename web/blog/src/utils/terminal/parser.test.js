import { describe, expect, it } from 'vitest'
import { parseCommand, tokenizeCommand } from './parser'

describe('terminal parser', () => {
  it('tokenizes whitespace and quoted arguments', () => {
    expect(tokenizeCommand('search blog "vue terminal"')).toEqual(['search', 'blog', 'vue terminal'])
    expect(tokenizeCommand("grep 'error message' blog/")).toEqual(['grep', 'error message', 'blog/'])
  })

  it('supports escaped characters and empty quoted arguments', () => {
    expect(tokenizeCommand('grep error\\ message ""')).toEqual(['grep', 'error message', ''])
    expect(tokenizeCommand('search blog "say \\"hello\\""')).toEqual(['search', 'blog', 'say "hello"'])
  })

  it('rejects malformed input instead of attempting to evaluate it', () => {
    expect(() => tokenizeCommand('search blog "unfinished')).toThrow(SyntaxError)
    expect(() => tokenizeCommand('open blog\\')).toThrow(SyntaxError)
    expect(tokenizeCommand('open "$(alert 1)"')).toEqual(['open', '$(alert 1)'])
  })

  it('returns a normalized command name and untouched arguments', () => {
    expect(parseCommand('OPEN article 12')).toEqual({ name: 'open', args: ['article', '12'] })
    expect(parseCommand('   ')).toEqual({ name: '', args: [] })
  })
})
