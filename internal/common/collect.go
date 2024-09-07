package common

type Collect struct {
	List      []string // 合集所有 URL 地址
	StartPath string   // 未处理任务本地存储路径
	EndPath   string   // 已处理任务存储路径
	Name      string   // 合集名称
}
