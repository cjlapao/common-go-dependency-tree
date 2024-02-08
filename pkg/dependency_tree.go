package dependency_tree

import (
	"fmt"
	"strconv"
)

func (d *DependencyTreeService[T]) Build() ([]*DependencyTreeItem[T], error) {
	if d.IsDebug() {
		d.logger.Debug("Dependency Tree Before:")
		for idx, service := range d.flatTree {
			d.logger.Debug("[%s] %s", strconv.Itoa(idx), service.Name)
		}
	}

	// Initial Pass to flatten our dependency tree based on linear dependency
	values, err := d.buildRightDependency()
	if err != nil {
		return values, err
	}

	// Reordering dependencies for stragglers (left shift)
	// First pass we ordered all of the dependencies in a forward method
	// putting it all in an order of a -> b -> c kind of dependency
	// we also might have services that might need sifting as they didn't fall into this like
	// a late service that has only dependency on b but not on c
	// }

	values, err = d.buildLeftDependency()
	if err != nil {
		return values, err
	}

	// Building parent dependency
	err = d.buildChildDependency()
	if err != nil {
		return values, err
	}

	// Last pass to make sure the left dependency did not make further issues
	values, err = d.buildRightDependency()
	if err != nil {
		return values, err
	}

	tree := d.buildTree("root")
	d.tree = tree

	if d.IsDebug() && d.IsVerbose() {
		d.logger.Debug(d.String())
	}

	if d.IsDebug() {
		d.logger.Debug("Dependency Tree After:")
		for idx, service := range values {
			d.logger.Debug("[%s] %s", strconv.Itoa(idx), service.Name)
		}
	}

	return values, nil
}

func (d *DependencyTreeService[T]) buildTree(parent string) []*DependencyTreeItem[T] {
	result := []*DependencyTreeItem[T]{}
	childs := d.GetItemByParent(parent)
	for _, child := range childs {
		if len(child.Children) > 0 {
			children := d.buildTree(child.Name)
			child.Children = children
		}
		result = append(result, child)
	}

	return result
}

func (d *DependencyTreeService[T]) buildChildDependency() error {
	idx := 0
parentLoop:
	for {
		shiftedItem := false
	treeLoop:
		for idx, item := range d.flatTree {
			if len(item.Children) == 0 {
				continue treeLoop
			}

			_, err := d.shiftChildItems(idx)
			if err != nil {
				return err
			}
		}

		idx += 1
		if !shiftedItem || idx > 1000 {
			if idx == 1000 {
				err := fmt.Errorf("something went wrong and we shifted more than 1000 items")
				return err
			}

			break parentLoop
		}
	}

	return nil
}

func (d *DependencyTreeService[T]) buildRightDependency() ([]*DependencyTreeItem[T], error) {
	idx := 0
outerLoop:
	for {
		needsShifting := false
	treeLoop:
		for svcIndex, item := range d.flatTree {
			dependencies := item.IsDependentOn()
		dependenciesLoop:
			for _, dependency := range dependencies {
				dependencyIndex, err := d.GetItemIndex(dependency)
				if err != nil {
					return nil, err
				}

				if dependencyIndex < 0 {
					err := fmt.Errorf("dependency on %s of service %s was not found in the context configuration", dependency, item.Name)
					return nil, err
				}

				if svcIndex < dependencyIndex {
					needsShifting = true
					if d.IsDebug() && d.IsVerbose() {
						d.logger.Debug("Shifting %s on index %s to index %s", item.Name, strconv.Itoa(svcIndex), strconv.Itoa(dependencyIndex))
					}

					// shiftRight(dependencyIndex - 1)
					_, err := d.shiftTo(svcIndex, dependencyIndex)
					if err != nil {
						return nil, err
					}
					// shifting all of the children to the same position as the parent
					children := d.GetItemChildren(item.ID)
					if len(children) == 0 {
						continue dependenciesLoop
					}

					currentIndex := 0
				childrenLoop:
					for {
						childIndex, err := d.GetItemIndex(children[currentIndex].ID)
						if err != nil {
							return nil, err
						}

						if childIndex < 0 {
							err := fmt.Errorf("dependency on %s of service %s was not found in the context configuration", children[currentIndex].Name, item.Name)
							return nil, err
						}

						if d.IsDebug() && d.IsVerbose() {
							d.logger.Debug("Shifting %s on index %s to index %s", children[currentIndex].Name, strconv.Itoa(childIndex), strconv.Itoa(dependencyIndex))
						}

						_, err = d.shiftTo(childIndex, dependencyIndex)
						if err != nil {
							return nil, err
						}

						currentIndex += 1
						if currentIndex == len(children) {
							break childrenLoop
						}
					}
				}
			}

			if needsShifting {
				break treeLoop
			}
		}

		idx += 1
		if !needsShifting || idx > 1000 {
			if idx == 1000 {
				err := fmt.Errorf("something went wrong and we shifted more than 1000 items")
				return nil, err
			}

			break outerLoop
		}
	}

	return d.flatTree, nil
}

func (d *DependencyTreeService[T]) buildLeftDependency() ([]*DependencyTreeItem[T], error) {
	idx := 0
	for {
		shiftHappened := false
		// For each loop after a shift we need to re-update our tree so we know of the position shifting
		tree, err := d.updateTree()
		if err != nil {
			return nil, err
		}
		for _, item := range d.flatTree {
			treeItem := tree[item.Name]
			highestIndex := -1
			changeTo := ""
			change := true

			if treeItem.highest == -1 || treeItem.lowest == -1 {
				continue
			}

			if treeItem.index < treeItem.lowest || treeItem.index-1 > treeItem.highest {
				for k, kv := range tree {
					if kv.highest == treeItem.highest {
						if kv.index > highestIndex && kv.index != treeItem.index {
							highestIndex = kv.index
							changeTo = k
						}

						if kv.index == treeItem.index {
							if treeItem.index-1 == highestIndex {
								change = false
								break
							}
						}
					}
				}

				if change && highestIndex != -1 && treeItem.index > highestIndex && treeItem.index != highestIndex+1 {
					if d.IsDebug() && d.IsVerbose() {
						d.logger.Debug("[%s] %s should move to %s after %s", fmt.Sprintf("%d", treeItem.index), item.Name, fmt.Sprintf("%d", highestIndex+1), changeTo)
					}
					_, err := d.shiftTo(treeItem.index, highestIndex+1)
					if err != nil {
						return nil, err
					}
					shiftHappened = true
					break
				}
			}
		}

		idx += 1
		if !shiftHappened || idx > 1000 {
			if idx == 1000 {
				err := fmt.Errorf("something went wrong and we shifted more than 1000 items")
				return nil, err
			}
			break
		}
	}

	return d.flatTree, nil
}

func (d *DependencyTreeService[T]) updateTree() (map[string]*DependencyTreeItem[T], error) {
	result := make(map[string]*DependencyTreeItem[T])
	for i, item := range d.flatTree {
		item.index = i
		item.highest = -1
		item.lowest = -1

		for _, dependency := range item.IsDependentOn() {
			dependencyIndex, err := d.GetItemIndex(dependency)
			if err != nil {
				return nil, err
			}

			if dependencyIndex > item.highest {
				item.highest = dependencyIndex
				if item.lowest == -1 {
					item.lowest = dependencyIndex
				}
			}
			if dependencyIndex < item.highest && dependencyIndex > item.lowest {
				item.lowest = dependencyIndex
			}
		}

		if d.IsDebug() && d.IsVerbose() {
			d.logger.Debug("%s [%s] %s -> HighestDependency: %s | LowestDependency: %s", fmt.Sprintf("%d", i), item.Name, fmt.Sprintf("%d", item.highest), fmt.Sprintf("%d", item.lowest))
		}

		result[item.Name] = item
	}

	return result, nil
}

func (d *DependencyTreeService[T]) shiftChildItems(index int) (bool, error) {
	currentChildIndex := 0
	didShuffle := false
	offsetParentIndex := index
	for {
		child := d.flatTree[index].Children[currentChildIndex]
		shiftedItem := false

		currentIdx, err := d.GetItemIndex(child.Name)
		if err != nil {
			return false, err
		}

		offset := offsetParentIndex + currentChildIndex + 1

		if currentIdx != offset {
			d.logger.Debug("Parent: %v, Child: %v\n", d.flatTree[index].Name, child.Name)
			shiftedItem = true
			if d.IsDebug() && d.IsVerbose() {
				d.logger.Debug("Shifting %s child %v on index %s to index %s", d.flatTree[index].Name, child.Name, strconv.Itoa(currentIdx), strconv.Itoa(offset))
			}

			_, err := d.shiftTo(currentIdx, offset)
			if err != nil {
				return false, err
			}

			if len(d.flatTree[offset].Children) > 0 {
				shuffledChilds, err := d.shiftChildItems(offset)
				if err != nil {
					return false, err
				}
				if shuffledChilds {
					offsetParentIndex = offsetParentIndex + len(d.flatTree[offset].Children)
					didShuffle = true
				}
			}
		}
		currentChildIndex += 1

		if shiftedItem {
			didShuffle = true
			if child.CallBack != nil {
				child.CallBack()
			}
		}

		if currentChildIndex >= len(d.flatTree[index].Children) {
			break
		}
	}

	return didShuffle, nil
}
