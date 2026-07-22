package service

import (
	"errors"
	"regexp"
	"sort"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

const publicWikiNodeLimit = 2000

var publicWikiHTMLPolicy = newPublicWikiHTMLPolicy()

type PublicWikiService struct{}

func NewPublicWikiService() *PublicWikiService { return &PublicWikiService{} }

func (s *PublicWikiService) Stats() (*models.PublicWikiStatsResponse, error) {
	var count int64
	err := database.DB.Model(&models.Doc{}).
		Joins("JOIN workspaces ON workspaces.id = docs.workspace_id AND workspaces.owner_id = docs.owner_id AND workspaces.is_public = ? AND workspaces.deleted_at IS NULL", true).
		Where("docs.status = ? AND docs.deleted_at IS NULL", models.DocStatusPublished).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	return &models.PublicWikiStatsResponse{PublicDocCount: count}, nil
}

func (s *PublicWikiService) Workspaces(page, pageSize int) ([]*models.PublicWorkspaceResponse, int64, error) {
	var total int64
	if err := database.DB.Model(&models.Workspace{}).
		Where("workspaces.is_public = ? AND workspaces.deleted_at IS NULL", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var workspaces []*models.PublicWorkspaceResponse
	err := database.DB.Model(&models.Workspace{}).
		Select("workspaces.id, workspaces.name, workspaces.description, workspaces.icon, workspaces.updated_at, workspaces.owner_id AS author_id, COALESCE(NULLIF(users.nickname, ''), users.username, '') AS author_name, COALESCE(users.avatar, '') AS author_avatar, COUNT(docs.id) AS doc_count").
		Joins("LEFT JOIN docs ON docs.workspace_id = workspaces.id AND docs.owner_id = workspaces.owner_id AND docs.status = ? AND docs.deleted_at IS NULL", models.DocStatusPublished).
		Joins("LEFT JOIN users ON users.id = workspaces.owner_id AND users.deleted_at IS NULL").
		Where("workspaces.is_public = ? AND workspaces.deleted_at IS NULL", true).
		Group("workspaces.id, workspaces.name, workspaces.description, workspaces.icon, workspaces.updated_at, workspaces.owner_id, users.id, users.nickname, users.username, users.avatar").
		Order("workspaces.updated_at DESC, workspaces.id DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Scan(&workspaces).Error
	return workspaces, total, err
}

func (s *PublicWikiService) Tree(workspaceID uint) (*models.PublicWorkspaceTreeResponse, error) {
	workspace, err := publicWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}

	var docs []models.Doc
	err = database.DB.Model(&models.Doc{}).
		Select("docs.id, docs.catalog_id, docs.title, docs.sort, docs.published_at, docs.updated_at").
		Joins("JOIN workspaces ON workspaces.id = docs.workspace_id AND workspaces.owner_id = docs.owner_id AND workspaces.is_public = ? AND workspaces.deleted_at IS NULL", true).
		Joins("LEFT JOIN catalogs ON catalogs.id = docs.catalog_id AND catalogs.workspace_id = docs.workspace_id AND catalogs.owner_id = docs.owner_id AND catalogs.deleted_at IS NULL").
		Where("docs.workspace_id = ? AND docs.owner_id = ? AND docs.status = ? AND docs.deleted_at IS NULL AND (docs.catalog_id IS NULL OR catalogs.id IS NOT NULL)", workspace.ID, workspace.OwnerID, models.DocStatusPublished).
		Order("docs.catalog_id IS NOT NULL ASC, docs.catalog_id ASC, docs.sort ASC, docs.id ASC").
		Limit(publicWikiNodeLimit + 1).Find(&docs).Error
	if err != nil {
		return nil, err
	}
	if len(docs) > publicWikiNodeLimit {
		return nil, ErrPublicWikiTooLarge
	}

	catalogs, err := publicWikiCatalogs(workspace.ID, workspace.OwnerID, docs)
	if err != nil {
		return nil, err
	}
	publicCatalogs := buildPublicCatalogTree(catalogs, docs)
	if len(docs)+countPublicCatalogs(publicCatalogs) > publicWikiNodeLimit {
		return nil, ErrPublicWikiTooLarge
	}

	publicDocs := make([]*models.PublicDocTreeResponse, 0, len(docs))
	for i := range docs {
		publicDocs = append(publicDocs, &models.PublicDocTreeResponse{
			ID: docs[i].ID, CatalogID: docs[i].CatalogID, Title: docs[i].Title, Sort: docs[i].Sort,
			PublishedAt: docs[i].PublishedAt, UpdatedAt: docs[i].UpdatedAt,
		})
	}
	return &models.PublicWorkspaceTreeResponse{
		Workspace: &models.PublicWorkspaceSummary{ID: workspace.ID, Name: workspace.Name, Description: workspace.Description, Icon: workspace.Icon},
		Catalogs:  publicCatalogs,
		Docs:      publicDocs,
	}, nil
}

func publicWikiCatalogs(workspaceID, ownerID uint, docs []models.Doc) ([]models.Catalog, error) {
	ids := make([]uint, 0)
	seen := make(map[uint]bool)
	for i := range docs {
		if docs[i].CatalogID != nil && !seen[*docs[i].CatalogID] {
			seen[*docs[i].CatalogID] = true
			ids = append(ids, *docs[i].CatalogID)
		}
	}

	catalogs := make([]models.Catalog, 0, len(ids))
	for len(ids) > 0 {
		if len(docs)+len(catalogs)+len(ids) > publicWikiNodeLimit {
			return nil, ErrPublicWikiTooLarge
		}
		var level []models.Catalog
		if err := database.DB.Select("id, parent_id, name, sort").
			Where("workspace_id = ? AND owner_id = ? AND deleted_at IS NULL AND id IN ?", workspaceID, ownerID, ids).
			Find(&level).Error; err != nil {
			return nil, err
		}
		catalogs = append(catalogs, level...)
		next := make([]uint, 0)
		for i := range level {
			if level[i].ParentID != nil && !seen[*level[i].ParentID] {
				seen[*level[i].ParentID] = true
				next = append(next, *level[i].ParentID)
			}
		}
		ids = next
	}
	return catalogs, nil
}

func (s *PublicWikiService) Doc(id uint) (*models.PublicDocResponse, error) {
	var doc models.PublicDocResponse
	err := database.DB.Model(&models.Doc{}).
		Select("docs.id, docs.workspace_id, docs.catalog_id, docs.title, docs.content_html, docs.view_count, docs.published_at, docs.updated_at").
		Joins("JOIN workspaces ON workspaces.id = docs.workspace_id AND workspaces.owner_id = docs.owner_id AND workspaces.is_public = ? AND workspaces.deleted_at IS NULL", true).
		Joins("LEFT JOIN catalogs ON catalogs.id = docs.catalog_id AND catalogs.workspace_id = docs.workspace_id AND catalogs.owner_id = docs.owner_id AND catalogs.deleted_at IS NULL").
		Where("docs.id = ? AND docs.status = ? AND docs.deleted_at IS NULL AND (docs.catalog_id IS NULL OR catalogs.id IS NOT NULL)", id, models.DocStatusPublished).Take(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	doc.ContentHTML = sanitizePublicWikiHTML(doc.ContentHTML)
	return &doc, nil
}

func publicWorkspace(id uint) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := database.DB.Where("id = ? AND is_public = ? AND deleted_at IS NULL", id, true).First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKnowledgeNotFound
		}
		return nil, err
	}
	return &workspace, nil
}

func buildPublicCatalogTree(catalogs []models.Catalog, docs []models.Doc) []*models.PublicCatalogResponse {
	catalogByID := make(map[uint]models.Catalog, len(catalogs))
	for i := range catalogs {
		catalogByID[catalogs[i].ID] = catalogs[i]
	}
	include := make(map[uint]bool)
	for i := range docs {
		if docs[i].CatalogID == nil {
			continue
		}
		id := *docs[i].CatalogID
		path := make(map[uint]bool)
		for id != 0 && !path[id] {
			catalog, ok := catalogByID[id]
			if !ok {
				break
			}
			path[id] = true
			include[id] = true
			if catalog.ParentID == nil {
				break
			}
			id = *catalog.ParentID
		}
	}

	filtered := make([]models.Catalog, 0, len(include))
	for i := range catalogs {
		if include[catalogs[i].ID] {
			filtered = append(filtered, catalogs[i])
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].Sort == filtered[j].Sort {
			return filtered[i].ID < filtered[j].ID
		}
		return filtered[i].Sort < filtered[j].Sort
	})

	nodes := make(map[uint]*models.PublicCatalogResponse, len(filtered))
	for i := range filtered {
		nodes[filtered[i].ID] = &models.PublicCatalogResponse{
			ID: filtered[i].ID, ParentID: filtered[i].ParentID, Name: filtered[i].Name,
			Sort: filtered[i].Sort, Children: make([]*models.PublicCatalogResponse, 0),
		}
	}
	roots := make([]*models.PublicCatalogResponse, 0)
	for i := range filtered {
		node := nodes[filtered[i].ID]
		if filtered[i].ParentID != nil {
			if parent, ok := nodes[*filtered[i].ParentID]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots
}

func countPublicCatalogs(catalogs []*models.PublicCatalogResponse) int {
	count := 0
	stack := append([]*models.PublicCatalogResponse(nil), catalogs...)
	for len(stack) > 0 {
		last := len(stack) - 1
		node := stack[last]
		stack = stack[:last]
		count++
		stack = append(stack, node.Children...)
	}
	return count
}

func sanitizePublicWikiHTML(content string) string {
	return publicWikiHTMLPolicy.Sanitize(content)
}

func newPublicWikiHTMLPolicy() *bluemonday.Policy {
	policy := bluemonday.UGCPolicy()
	policy.AllowElements("table", "thead", "tbody", "tfoot", "tr", "th", "td", "caption", "colgroup", "col")
	policy.AllowAttrs("colspan", "rowspan", "scope").OnElements("th", "td")
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^[a-zA-Z0-9_ -]{1,200}$`)).OnElements("pre", "code")
	policy.AllowAttrs("loading").Matching(regexp.MustCompile(`^(lazy|eager)$`)).OnElements("img")
	return policy
}
