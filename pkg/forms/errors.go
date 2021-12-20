package forms

// 自定义一种错误类型，用来存放表格校验错误，key 值使用表格中的key
type errors map[string][]string

// Add 添加错误
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get 从错误信息中获取第一个错误
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
