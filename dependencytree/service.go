package dependencytree

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	log "github.com/cjlapao/common-go-logger"
)

var (
	globalDependencyTreeService []interface{}
	lock                        = &sync.Mutex{}
)

type DependencyTreeService[T interface{}] struct {
	logger   *log.LoggerService
	debug    bool
	verbose  bool
	flatTree []*DependencyTreeItem[T]
	tree     []*DependencyTreeItem[T]
}

func Get[T interface{}](v T) *DependencyTreeService[T] {
	lock.Lock()
	for _, item := range globalDependencyTreeService {
		mPointer := &DependencyTreeService[T]{}
		itemKind := reflect.TypeOf(item).String()
		pointerKind := reflect.TypeOf(mPointer).String()
		if pointerKind == itemKind {
			lock.Unlock()
			return item.(*DependencyTreeService[T])
		}
	}

	newTreeType := DependencyTreeService[T]{
		debug:    false,
		verbose:  false,
		logger:   log.Get(),
		flatTree: []*DependencyTreeItem[T]{},
		tree:     []*DependencyTreeItem[T]{},
	}
	globalDependencyTreeService = append(globalDependencyTreeService, &newTreeType)

	lock.Unlock()

	return &newTreeType
}

func (d *DependencyTreeService[T]) String() string {
	lines := d.printTree(d.tree, 0, "", false)
	return strings.Join(lines, "\n")
}

func (d *DependencyTreeService[T]) Clear() {
	d.flatTree = []*DependencyTreeItem[T]{}
	d.tree = []*DependencyTreeItem[T]{}
}

func (d *DependencyTreeService[T]) IsDebug() bool {
	return d.debug
}

func (d *DependencyTreeService[T]) IsVerbose() bool {
	return d.verbose
}

func (d *DependencyTreeService[T]) SetDebug(debug bool) {
	d.debug = debug
}

func (d *DependencyTreeService[T]) SetVerbose(verbose bool) {
	d.verbose = verbose
}

func (d *DependencyTreeService[T]) SetLogger(logger *log.LoggerService) {
	d.logger = logger
}

func (d *DependencyTreeService[T]) AddRootItem(id string, name string, value T) (*DependencyTreeItem[T], error) {
	teeItem, err := NewDependencyTreeItem[T](id, name, value)
	if err != nil {
		return nil, err
	}
	teeItem.SetParent("root")

	err = d.AddDependencyTreeItem(teeItem)
	if err != nil {
		return nil, err
	}

	return teeItem, nil
}

func (d *DependencyTreeService[T]) DependsOn(id string, dependencyId string) error {
	item := d.GetItem(id)
	if item == nil {
		return fmt.Errorf("item %v not found", id)
	}

	dependency := d.GetItem(dependencyId)
	if dependency == nil {
		return fmt.Errorf("dependency %v not found", dependencyId)
	}

	item.isDependentOn = append(item.isDependentOn, dependency.ID)
	dependency.requiredBy = append(dependency.requiredBy, item.ID)

	return nil
}

func (d *DependencyTreeService[T]) AddItem(id string, name string, parent string, value T) (*DependencyTreeItem[T], error) {
	treeItem, err := NewDependencyTreeItem[T](id, name, value)
	if err != nil {
		return nil, err
	}

	treeItem.SetParent(parent)

	err = d.AddDependencyTreeItem(treeItem)
	if err != nil {
		return nil, err
	}

	return treeItem, nil
}

func (d *DependencyTreeService[T]) AddDependencyTreeItem(item *DependencyTreeItem[T]) error {
	for _, i := range d.flatTree {
		if strings.EqualFold(i.ID, item.ID) || strings.EqualFold(i.Name, item.Name) {
			return fmt.Errorf("item with id %v already exists", item.ID)
		}
	}

	d.flatTree = append(d.flatTree, item)
	return nil
}

func (d *DependencyTreeService[T]) RemoveDependencyTreeItem(item *DependencyTreeItem[T]) error {
	for idx, i := range d.flatTree {
		if strings.EqualFold(i.ID, item.ID) || strings.EqualFold(i.Name, item.Name) {
			d.flatTree = append(d.flatTree[:idx], d.flatTree[idx+1:]...)
			return nil
		}
	}

	return fmt.Errorf("item with id %v not found", item.ID)
}

func (d *DependencyTreeService[T]) GetItem(nameOrId string) *DependencyTreeItem[T] {
	for _, item := range d.flatTree {
		if strings.EqualFold(item.ID, nameOrId) || strings.EqualFold(item.Name, nameOrId) {
			return item
		}
	}

	return nil
}

func (d *DependencyTreeService[T]) GetItemIndex(nameOrId string) (int, error) {
	for idx, item := range d.flatTree {
		if strings.EqualFold(item.ID, nameOrId) || strings.EqualFold(item.Name, nameOrId) {
			return idx, nil
		}
	}

	return -1, fmt.Errorf("item with id %v not found", nameOrId)
}

func (d *DependencyTreeService[T]) GetItemChildren(nameOrId string) []*DependencyTreeItem[T] {
	for _, item := range d.flatTree {
		if strings.EqualFold(item.ID, nameOrId) || strings.EqualFold(item.Name, nameOrId) {
			return item.Children
		}
	}

	return []*DependencyTreeItem[T]{}
}

func (d *DependencyTreeService[T]) GetItemByParent(parent string) []*DependencyTreeItem[T] {
	result := []*DependencyTreeItem[T]{}

	for _, item := range d.flatTree {
		if strings.EqualFold(item.GetParentName(), parent) || strings.EqualFold(item.GetParentId(), parent) {
			result = append(result, item)
		}
	}

	return result
}

func (d *DependencyTreeService[T]) GetItemDependencies(nameOrId string) []*DependencyTreeItem[T] {
	for _, item := range d.flatTree {
		if strings.EqualFold(item.ID, nameOrId) || strings.EqualFold(item.Name, nameOrId) {
			return item.Children
		}
	}

	return []*DependencyTreeItem[T]{}
}

func (d *DependencyTreeService[T]) Tree() []*DependencyTreeItem[T] {
	buildTree := d.buildTree("root")
	d.tree = buildTree
	return d.tree
}

func (d *DependencyTreeService[T]) FlatTree() []*DependencyTreeItem[T] {
	return d.flatTree
}

func (d *DependencyTreeService[T]) PrintFlatTree() {
	for _, item := range d.flatTree {
		d.logger.Info("Id: %v, Name: %v", item.ID, item.Name)
	}
}
