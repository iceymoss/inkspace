package service

import "errors"

var (
	ErrKnowledgeNotFound = errors.New("资源不存在")
	ErrCatalogCycle      = errors.New("不能将目录移动到自身或其子目录下")
	ErrShareDisabled     = errors.New("分享链接已被作者关闭")
	ErrShareExpired      = errors.New("分享链接已过期")
	ErrKnowledgeInvalid  = errors.New("请求参数无效")
)
