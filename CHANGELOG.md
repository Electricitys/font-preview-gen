# Changelog

All notable changes to this project will be documented in this file.


## 0.2.0 (2026-07-07)
### Minor Changes
- feat: Flag output, and stdout support
  - Refactor cmd/font-gen to use flag parsing, support an output flag
  and stdout when omitted.
  - Extract generation logic into internal/generator for reuse.
  - Add example binary and update file permissions.

## 0.1.2 (2026-07-07)
### Patch Changes
- chore: Refactor release workflow to use dynamic versioning

## 0.1.1 (2026-07-07)
### Patch Changes
- build: Update Makefile with dynamic versioning and changeset support
  - Implement dynamic version lookup from version.txt
  - Add changeset and version-bump targets
  - Silence build commands output

