## 题目元信息

- 来源：`self_built`
- 处理流程：`Codex Gold 修复`
- Bug 分类：`nil`
- 复现稳定性：`stable`
- 难度：`advanced`
- 目标 Bug 数量：`1`
- 上游 Issue：无（0-1 自建项目）

# BUG_TASK

- 来源 (source): self_built
- 题目 (title): 无完成事件的作业归档后留存信息为空，回收或查询时触发 nil 指针 panic
- Bug 分类 (category): nil / panic
- 用户可见症状 (symptom):
  对只有启动（started）或断点（checkpoint）事件、从未收到完成（finished）事件的作业执行归档后，
  后续对该归档作业进行回收或查询时，服务发生 nil 指针解引用 panic：
  - 回收路径直接触发 panic（进程崩溃）；
  - 查询路径在 HTTP 层被 recover 后返回 500 internal error，而非正常的 200 与留存时间戳。
  正常收到 finished 事件的作业归档后，回收与查询均正常，不受影响。
- 正确行为 (expected):
  所有被归档的作业（无论是否收到过 finished 事件）都必须携带完整留存信息
  （归档时间、过期时间等）。回收器与查询接口不得 panic；孤儿作业可被正常删除，
  也可被正常查询到状态与留存时间戳。
- 稳定复现命令 (reproduce):
  go test -run TestEraseArchivedOrphanJob -count=1 ./internal/retention/
- 验收标准 (acceptance):
  1. 正常完成作业归档后可被回收删除，回收与查询均不 panic。
  2. 无 finished 事件的孤儿作业归档后，回收与查询均不 panic，且返回正确状态。
  3. 所有归档作业均携带非空留存信息；归档计数与索引保持一致。
