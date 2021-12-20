package main

import (
	"github.com/liliang-cn/snippetbox/pkg/forms"
	"github.com/liliang-cn/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// 定义一个类型用来存放需要插入到html模版中的动态数据
type templateData struct {
	Snippet     *models.Snippet   // 某一条 Snippet
	Snippets    []*models.Snippet // Snippet 列表
	CurrentYear int               // 当前年
	Form        *forms.Form       // 提交的数据和校验信息
}

// 定义一个格式化时间的函数
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// 初始化一个 template.FuncMap 对象，将其存为全局变量，相当于一个查找映射
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// 初始化一个新的map用作缓存
	cache := map[string]*template.Template{}
	// 使用filepath.Glob函数获得一个所有以.page.tmpl结尾的文件路径的slice
	// 此处就是获取到所有的page模版
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// 遍历所有的page模版
	for _, page := range pages {
		// 将文件名（如'home.page.tmpl'）从完整路径中提取出来并赋给name变量
		name := filepath.Base(page)
		// 解析页面模版为模版集合
		// FuncMap 必须在ParseFiles()之前注册进去，意味着我们需要通过template.New()来新建一个空模版
		// 然后使用Func()方法来注册FuncMap
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// 使用 ParseGlob 方法将所有的'layout'模版添加到模版集合中
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// 使用 ParseGlob 方法将所有的'partial'模版添加到模版集合中
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// 把模版集合添加到缓存中，key值为页面的文件名
		cache[name] = ts
	}

	// 将map返回
	return cache, nil
}
