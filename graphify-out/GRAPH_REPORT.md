# Graph Report - goxide  (2026-09-30)

## Corpus Check
- 47 files · ~155,889 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 6, .css 4, .toml 1)

## Summary
- 696 nodes · 1538 edges · 37 communities (20 shown, 17 thin omitted)
- Extraction: 96% EXTRACTED · 4% INFERRED · 0% AMBIGUOUS · INFERRED: 55 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `662be091`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Result
- highlight.min.js
- testing.B
- result_test.go
- Result Chaining
- context.Context
- T
- Goxide README
- Fluent Chain Pipeline
- app.js
- ApplyToResult
- script.js
- Goxide Logo and Tagline
- goxide Logo Image
- github.com/seyallius/goxide
- Option Type
- Result Type in Option Documentation
- Repository Guidelines
- metrics
- compare.js
- se
- ie
- .parseInline
- ce
- re
- compare.md
- timeline.md
- sqlEngine
- marked.min.js

## God Nodes (most connected - your core abstractions)
1. `Result` - 65 edges
2. `Catch()` - 37 edges
3. `Wrap()` - 29 edges
4. `se` - 28 edges
5. `ie` - 25 edges
6. `repos()` - 23 edges
7. `clearUsers()` - 20 edges
8. `CatchWith()` - 18 edges
9. `metrics` - 16 edges
10. `byId()` - 14 edges

## Surprising Connections (you probably didn't know these)
- `CPU Profile: BenchmarkResultWithTry` --semantically_similar_to--> `BubbleUp Early Return Pattern`  [INFERRED] [semantically similar]
  docs/profiles/cpu_BenchmarkResultWithTry.txt → rusty/result/README_RESULT.md
- `Result Type` --semantically_similar_to--> `Go Result Benchmarks`  [INFERRED] [semantically similar]
  README.md → docs/benchmarks_raw.txt
- `Chain Package in Option Documentation` --semantically_similar_to--> `Chain Package`  [AMBIGUOUS] [semantically similar]
  rusty/option/README_OPTION.md → rusty/chain/README_CHAIN.md
- `Goxide README` --references--> `MIT License`  [EXTRACTED]
  README.md → THIRD_PARTY_LICENSES/MIT.txt
- `Full CPU Profile: Result Benchmarks` --references--> `Result Chained Operations Benchmarks`  [INFERRED]
  docs/profiles/cpu_full.txt → rusty/result/raw_benchmarks.txt

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **BenchmarkResultWithTry profiling family** — docs_profiles_cpu_benchmarkresultwithtry_benchmark_result_with_try_profile, docs_profiles_mem_objects_benchmarkresultwithtry_allocation_profile, docs_profiles_mem_space_benchmarkresultwithtry_allocation_space_profile, rusty_result_raw_benchmarks_benchmark_result_with_try [EXTRACTED 1.00]
- **Rust-inspired Ergonomics for Go** — goxide2_rust_inspired_ergonomics, goxide2_go, goxide2_option, goxide2_result, goxide2_fluent_chaining [EXTRACTED 1.00]
- **Rust-Inspired Safety and Expressiveness** — readme_result_type, readme_option_type, readme_chain_pattern, readme_type_safety [EXTRACTED 1.00]
- **Gradual Goxide Adoption Flow** — wrap_bridge, bubbleup_pattern, catcherr_boundary_adapter, chain_pipeline [EXTRACTED 1.00]
- **Recovery and Boundary Adaptation Mechanisms** — catch_recovery, error_recovery_tools, catcherr_boundary_adapter [EXTRACTED 1.00]
- **Goxide Result, Option, and Chain Ecosystem** — result_type, option_type, chain_pipeline [EXTRACTED 1.00]
- **Composable functional pipeline concepts** — rusty_chain_readme_chain_fluent_method_chaining, rusty_result_readme_result_and_then, rusty_types_readme_types_function_composition, rusty_types_readme_types_higher_order_functions [INFERRED 0.75]
- **Result benchmark profiling family** — docs_profiles_cpu_full_result_benchmark_profile, docs_profiles_mem_objects_full_result_benchmark_allocation_profile, docs_profiles_mem_space_full_result_benchmark_allocation_space_profile, rusty_result_raw_benchmarks_result_chained_operations [INFERRED 0.75]

## Communities (37 total, 17 thin omitted)

### Community 1 - "Result"
Cohesion: 0.07
Nodes (51): Cache, Order, User, Result[T], tryError, TestResultChain_MultipleOperations(), chargePayment(), chargePaymentTraditional() (+43 more)

### Community 2 - "highlight.min.js"
Cohesion: 0.06
Nodes (36): a(), b(), be(), d(), e(), f(), b(), c() (+28 more)

### Community 3 - "testing.B"
Cohesion: 0.06
Nodes (52): ApplyToResult2, DB(), RunGoxideTestMain(), Chain(), TestProfile(), BenchmarkChainedSuccess(), BenchmarkDBChainedOperations(), BenchmarkDBCreateUser() (+44 more)

### Community 4 - "result_test.go"
Cohesion: 0.10
Nodes (52): Profile, User, Config, failOnOdd(), TestResultChain_AndThen(), TestResultChain_AndThen2(), TestResultChain_AndThen3_ErrorFromFunction(), TestResultChain_AndThen4_ErrorPropagation() (+44 more)

### Community 5 - "Result Chaining"
Cohesion: 0.07
Nodes (29): CPU Profile: BenchmarkResultWithTry, Full CPU Profile: Result Benchmarks, Allocation Objects Profile: BenchmarkResultWithTry, Full Allocation Objects Profile: Result Benchmarks, Allocation Space Profile: BenchmarkResultWithTry, Full Allocation Space Profile: Result Benchmarks, Chain2 Multi-Step Chaining, Chain Package (+21 more)

### Community 6 - "context.Context"
Cohesion: 0.31
Nodes (4): ResultUserRepo, TraditionalUserRepo, User, NewTraditionalUserRepo()

### Community 7 - "T"
Cohesion: 0.16
Nodes (13): Option[T], Cast(), FlatMap(), Option, If(), Map(), None(), Some() (+5 more)

### Community 8 - "Goxide README"
Cohesion: 0.11
Nodes (20): Go Result Benchmarks, Benchmark Performance Measurement, Raw Benchmark Results, Benchmark Hardware Environment, Go 1.25.4 Linux AMD64 Environment, Goxide Documentation Site, Highlight.js Syntax Highlighting, Markdown Documentation Content (+12 more)

### Community 9 - "Fluent Chain Pipeline"
Cohesion: 0.13
Nodes (22): BubbleUp Pattern, Catch Recovery, CatchErr Boundary Adapter, Fluent Chain Pipeline, Performance Benchmarks, Quick Cookbook, Database Operations Example, HTTP Handlers Example (+14 more)

### Community 10 - "app.js"
Cohesion: 0.10
Nodes (46): addCopyButtons(), addHeadingAnchors(), appendCell(), applyFilter(), applySiteSettings(), applyTheme(), bindGlobalHandlers(), boot() (+38 more)

### Community 11 - "ApplyToResult"
Cohesion: 0.24
Nodes (4): ApplyToResult, ApplyToResult2[Out1, Out2, T], ApplyToResult[Out, In], Chain2()

### Community 12 - "script.js"
Cohesion: 0.33
Nodes (9): CONFIG, configureMarked(), generateTOC(), handleRouteChange(), initDocs(), loadContent(), setupNavigation(), setupScrollSpy() (+1 more)

### Community 13 - "Goxide Logo and Tagline"
Cohesion: 0.33
Nodes (6): Fluent Chaining, Go Programming Language, Goxide Logo and Tagline, Option, Result, Rust-inspired Ergonomics

### Community 14 - "goxide Logo Image"
Cohesion: 0.67
Nodes (3): goxide Logo Image, Blue Gopher-Crab Mascot, goxide Wordmark

### Community 19 - "Repository Guidelines"
Cohesion: 0.29
Nodes (6): Build, Test, and Development Commands, Coding Style & Naming Conventions, Commit & Pull Request Guidelines, Project Structure & Module Organization, Repository Guidelines, Testing Guidelines

### Community 20 - "metrics"
Cohesion: 0.05
Nodes (40): data, index, runs_dir, default_path, favicon, generated_by, metrics, allocs_op (+32 more)

### Community 21 - "compare.js"
Cohesion: 0.14
Nodes (29): clampProbability(), combinations(), compareRuns(), direction(), displayName(), erfc(), erfcContinuedFraction(), erfSeries() (+21 more)

### Community 22 - "se"
Cohesion: 0.08
Nodes (3): ne(), se, te()

### Community 23 - "ie"
Cohesion: 0.12
Nodes (3): ie, W(), Y()

### Community 35 - "sqlEngine"
Cohesion: 0.20
Nodes (4): GoxideDataInit(), GoxideEngineInit(), GoxideSchemaInit(), sqlEngine

## Ambiguous Edges - Review These
- `Chain Package` → `Chain Package in Option Documentation`  [AMBIGUOUS]
  rusty/option/README_OPTION.md · relation: semantically_similar_to

## Knowledge Gaps
- **80 isolated node(s):** `CONFIG`, `github.com/seyallius/goxide`, `Build, Test, and Development Commands`, `Coding Style & Naming Conventions`, `Commit & Pull Request Guidelines` (+75 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 169 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Chain Package` and `Chain Package in Option Documentation`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **Why does `Result` connect `Result` to `testing.B`, `result_test.go`, `context.Context`, `T`, `ApplyToResult`?**
  _High betweenness centrality (0.042) - this node is a cross-community bridge._
- **Why does `pe()` connect `highlight.min.js` to `.parseInline`, `marked.min.js`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `CONFIG`, `github.com/seyallius/goxide`, `Build, Test, and Development Commands` to the rest of the system?**
  _80 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Result` be split into smaller, more focused modules?**
  _Cohesion score 0.07314814814814814 - nodes in this community are weakly interconnected._
- **Should `highlight.min.js` be split into smaller, more focused modules?**
  _Cohesion score 0.06151062867480778 - nodes in this community are weakly interconnected._
- **Should `testing.B` be split into smaller, more focused modules?**
  _Cohesion score 0.062317429406037 - nodes in this community are weakly interconnected._