package dependencytree

import (
	"errors"
	"fmt"
	"strings"
)

type DependencyTreeItem[T interface{}] struct {
	ID            string
	Name          string
	FlatIndex     int
	index         int
	highest       int
	lowest        int
	isDependentOn []string
	parentName    string
	Parent        *DependencyTreeItem[T]
	obj           T
	requiredBy    []string
	Children      []*DependencyTreeItem[T]
	CallBack      func()
	Metadata      map[string]interface{}
}

func NewDependencyTreeItem[T interface{}](id string, name string, value T) (*DependencyTreeItem[T], error) {
	if id == "" {
		return nil, errors.New("id must not be empty")
	}

	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	result := DependencyTreeItem[T]{
		ID:            id,
		Name:          name,
		index:         0,
		FlatIndex:     0,
		highest:       0,
		lowest:        0,
		isDependentOn: []string{},
		parentName:    "root",
		Parent:        nil,
		obj:           value,
		requiredBy:    []string{},
		Children:      []*DependencyTreeItem[T]{},
		CallBack:      nil,
		Metadata:      make(map[string]interface{}),
	}

	return &result, nil
}

func (dt *DependencyTreeItem[T]) AddChild(child *DependencyTreeItem[T]) {
	dt.Children = append(dt.Children, child)
	dt.AddRequiredBy(child.ID)
}

func (dt *DependencyTreeItem[T]) SetParent(parent string) {
	dt.parentName = parent
}

func (dt *DependencyTreeItem[T]) DependsOn(idOrName string) error {
	for _, item := range dt.isDependentOn {
		if strings.EqualFold(item, idOrName) {
			return fmt.Errorf("item %v already exists", idOrName)
		}
	}

	dt.isDependentOn = append(dt.isDependentOn, idOrName)
	return nil
}

func (dt *DependencyTreeItem[T]) IsDependentOn() []string {
	return dt.isDependentOn
}

func (dt *DependencyTreeItem[T]) GetParentName() string {
	if dt.Parent != nil {
		return dt.Parent.Name
	}

	return dt.parentName
}

func (dt *DependencyTreeItem[T]) GetParentId() string {
	if dt.Parent != nil {
		return dt.Parent.ID
	}

	return ""
}

func (dt *DependencyTreeItem[T]) RequiredBy() []string {
	return dt.requiredBy
}

func (dt *DependencyTreeItem[T]) AddRequiredBy(id string) {
	for _, item := range dt.requiredBy {
		if strings.EqualFold(item, id) {
			return
		}
	}

	dt.requiredBy = append(dt.requiredBy, id)
}

func (p *DependencyTreeItem[T]) GetProperty(key string, defaultValue interface{}) interface{} {
	if p.Metadata == nil {
		return defaultValue
	}

	value, ok := p.Metadata[key]
	if !ok {
		return defaultValue
	}

	return value
}

func (p *DependencyTreeItem[T]) SetProperty(key string, value interface{}) error {
	if p.Metadata == nil {
		p.Metadata = make(map[string]interface{})
	}

	p.Metadata[key] = value
	return nil
}
