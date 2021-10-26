package main

import (
	"github.com/liliang-cn/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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
		ts, err := template.ParseFiles(page)
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
