# Graph Report - .  (2026-06-22)

## Corpus Check
- 138 files · ~92,180 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1373 nodes · 2510 edges · 84 communities (63 shown, 21 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 250 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Bet Configuration|Bet Configuration]]
- [[_COMMUNITY_Miner & Stream Engine|Miner & Stream Engine]]
- [[_COMMUNITY_Frontend Dashboard UI|Frontend Dashboard UI]]
- [[_COMMUNITY_GraphQL Client|GraphQL Client]]
- [[_COMMUNITY_Authentication Core|Authentication Core]]
- [[_COMMUNITY_Account Config & Features|Account Config & Features]]
- [[_COMMUNITY_Config Schema|Config Schema]]
- [[_COMMUNITY_Twitch Chat|Twitch Chat]]
- [[_COMMUNITY_Debug & Scheduler|Debug & Scheduler]]
- [[_COMMUNITY_Channel Points & Goals|Channel Points & Goals]]
- [[_COMMUNITY_Twitch Client Layer|Twitch Client Layer]]
- [[_COMMUNITY_Notification Batching|Notification Batching]]
- [[_COMMUNITY_Matrix Notifier|Matrix Notifier]]
- [[_COMMUNITY_PubSub Connection|PubSub Connection]]
- [[_COMMUNITY_Analytics Server|Analytics Server]]
- [[_COMMUNITY_Prediction Handler|Prediction Handler]]
- [[_COMMUNITY_Drops Testing|Drops Testing]]
- [[_COMMUNITY_Structured Logger|Structured Logger]]
- [[_COMMUNITY_Graphify Skill|Graphify Skill]]
- [[_COMMUNITY_PubSub Testing|PubSub Testing]]
- [[_COMMUNITY_GitHub Community Docs|GitHub Community Docs]]
- [[_COMMUNITY_Service Installer|Service Installer]]
- [[_COMMUNITY_Drops & Campaigns|Drops & Campaigns]]
- [[_COMMUNITY_Config Editor Server|Config Editor Server]]
- [[_COMMUNITY_PubSub Pool|PubSub Pool]]
- [[_COMMUNITY_Miner Dispatcher|Miner Dispatcher]]
- [[_COMMUNITY_Streamer Model|Streamer Model]]
- [[_COMMUNITY_Community Goal Handler|Community Goal Handler]]
- [[_COMMUNITY_Analytics Backend|Analytics Backend]]
- [[_COMMUNITY_Watcher Config|Watcher Config]]
- [[_COMMUNITY_Notification Docs|Notification Docs]]
- [[_COMMUNITY_Config Loading|Config Loading]]
- [[_COMMUNITY_Self-Updater|Self-Updater]]
- [[_COMMUNITY_Logs Frontend|Logs Frontend]]
- [[_COMMUNITY_Dashboard Frontend|Dashboard Frontend]]
- [[_COMMUNITY_Streamer Settings|Streamer Settings]]
- [[_COMMUNITY_OAuth Device Flow|OAuth Device Flow]]
- [[_COMMUNITY_Runtime Config|Runtime Config]]
- [[_COMMUNITY_GitHub Releases|GitHub Releases]]
- [[_COMMUNITY_Password Auth Flow|Password Auth Flow]]
- [[_COMMUNITY_Deployment Infra|Deployment Infra]]
- [[_COMMUNITY_Entrypoint & Logging|Entrypoint & Logging]]
- [[_COMMUNITY_PubSub Topics|PubSub Topics]]
- [[_COMMUNITY_Version Parsing|Version Parsing]]
- [[_COMMUNITY_Contribution Guidelines|Contribution Guidelines]]
- [[_COMMUNITY_README & Config Editor|README & Config Editor]]
- [[_COMMUNITY_Miner Core Loop|Miner Core Loop]]
- [[_COMMUNITY_PubSub Messages|PubSub Messages]]
- [[_COMMUNITY_Notifier Base|Notifier Base]]
- [[_COMMUNITY_Discord Notifier|Discord Notifier]]
- [[_COMMUNITY_Gotify Notifier|Gotify Notifier]]
- [[_COMMUNITY_Pushover Notifier|Pushover Notifier]]
- [[_COMMUNITY_Telegram Notifier|Telegram Notifier]]
- [[_COMMUNITY_Webhook Notifier|Webhook Notifier]]
- [[_COMMUNITY_PubSub Debug|PubSub Debug]]
- [[_COMMUNITY_Worker Pool|Worker Pool]]
- [[_COMMUNITY_Linting Config|Linting Config]]
- [[_COMMUNITY_GQL Constants|GQL Constants]]
- [[_COMMUNITY_Raid Model|Raid Model]]
- [[_COMMUNITY_Auth Interface|Auth Interface]]
- [[_COMMUNITY_Config Edit Script|Config Edit Script]]
- [[_COMMUNITY_GQL Interface|GQL Interface]]
- [[_COMMUNITY_Run Script|Run Script]]
- [[_COMMUNITY_Dashboard Auth Script|Dashboard Auth Script]]
- [[_COMMUNITY_Git Hook Installer|Git Hook Installer]]
- [[_COMMUNITY_Static HTML Pages|Static HTML Pages]]
- [[_COMMUNITY_Twitch API Interface|Twitch API Interface]]
- [[_COMMUNITY_Ideas Template|Ideas Template]]
- [[_COMMUNITY_Show & Tell Template|Show & Tell Template]]
- [[_COMMUNITY_Feature Request Template|Feature Request Template]]
- [[_COMMUNITY_PR Template|PR Template]]
- [[_COMMUNITY_Go Module Entry|Go Module Entry]]
- [[_COMMUNITY_Cross-Repo Merge|Cross-Repo Merge]]
- [[_COMMUNITY_GitHub Clone|GitHub Clone]]
- [[_COMMUNITY_Graph Node Explain|Graph Node Explain]]
- [[_COMMUNITY_Project CLAUDE|Project CLAUDE.md]]

## God Nodes (most connected - your core abstractions)
1. `DefaultBetSettings()` - 41 edges
2. `T` - 37 edges
3. `makeBet()` - 36 edges
4. `Miner` - 31 edges
5. `Connection` - 31 edges
6. `Streamer` - 30 edges
7. `Client` - 27 edges
8. `twoOutcomes()` - 27 edges
9. `Client` - 24 edges
10. `Authenticator` - 23 edges

## Surprising Connections (you probably didn't know these)
- `Environment Variables Configuration` --semantically_similar_to--> `Per-Account vs Global Notification Credentials`  [INFERRED] [semantically similar]
  README.md → .github/wiki/Notifications.md
- `Authentication Priority Chain` --conceptually_related_to--> `Troubleshooting Wiki Page`  [INFERRED]
  README.md → .github/wiki/Troubleshooting.md
- `main()` --calls--> `Bool`  [INFERRED]
  cmd/twitch-miner-go/main.go → internal/miner/miner.go
- `main()` --calls--> `getEnv()`  [INFERRED]
  cmd/twitch-miner-go/main.go → internal/config/config.go
- `main()` --calls--> `LoadAllAccountConfigs()`  [INFERRED]
  cmd/twitch-miner-go/main.go → internal/config/config.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Graphify Extraction Pipeline (AST + Semantic + Merge)** — graphify_skill_ast_extraction, graphify_skill_semantic_extraction, graphify_skill_extraction_cache, graphify_skill_knowledge_graph [EXTRACTED 1.00]
- **twitch-miner-go Core Package Architecture** — wiki_architecture_miner_package, wiki_architecture_auth_package, wiki_architecture_pubsub_package, wiki_architecture_twitch_package, wiki_architecture_model_package [EXTRACTED 1.00]
- **Graphify Query Traversal Modes** — references_query_vocab_expansion, references_query_bfs_traversal, references_query_dfs_traversal, references_query_networkx_fallback [EXTRACTED 1.00]
- **CI/CD Pipeline: Build → Version → Deploy** — workflows_ci_build_job, workflows_ci_version_job, workflows_ci_deploy_job, workflows_fly_deploy_workflow [EXTRACTED 1.00]
- **Notification Provider System (6 providers + batching + events)** — wiki_notifications_telegram_provider, wiki_notifications_discord_provider, wiki_notifications_webhook_provider, wiki_notifications_matrix_provider, wiki_notifications_pushover_provider, wiki_notifications_gotify_provider, wiki_notifications_batching, wiki_notifications_events [INFERRED 0.95]
- **Prediction Strategy Suite** — wiki_prediction_strategies_smart, wiki_prediction_strategies_most_voted, wiki_prediction_strategies_high_odds, wiki_prediction_strategies_percentage, wiki_prediction_strategies_smart_money, wiki_prediction_strategies_fixed, wiki_prediction_strategies_stealth, wiki_prediction_strategies_delay_modes [EXTRACTED 1.00]

## Communities (84 total, 21 thin omitted)

### Community 0 - "Bet Configuration"
Cohesion: 0.06
Nodes (76): B, Bet, BetSettingsConfig, FilterConditionConfig, BetSettings, Mutex, Streamer, Time (+68 more)

### Community 1 - "Miner & Stream Engine"
Cohesion: 0.05
Nodes (38): Context, Miner, Streamer, Campaign, Duration, GameInfo, Tag, Time (+30 more)

### Community 2 - "Frontend Dashboard UI"
Cohesion: 0.08
Nodes (58): addCategoryItem(), addStreamerItem(), addTag(), addTeamItem(), api(), assignNum(), assignTriToggle(), collectCategories() (+50 more)

### Community 3 - "GraphQL Client"
Cohesion: 0.08
Nodes (29): circuitBreaker, isRetryableGQLError(), IsTransientError(), NewClient(), NewClientForTest(), TestIsRetryableGQLError(), TestIsTransientError(), wrapTransientGQLError() (+21 more)

### Community 4 - "Authentication Core"
Cohesion: 0.07
Nodes (19): generateDeviceID(), GenerateHex(), NewAuthenticator(), NewForTest(), Cookie, CookieJar, CookieFileExists(), NewCookieJar() (+11 more)

### Community 5 - "Account Config & Features"
Cohesion: 0.07
Nodes (28): Batcher, FeaturesConfig, FollowersConfig, AuthConfig, CategoryWatcherConfig, AccountConfig, NotificationsConfig, Priority (+20 more)

### Community 6 - "Config Schema"
Cohesion: 0.08
Nodes (36): BetSettingsConfig, CategoryConfig, ResolveBatchConfig(), AuthConfig, boolPtr(), TestIsBatchEnabled_Nil(), TestResolveBatchConfig_BothNil(), TestResolveBatchConfig_GlobalOnly() (+28 more)

### Community 7 - "Twitch Chat"
Cohesion: 0.06
Nodes (22): NewManager(), Handler, NewHandler(), Manager, Client, Context, Handler, Logger (+14 more)

### Community 8 - "Debug & Scheduler"
Cohesion: 0.09
Nodes (27): ConnectionSnapshot, DebugPredictionEntry, DebugWatchingEntry, Miner, Time, Context, Miner, Streamer (+19 more)

### Community 9 - "Channel Points & Goals"
Cohesion: 0.09
Nodes (16): GoalContribution, ChannelPointsContext, GameResp, GoalContribution, PlaybackAccessToken, StreamInfoResponse, TeamMember, TopStream (+8 more)

### Community 10 - "Twitch Client Layer"
Cohesion: 0.10
Nodes (17): Authenticator, AccountConfig, Context, Logger, Provider, RWMutex, Streamer, Time (+9 more)

### Community 11 - "Notification Batching"
Cohesion: 0.12
Nodes (28): BatchConfig, batchKey, Context, Duration, Event, Logger, Mutex, Once (+20 more)

### Community 12 - "Matrix Notifier"
Cohesion: 0.12
Nodes (31): Int64, Time, T, baseNotifier, Client, Context, Event, Message (+23 more)

### Community 13 - "PubSub Connection"
Cohesion: 0.12
Nodes (14): Conn, Context, Logger, Message, Mutex, Once, Provider, Connection (+6 more)

### Community 14 - "Analytics Server"
Cohesion: 0.13
Nodes (22): historyAggregate, HistoryEntry, Request, AnalyticsServer, Streamer, T, ResponseWriter, errorResponse (+14 more)

### Community 15 - "Prediction Handler"
Cohesion: 0.14
Nodes (18): Outcome, Context, EventPrediction, Message, Miner, Streamer, BoolFromMap(), FloatFromAny() (+10 more)

### Community 16 - "Drops Testing"
Cohesion: 0.13
Nodes (21): Client, Context, Mutex, Request, Response, T, inventoryDrop, claimFailedResponse() (+13 more)

### Community 17 - "Structured Logger"
Cohesion: 0.13
Nodes (19): Attr, Context, Event, Handler, Mutex, Level, colorHandler, Config (+11 more)

### Community 18 - "Graphify Skill"
Cohesion: 0.07
Nodes (32): Graphify Skill Entry Point, AST Structural Extraction, Community Detection, Extraction Cache, God Nodes Analysis, Graphify Full Pipeline, HTML Graph Visualization, Knowledge Graph (+24 more)

### Community 19 - "PubSub Testing"
Cohesion: 0.09
Nodes (13): Connection, Context, Mutex, Provider, T, newTestConnection(), TestHandleResponse_ERR_BADAUTH_AlreadyRefreshedByAnother(), TestHandleResponse_ERR_BADAUTH_RefreshesAndResubscribes() (+5 more)

### Community 20 - "GitHub Community Docs"
Cohesion: 0.09
Nodes (31): Q&A Discussion Template, Bug Report Issue Template, Configuration Help Issue Template, auth Package (OAuth2 + Token), chat Package (IRC Manager), gql Package (Twitch GraphQL Client), Per-Account Miner Lifecycle, miner Package (Core Orchestrator) (+23 more)

### Community 21 - "Service Installer"
Cohesion: 0.18
Nodes (28): ask(), banner(), confirm(), DEFAULT_CONFIG_DIR, DEFAULT_DATA_DIR, DEFAULT_ENV_FILE, DEFAULT_INSTALL_DIR, DEFAULT_LOG_LEVEL (+20 more)

### Community 22 - "Drops & Campaigns"
Cohesion: 0.11
Nodes (16): Drop, GameInfo, Time, Time, Campaign, Context, RawMessage, Streamer (+8 more)

### Community 23 - "Config Editor Server"
Cohesion: 0.12
Nodes (25): args, atomicWrite(), cleanConfig(), configDir, configPath(), fs, handleRequest(), http (+17 more)

### Community 24 - "PubSub Pool"
Cohesion: 0.17
Nodes (11): Connection, Context, Logger, Message, Mutex, Provider, Pool, PubSubTopic (+3 more)

### Community 25 - "Miner Dispatcher"
Cohesion: 0.11
Nodes (18): API, CategoryWatcher, Dispatcher, AccountConfig, EventPrediction, Logger, Miner, Mutex (+10 more)

### Community 26 - "Streamer Model"
Cohesion: 0.09
Nodes (11): CommunityGoal, HistoryEntry, PointsMultiplier, RWMutex, StreamerSettings, Time, HistoryEntry, PointsMultiplier (+3 more)

### Community 27 - "Community Goal Handler"
Cohesion: 0.30
Nodes (9): CommunityGoal, Context, Event, Message, Miner, Streamer, extractNestedInt(), mapReasonToEvent() (+1 more)

### Community 28 - "Analytics Backend"
Cohesion: 0.18
Nodes (15): Handler, Logger, RWMutex, AnalyticsServer, Streamer, Server, checkCredentials(), NewAnalyticsServer() (+7 more)

### Community 29 - "Watcher Config"
Cohesion: 0.19
Nodes (15): Category Watcher Feature, Team Watcher Feature, guliveer_2 Account Config, guliveer_ Account Config, Bet Amount Calculation Formula, Bet Delay Modes, Prediction Strategies Wiki Page, Prediction Filter Conditions (+7 more)

### Community 30 - "Notification Docs"
Cohesion: 0.16
Nodes (15): POST /api/test-notification Endpoint, Wiki Footer Navigation, Notification Batching Mechanism, Discord Notification Provider, Notifications Wiki Page, Per-Account vs Global Notification Credentials, Notification Event Types, Gotify Notification Provider (+7 more)

### Community 31 - "Config Loading"
Cohesion: 0.26
Nodes (12): applyDefaults(), applyEnvOverrides(), getEnv(), LoadAccountConfig(), LoadAllAccountConfigs(), parseProxyURL(), TestApplyDefaultsSetsMaxWatchStreams(), TestValidateRejectsInvalidMaxWatchStreams() (+4 more)

### Community 32 - "Self-Updater"
Cohesion: 0.29
Nodes (14): Context, T, checkWithURL(), DownloadAsset(), TestCheckForUpdate_DevVersion(), TestCheckForUpdate_NewerAvailable(), TestCheckForUpdate_NoMatchingAsset(), TestCheckForUpdate_PopulatesAssetURL() (+6 more)

### Community 33 - "Logs Frontend"
Cohesion: 0.27
Nodes (13): buildFilterParams(), clearFilters(), escapeHTML(), fetchJSON(), formatPoints(), handleSort(), loadFilters(), populateSelect() (+5 more)

### Community 34 - "Dashboard Frontend"
Cohesion: 0.31
Nodes (13): buildFilterParams(), clearFilters(), debounce(), escapeHTML(), fetchJSON(), formatPoints(), initFilterListeners(), loadFilters() (+5 more)

### Community 35 - "Streamer Settings"
Cohesion: 0.17
Nodes (8): StreamerSettings, StreamerSettings, BetSettings, ChatPresence, DefaultStreamerSettings(), ParseChatPresence(), ShouldJoinChat(), StreamerSettings

### Community 36 - "OAuth Device Flow"
Cohesion: 0.35
Nodes (5): DeviceCodeResponse, TokenErrorResponse, TokenResponse, Authenticator, Context

### Community 37 - "Runtime Config"
Cohesion: 0.27
Nodes (8): Logger, T, Twitch, envOrDefault(), LoadTwitchFromEnv(), TestClientIDsForGQL_Dedup(), TestLoadTwitchFromEnv_Defaults(), TestLoadTwitchFromEnv_EnvOverride()

### Community 38 - "GitHub Releases"
Cohesion: 0.25
Nodes (10): ghAsset, ghAsset, ghRelease, UpdateInfo, CheckForUpdate(), findAssetURL(), FormatNotification(), isGitRepo() (+2 more)

### Community 39 - "Password Auth Flow"
Cohesion: 0.42
Nodes (4): promptLine(), loginResponse, Authenticator, Context

### Community 40 - "Deployment Infra"
Cohesion: 0.24
Nodes (8): GHCR Docker Image (ghcr.io/guliveer/twitch-miner-go), Docker Compose Deployment, Fly.io Deployment, Linux Service Deployment (systemd / OpenRC), Windows Service Deployment (NSSM), CI Deploy Job (calls fly-deploy), Docker Publish GitHub Actions Workflow, Fly.io Deploy GitHub Actions Workflow

### Community 41 - "Entrypoint & Logging"
Cohesion: 0.28
Nodes (7): Bool, Logger, ColorSupported(), main(), playStartupAnimation(), runHealthcheck(), ExitForRestart()

### Community 42 - "PubSub Topics"
Cohesion: 0.42
Nodes (5): Streamer, NewStreamerTopic(), NewUserTopic(), PubSubTopic, PubSubTopicType

### Community 43 - "Version Parsing"
Cohesion: 0.36
Nodes (7): T, Version, Compare(), Parse(), TestCompare(), TestParse(), TestString()

### Community 44 - "Contribution Guidelines"
Cohesion: 0.29
Nodes (7): VERSION File, Conventional Commits Convention, Git Hooks Setup (install-hooks.sh), CI Commit Lint Job, Conventional Commits Enforcement in CI, CI Automated Version Bump Job, CI GitHub Actions Workflow

### Community 45 - "README & Config Editor"
Cohesion: 0.29
Nodes (6): Config Editor Frontend HTML, Authentication Priority Chain, Auto-Update Mechanism (-auto-update flag), Config Editor GUI Tool, Environment Variables Configuration, Twitch Runtime Identifiers (Client IDs)

### Community 47 - "PubSub Messages"
Cohesion: 0.29
Nodes (6): RawMessage, MessageData, Request, RequestData, Response, RequestData

### Community 49 - "Discord Notifier"
Cohesion: 0.40
Nodes (5): baseNotifier, Client, Context, Event, Discord

### Community 50 - "Gotify Notifier"
Cohesion: 0.40
Nodes (5): baseNotifier, Client, Context, Event, Gotify

### Community 51 - "Pushover Notifier"
Cohesion: 0.40
Nodes (5): baseNotifier, Client, Context, Event, Pushover

### Community 52 - "Telegram Notifier"
Cohesion: 0.40
Nodes (5): baseNotifier, Client, Context, Event, Telegram

### Community 53 - "Webhook Notifier"
Cohesion: 0.40
Nodes (5): baseNotifier, Client, Context, Event, Webhook

### Community 54 - "PubSub Debug"
Cohesion: 0.40
Nodes (3): Connection, Pool, ConnectionSnapshot

### Community 55 - "Worker Pool"
Cohesion: 0.50
Nodes (3): Context, T, Run()

### Community 56 - "Linting Config"
Cohesion: 1.00
Nodes (3): Enabled Linters (govet, ineffassign, staticcheck, unused), GolangCI Lint Configuration, CI Build Job

## Knowledge Gaps
- **290 isolated node(s):** `edit-config.sh script`, `github.com/Guliveer/twitch-miner-go`, `DEFAULT_SERVICE_NAME`, `DEFAULT_INSTALL_DIR`, `DEFAULT_CONFIG_DIR` (+285 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **21 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Entrypoint & Logging` to `Self-Updater`, `Runtime Config`, `GitHub Releases`, `Twitch Chat`, `Version Parsing`, `Structured Logger`, `Miner Dispatcher`, `Analytics Backend`, `Config Loading`?**
  _High betweenness centrality (0.137) - this node is a cross-community bridge._
- **Why does `Parse()` connect `Version Parsing` to `Self-Updater`, `Entrypoint & Logging`, `Matrix Notifier`, `Prediction Handler`, `Webhook Notifier`, `Drops & Campaigns`, `Config Loading`?**
  _High betweenness centrality (0.124) - this node is a cross-community bridge._
- **Why does `Setup()` connect `Structured Logger` to `PubSub Testing`, `Drops Testing`, `Notification Batching`, `Entrypoint & Logging`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Are the 39 inferred relationships involving `DefaultBetSettings()` (e.g. with `BenchmarkBetCalculate()` and `BenchmarkFilterConditionSkip()`) actually correct?**
  _`DefaultBetSettings()` has 39 INFERRED edges - model-reasoned connections that need verification._
- **What connects `edit-config.sh script`, `github.com/Guliveer/twitch-miner-go`, `DEFAULT_SERVICE_NAME` to the rest of the system?**
  _293 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Bet Configuration` be split into smaller, more focused modules?**
  _Cohesion score 0.06013745704467354 - nodes in this community are weakly interconnected._
- **Should `Miner & Stream Engine` be split into smaller, more focused modules?**
  _Cohesion score 0.05069124423963134 - nodes in this community are weakly interconnected._