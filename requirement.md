# RepoSync 需求文档

> 本文档汇总了用户对 reposync 项目提出的所有历史需求和改进。

---

## 1. 配置行为理解与修正

**日期**: 2026-06-23  
**问题**: 用户发现配置了 5 个具体仓库，但执行后所有仓库都被下载了。

**分析结论**:
- `default-action: include` 意味着 GitHub 账户下的 **所有仓库** 都会被同步
- rules 中的 `action: include` 规则是多余的（默认行为已经是 include）
- 要只同步特定仓库，应将 `default-action` 改为 `exclude`

**解决方案**: 修改配置文件，将 `default-action` 改为 `exclude`，rules 中保留 `action: include`。

---

## 2. 支持其他用户的公共仓库同步（URL 规则）

**日期**: 2026-06-23  
**需求**: 用户想要同步其他 GitHub 用户的公共仓库（如 `rust-kotlin/ashell`），但 `platform.Repositories()` 只返回认证用户的仓库。

**解决方案**: 添加 URL-based 直接克隆规则。
- 当 rule 中的 `path` 包含 `https://` 或 `http://` 时，视为直接克隆 URL
- 新增 `DirectClone` 结构体、`ExtractDirectClones()`、`isURLRule()`、`parseGitHubURL()` 等函数
- `EvaluateRules()` 函数跳过 URL 规则（由单独的逻辑处理）

**配置示例**:
```yaml
rules:
  - rule: path == "https://github.com/rust-kotlin/ashell"
    action: include
```

---

## 3. 添加 `--config` 参数支持

**日期**: 2026-06-23  
**问题**: `config.Load()` 没有使用 `--config` CLI 参数值。

**解决方案**: 修改 `config.Load()` 签名为 `Load(files ...string)`，支持传入配置文件路径。更新 `clone.go` 传递 `configFile` 变量。

---

## 4. 修复陈旧状态目录问题

**日期**: 2026-06-23  
**问题**: 当 state 中有条目但目录不存在时，代码跳过克隆直接执行 `updateRemote`，导致失败。

**解决方案**: 在 `clone.go` 中检查 `repository.Exists()`，即使 state 说目录匹配，如果目录不存在则重新克隆。

---

## 5. 添加 `sync` 命令

**日期**: 2026-06-24  
**需求**: 用户希望有一个命令可以完成完整的同步流程。

**解决方案**: 创建 `pkg/cmd/sync.go`，实现三步流程：
1. **clone** — 同步仓库列表（添加新仓库、移动错位的仓库、更新远程地址）
2. **pull** — 同步代码（fetch + pull --ff-only）
3. **cleanup** — 清理不存在的仓库的状态条目

**设计决策**:
- `sync` 命令复用现有的 `cloneCmd()` 和 `pullCmd()` 实例
- 创建命令实例后调用 `Execute()` 执行

---

## 6. 设置版本号为 0.2.0

**日期**: 2026-06-24  
**需求**: 将当前版本号设置为 `0.2.0`。

**解决方案**:
- `app.go`: `version = "0.2.0"`
- `pkg/cmd/version.go`: `var Version = "0.2.0"`

**注意**: build ldflags 会覆盖 Go 代码中的默认值，需要使用显式 ldflags 才能生效。

---

## 7. 添加 `--log-file` 日志文件输出功能

**日期**: 2026-06-25  
**需求**: 添加 `--log-file dir` 参数，将日志保存到文件。
- 文件名格式：`reposync-yyyymmdd-hhmm.log`
- 如果未指定 `dir`，保存到 config.yaml 同级目录
- 如果未指定 `--log-file`，保持默认行为（仅输出到 stderr）

**解决方案**:
- `pkg/cmd/root.go`: 添加 `logFileDir` 变量和 `--log-file` CLI 参数
- 在 `PersistentPreRun` 中创建日志文件，使用 `zerolog.MultiLevelWriter` 同时输出到 stderr 和文件
- `pkg/config/util.go`: 新增 `GetConfigDir()` 函数获取配置文件目录

---

## 8. 探索日志功能（只读分析）

**日期**: 2026-06-25  
**需求**: 用户想了解项目的日志功能实现。

**分析结果**:
- 日志库：zerolog v1.35.1 + zerologconfig v0.1.1
- 默认输出：`os.Stderr`，使用 `zerolog.ConsoleWriter` 格式化
- 日志级别控制：`--log-level` CLI 参数
- 日志格式控制：`--log-format` CLI 参数
- 调用者信息：`--log-caller` CLI 参数
- 文件输出：**原来没有**，已在需求 7 中添加

---

## 9. 编译验证

**日期**: 2026-06-25  
**需求**: 修改代码后编译验证。

**验证结果**: 
- `go build -o reposync .` 编译成功，无错误
- `go test ./...` 所有测试通过

---

## 涉及修改的文件汇总

| 文件 | 修改内容 |
|------|----------|
| `app.go` | 版本号设为 `0.2.0` |
| `pkg/cmd/clone.go` | 添加 URL 规则直接克隆支持、`--config` 参数传递、陈旧目录检测 |
| `pkg/cmd/root.go` | 添加 `--log-file` 参数、注册 `sync` 子命令 |
| `pkg/cmd/sync.go` | **新文件** — `sync` 子命令实现 |
| `pkg/config/config.go` | 添加 `DirectClone` 结构体、URL 规则解析逻辑 |
| `pkg/config/load.go` | `Load()` 支持传入配置文件路径参数 |
| `pkg/config/util.go` | 新增 `GetConfigDir()` 工具函数 |
| `pkg/config/config_test.go` | **新文件** — 配置相关单元测试 |
