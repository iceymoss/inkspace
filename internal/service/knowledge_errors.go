package service

import "errors"

var (
	ErrKnowledgeNotFound     = errors.New("资源不存在")
	ErrKnowledgeForbidden    = errors.New("无权执行此操作")
	ErrWorkspaceMemberExists = errors.New("用户已是工作空间成员")
	ErrWorkspaceOwner        = errors.New("不能修改或移除工作空间所有者")
	ErrCatalogCycle          = errors.New("不能将目录移动到自身或其子目录下")
	ErrShareDisabled         = errors.New("分享链接已被作者关闭")
	ErrShareExpired          = errors.New("分享链接已过期")
	ErrKnowledgeInvalid      = errors.New("请求参数无效")
	ErrPublicWikiTooLarge    = errors.New("公开知识库节点超过2000个，请拆分工作区")
)
