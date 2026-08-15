const typeMap = {
  md: ['markdown', 'markdown', 'Markdown'], markdown: ['markdown', 'markdown', 'Markdown'],
  txt: ['text', 'text', '纯文本'], log: ['text', 'text', '日志'], ini: ['text', 'ini', 'INI 配置'], conf: ['text', 'text', '配置文件'], env: ['text', 'dotenv', '环境变量'],
  json: ['code', 'json', 'JSON'], jsonl: ['code', 'jsonl', 'JSON Lines'], yaml: ['code', 'yaml', 'YAML'], yml: ['code', 'yaml', 'YAML'],
  xml: ['code', 'xml', 'XML'], toml: ['code', 'toml', 'TOML'], csv: ['code', 'csv', 'CSV'],
  go: ['code', 'go', 'Go'], js: ['code', 'javascript', 'JavaScript'], jsx: ['code', 'jsx', 'JSX'], ts: ['code', 'typescript', 'TypeScript'], tsx: ['code', 'tsx', 'TSX'], vue: ['code', 'vue', 'Vue'],
  py: ['code', 'python', 'Python'], java: ['code', 'java', 'Java'], kt: ['code', 'kotlin', 'Kotlin'], kts: ['code', 'kotlin', 'Kotlin Script'],
  c: ['code', 'c', 'C'], h: ['code', 'c', 'C Header'], cpp: ['code', 'cpp', 'C++'], hpp: ['code', 'cpp', 'C++ Header'], cs: ['code', 'csharp', 'C#'], rs: ['code', 'rust', 'Rust'],
  php: ['code', 'php', 'PHP'], rb: ['code', 'ruby', 'Ruby'], swift: ['code', 'swift', 'Swift'], scala: ['code', 'scala', 'Scala'],
  sh: ['code', 'shell', 'Shell'], bash: ['code', 'shell', 'Bash'], zsh: ['code', 'shell', 'Zsh'], fish: ['code', 'shell', 'Fish'], sql: ['code', 'sql', 'SQL'],
  html: ['code', 'html', 'HTML'], css: ['code', 'css', 'CSS'], scss: ['code', 'scss', 'SCSS'], less: ['code', 'less', 'Less'], dockerfile: ['code', 'dockerfile', 'Dockerfile']
}

export const textFileAccept = Object.keys(typeMap).map(extension => `.${extension}`).join(',')

export function inferCreatableFileType(fileName) {
  const normalized = fileName.trim()
  const baseName = normalized.split(/[\\/]/).pop()
  const extension = baseName?.toLowerCase() === 'dockerfile' ? 'dockerfile' : baseName?.split('.').pop()?.toLowerCase()
  const match = typeMap[extension]
  if (!match) return null
  return { kind: match[0], language: match[1], label: match[2] }
}
