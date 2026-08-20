# 知识库文件上传样例

本目录由 go run scripts/generate_knowledge_file_samples.go 生成，用于知识库文件上传、类型识别和预览联调。

| 文件 | 预期 kind | 预期状态 |
|------|-----------|----------|
| sample.md | markdown | ready，可在线编辑 |
| sample.txt | text | ready，可在线编辑 |
| sample.go | code | ready，可在线编辑 |
| sample.json | code/json | ready，可在线编辑 |
| sample.yaml | code/yaml | ready，可在线编辑 |
| sample.csv | code/csv | ready，可在线编辑 |
| sample.pdf | file | ready |
| sample.png | file/image | ready |
| sample.svg | file/image | pending，等待安全清洗 |
| sample.zip | file/archive | pending |
| sample.gz | file/archive | pending |
| sample.docx | file/office | pending |
| sample.xlsx | file/office | pending |
| sample.pptx | file/office | pending |
| sample.odt | file/office | pending |
| sample.bin | fallback file | 默认拒绝；配置允许未知扩展名后接收 |
| invalid-disguised.pdf | - | 应拒绝，内容不是 PDF |
| invalid-binary.txt | - | 应拒绝，包含 NUL 字节 |

Office 文件是满足当前上传策略检测的最小容器，不用于验证 LibreOffice 的视觉转换质量。
