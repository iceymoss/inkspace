package service

import (
	"errors"
	"fmt"
)

var (
	ErrKnowledgeNotFound     = errors.New("资源不存在")
	ErrKnowledgeForbidden    = errors.New("无权执行此操作")
	ErrWorkspaceMemberExists = errors.New("用户已是工作空间成员")
	ErrWorkspaceOwner        = errors.New("不能修改或移除工作空间所有者")
	ErrCatalogCycle          = errors.New("不能将目录移动到自身或其子目录下")
	ErrShareDisabled         = errors.New("分享链接已被作者关闭")
	ErrShareExpired          = errors.New("分享链接已过期")
	ErrKnowledgeInvalid      = errors.New("请求参数无效")
	ErrDocRevisionConflict   = errors.New("文档已被其他请求更新")
	ErrDocNotEditable        = errors.New("该文档类型不支持在线编辑")
	ErrDocPreviewUnavailable = errors.New("该文件暂不支持在线预览")
	ErrKnowledgeFileTooLarge = errors.New("文件大小超过限制")
	ErrKnowledgeTextTooLarge = errors.New("在线编辑内容不能超过 2 MiB")
	ErrKnowledgeFileType     = errors.New("不支持或文件类型不匹配")
	ErrPublicWikiTooLarge    = errors.New("公开知识库节点超过2000个，请拆分工作区")
)

type DocRevisionConflictError struct {
	Revision uint64
}

func (e *DocRevisionConflictError) Error() string {
	return fmt.Sprintf("%s，当前版本为 %d", ErrDocRevisionConflict, e.Revision)
}

func (e *DocRevisionConflictError) Unwrap() error { return ErrDocRevisionConflict }
