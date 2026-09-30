# Graph Report - goxide  (2026-09-30)

## Corpus Check
- 43 files · ~147,134 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 7 file(s) not represented in the graph (top: (none) 6, .css 1)

## Summary
- 460 nodes · 1145 edges · 20 communities (17 shown, 3 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 38 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2ca3a658`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- render.go
- examples.go
- sqlEngine
- result_benchmark_with_db_test.go
- result_test.go
- Result Chaining
- Result
- T
- Goxide README
- Fluent Chain Pipeline
- ApplyToResult
- context.Context
- script.js
- Goxide Logo and Tagline
- goxide Logo Image
- github.com/seyallius/goxide
- Option Type
- Result Type in Option Documentation
- profparse/main.go
- Repository Guidelines

## God Nodes (most connected - your core abstractions)
1. `Result` - 65 edges
2. `Catch()` - 37 edges
3. `Wrap()` - 29 edges
4. `repos()` - 23 edges
5. `clearUsers()` - 20 edges
6. `CatchWith()` - 18 edges
7. `If()` - 14 edges
8. `User` - 13 edges
9. `Fallback()` - 12 edges
10. `Result[T]` - 12 edges

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
- **Performance Documentation Flow** — docs_hardware_environment, docs_benchmarks_raw, docs_tools_benchmark_documentation_pipeline [EXTRACTED 1.00]
- **Recovery and Boundary Adaptation Mechanisms** — catch_recovery, error_recovery_tools, catcherr_boundary_adapter [EXTRACTED 1.00]
- **Goxide Result, Option, and Chain Ecosystem** — result_type, option_type, chain_pipeline [EXTRACTED 1.00]
- **Composable functional pipeline concepts** — rusty_chain_readme_chain_fluent_method_chaining, rusty_result_readme_result_and_then, rusty_types_readme_types_function_composition, rusty_types_readme_types_higher_order_functions [INFERRED 0.75]
- **Result benchmark profiling family** — docs_profiles_cpu_full_result_benchmark_profile, docs_profiles_mem_objects_full_result_benchmark_allocation_profile, docs_profiles_mem_space_full_result_benchmark_allocation_space_profile, rusty_result_raw_benchmarks_result_chained_operations [INFERRED 0.75]

## Communities (20 total, 3 thin omitted)

### Community 0 - "render.go"
Cohesion: 0.08
Nodes (37): BenchResult, Category, Config, GroupRule, AggregateByMedian(), applyGroupRules(), DetectStyle(), GroupByWorkload() (+29 more)

### Community 1 - "examples.go"
Cohesion: 0.09
Nodes (39): Cache, Config, Order, User, TestResultChain_MultipleOperations(), chargePayment(), chargePaymentTraditional(), enrichUserData() (+31 more)

### Community 2 - "sqlEngine"
Cohesion: 0.20
Nodes (4): GoxideDataInit(), GoxideEngineInit(), GoxideSchemaInit(), sqlEngine

### Community 3 - "result_benchmark_with_db_test.go"
Cohesion: 0.10
Nodes (37): DB(), RunGoxideTestMain(), TestProfile(), BenchmarkResultChainedSuccess(), BenchmarkResultError(), BenchmarkResultSuccess(), BenchmarkResultSuccessUnwrapOr(), BenchmarkTraditionalChainedSuccess() (+29 more)

### Community 4 - "result_test.go"
Cohesion: 0.11
Nodes (47): Profile, User, failOnOdd(), TestResultChain_AndThen(), TestResultChain_AndThen2(), TestResultChain_AndThen3_ErrorFromFunction(), TestResultChain_AndThen4_ErrorPropagation(), TestResultChain_Map() (+39 more)

### Community 5 - "Result Chaining"
Cohesion: 0.07
Nodes (29): CPU Profile: BenchmarkResultWithTry, Full CPU Profile: Result Benchmarks, Allocation Objects Profile: BenchmarkResultWithTry, Full Allocation Objects Profile: Result Benchmarks, Allocation Space Profile: BenchmarkResultWithTry, Full Allocation Space Profile: Result Benchmarks, Chain2 Multi-Step Chaining, Chain Package (+21 more)

### Community 6 - "Result"
Cohesion: 0.18
Nodes (17): Result[T], tryError, AndThen(), CatchErr(), Err(), FlatMap(), Result, If() (+9 more)

### Community 7 - "T"
Cohesion: 0.16
Nodes (13): Option[T], Cast(), FlatMap(), Option, If(), Map(), None(), Some() (+5 more)

### Community 8 - "Goxide README"
Cohesion: 0.09
Nodes (23): Go Result Benchmarks, Benchmark Performance Measurement, Raw Benchmark Results, Benchmark Hardware Environment, Go 1.25.4 Linux AMD64 Environment, Goxide Documentation Site, Highlight.js Syntax Highlighting, Markdown Documentation Content (+15 more)

### Community 9 - "Fluent Chain Pipeline"
Cohesion: 0.13
Nodes (22): BubbleUp Pattern, Catch Recovery, CatchErr Boundary Adapter, Fluent Chain Pipeline, Performance Benchmarks, Quick Cookbook, Database Operations Example, HTTP Handlers Example (+14 more)

### Community 10 - "ApplyToResult"
Cohesion: 0.18
Nodes (6): ApplyToResult, ApplyToResult2, ApplyToResult2[Out1, Out2, T], ApplyToResult[Out, In], Chain2(), Chain()

### Community 11 - "context.Context"
Cohesion: 0.28
Nodes (5): ResultUserRepo, TraditionalUserRepo, User, NewResultUserRepo(), NewTraditionalUserRepo()

### Community 12 - "script.js"
Cohesion: 0.33
Nodes (9): CONFIG, configureMarked(), generateTOC(), handleRouteChange(), initDocs(), loadContent(), setupNavigation(), setupScrollSpy() (+1 more)

### Community 13 - "Goxide Logo and Tagline"
Cohesion: 0.33
Nodes (6): Fluent Chaining, Go Programming Language, Goxide Logo and Tagline, Option, Result, Rust-inspired Ergonomics

### Community 14 - "goxide Logo Image"
Cohesion: 0.67
Nodes (3): goxide Logo Image, Blue Gopher-Crab Mascot, goxide Wordmark

### Community 18 - "profparse/main.go"
Cohesion: 0.25
Nodes (9): pprofRowData, section, humanKind(), kindOrder(), main(), parse(), render(), run() (+1 more)

### Community 19 - "Repository Guidelines"
Cohesion: 0.29
Nodes (6): Build, Test, and Development Commands, Coding Style & Naming Conventions, Commit & Pull Request Guidelines, Project Structure & Module Organization, Repository Guidelines, Testing Guidelines

## Ambiguous Edges - Review These
- `Chain Package` → `Chain Package in Option Documentation`  [AMBIGUOUS]
  rusty/option/README_OPTION.md · relation: semantically_similar_to

## Knowledge Gaps
- **46 isolated node(s):** `Project Structure & Module Organization`, `Build, Test, and Development Commands`, `Coding Style & Naming Conventions`, `Testing Guidelines`, `Commit & Pull Request Guidelines` (+41 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 83 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What is the exact relationship between `Chain Package` and `Chain Package in Option Documentation`?**
  _Edge tagged AMBIGUOUS (relation: semantically_similar_to) - confidence is low._
- **Why does `Result` connect `Result` to `examples.go`, `result_benchmark_with_db_test.go`, `result_test.go`, `T`, `ApplyToResult`, `context.Context`?**
  _High betweenness centrality (0.104) - this node is a cross-community bridge._
- **Why does `Option` connect `T` to `Result`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Why does `Catch()` connect `result_test.go` to `result_benchmark_with_db_test.go`, `examples.go`, `context.Context`, `Result`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **What connects `Project Structure & Module Organization`, `Build, Test, and Development Commands`, `Coding Style & Naming Conventions` to the rest of the system?**
  _46 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `render.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07894736842105263 - nodes in this community are weakly interconnected._
- **Should `examples.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08672699849170437 - nodes in this community are weakly interconnected._