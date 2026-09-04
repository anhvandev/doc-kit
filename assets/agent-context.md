## Working rules

**Do exactly what is asked. Nothing more, nothing less.** Bias toward caution over speed; use judgment on trivial tasks.

### Think before coding
- State assumptions explicitly. If uncertain, ask.
- If several interpretations exist, present them; do not pick one silently.
- If a simpler approach exists, say so and push back when warranted.
- If something is unclear, stop, name what is confusing, and ask.

### Simplicity first
- Minimum code that solves the problem. No speculative features, abstractions, or configurability.
- No error handling for impossible scenarios. Pick the practical option, not the elegant one.
- If 200 lines could be 50, rewrite.

### Surgical changes
- Touch only what you must. Match existing style. Do not "improve" adjacent code or formatting.
- Remove imports, variables, and functions that your change made unused; leave pre-existing dead code and mention it.
- Every changed line should trace directly to the request.

### Goal-driven execution
- Turn tasks into verifiable goals: a failing test, a passing check, a reproducible command.
- For multi-step work, state a brief plan with a verify step per item, then loop until verified.

### Answer the instruction, nothing else
- Lead with the result. Short sentences, plain words, one solution, no menus of alternatives.
- No tips, warnings, "also consider", or next steps unless asked. No closing summary.
- Do not spawn subagents for work you can do directly.

**When in doubt, ask. Do not guess and do not expand scope.**

## Project documentation: `dk`

Every document under `docs/` is created by the `dk` CLI from embedded templates; never create files there by hand. Three layers: **templates** (structure, changed only by releasing a new `dk`), **CLI** (deterministic work: create files, changelog, render, index, check), **skills** (work that needs reasoning: interviews, drafting, impact analysis). If `dk` is missing or fails, stop and tell the user to install it.

### Skills (installed with `dk skill install`)

| Skill | Use when |
|---|---|
| `doc-intake` | a new idea with nothing written yet: idea, interview, brief awaiting approval (`docs/intake/`) |
| `doc-cr` | a change to something that exists: Change Request with impact table awaiting approval (`docs/cr/`) |
| `doc-overview` | Product overview, Architecture overview, Glossary (`docs/overview/`) |
| `doc-feature-list` | planned feature table split from an approved Product brief: temporary FL-xx codes, groups, MoSCoW priority, links back to brief sections (`docs/overview/feature-list.md`) |
| `doc-adr` | a significant technical decision, numbered, immutable once accepted (`docs/adr/`) |
| `doc-feature-spec` | 11-section Feature Spec from an approved brief or accepted CR (`docs/features/`) |
| `doc-design-system` | tokens, foundations, atoms through templates, patterns (`docs/design/`) |
| `doc-design-flow` | user flow, wireframe, HTML mockup, prototype for one feature |
| `doc-test` | testing strategy, Gherkin or table test cases, UI checklist, test report (`docs/test/`) |
| `doc-plan-report` | roadmap, plan and phases, evidence-backed report, decision log, product CHANGELOG |
| `doc-release` | release brief, release notes, user guide, FAQ for end users (`docs/release/`) |
| `doc-ops` | deployment, environment, runbook, monitoring, postmortem, backup and DR (`docs/ops/`) |

### Main `dk` commands

```
dk new <type> <slug> [--from <file>] [--set k=v]    # create from template, fill frontmatter
dk changelog add <file> --summary "<real summary>" --source <CR-id|brief>
dk render <file> | --all --index                    # self-contained HTML into docs/html/
dk index [features|adr|cr|intake|user-guide|all]    # generated README.md indexes
dk check [<file>]                                   # exit code 3 on errors
dk refs <file>                                      # outgoing and incoming links
dk status | dk doctor                               # docs overview | installation check
```

### Documentation rules

- After every edit to a document: `changelog add` with a summary that describes the content, `render`, `index` for that directory, `check` for that file. The pre-commit hook blocks commits that lack a changelog line.
- Statuses use the English keys from `types.toml` (`draft`, `review`, `approved`, ...). Only a person moves a document to a final status; skills stop at "awaiting approval".
- Never edit `created`, `created_by`, `dk_version` in frontmatter; never edit the body of an ADR that is `accepted`; never edit files marked `generated: true` or anything under `docs/html/`.
- Line limits: warning above 500, error above 800 (`dk.toml` `[check]`); per-type limits live in `types.toml`.
