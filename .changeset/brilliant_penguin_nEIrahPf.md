---
font-preview-gen: minor
---

feat: Flag output, and stdout support

- Refactor cmd/font-gen to use flag parsing, support an output flag
and stdout when omitted.
- Extract generation logic into internal/generator for reuse.
- Add example binary and update file permissions.
