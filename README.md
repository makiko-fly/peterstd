# peterstd

peterstd 中存放着一些独立的，与业务无关的基础库。

使用时优先使用每个独立package，尽量不要使用peterstd命名空间下的代码。

peterstd命名空间下的代码应该尽量组织成独立的package。

## packages

- cache: 内存缓存
- config: 配置
- copier: 结构体复制
- datasource: 数据源，mysql, redis等
- datetime: 时间日期
- doc: 文档
- dtask: 分布式任务
- env: 环境变量
- http: http 服务
- json: josn序列化反序列化
- logging: 日志
- writers: io.Writer实现
- nsq: nsq生产者消费者
- rpc: rpc 服务
- server: server层抽象
- utils: 工具库
