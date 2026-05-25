package dto

import "product-mall/internal/model"

type Category struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	ParentID uint        `json:"parent_id"`
	Level    uint        `json:"level"`
	Sort     int         `json:"sort"`
	Children []*Category `json:"children,omitempty"`
}

func BuildCategory(item model.Category) Category {
	return Category{
		ID:       item.ID,
		Name:     item.Name,
		ParentID: item.ParentID,
		Level:    item.Level,
		Sort:     item.Sort,
	}
}

func BuildCategories(items []model.Category) []Category {
	categories := make([]Category, 0, len(items))
	for _, item := range items {
		categories = append(categories, BuildCategory(item))
	}
	return categories
}

func BuildCategoryTree(items []model.Category) []*Category {
	categoryMap := make(map[uint]*Category)
	var roots []*Category

	for _, item := range items {
		c := BuildCategory(item)
		categoryMap[c.ID] = &c
	}

	for _, c := range categoryMap {
		if c.ParentID == 0 {
			roots = append(roots, c)
		} else if parent, ok := categoryMap[c.ParentID]; ok {
			parent.Children = append(parent.Children, c)
		}
	}
	return roots
}
