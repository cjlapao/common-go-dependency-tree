package dependencytree

import (
	"strconv"
	"strings"
)

func (d *DependencyTreeService[T]) shiftTo(from, to int) ([]*DependencyTreeItem[T], error) {
	switch size := len(d.flatTree); {
	case size == 0:
		// Nothing to shift, returning
		return d.flatTree, nil
	case size > 1:
		// not inside the range of the array
		if from > len(d.flatTree) || from == -1 || to > len(d.flatTree) || to == -1 {
			return d.flatTree, nil
		}

		// Deciding the direction of the shift
		pos := from - to
		switch itemPosition := pos; {
		case itemPosition < 0:
			if d.IsDebug() && d.IsVerbose() {
				d.logger.Debug("Shifting forwards from %s to %s", strconv.Itoa(from), strconv.Itoa(to))
			}
			for {
				if from == to {
					break
				}
				d.flatTree[from], d.flatTree[from+1] = d.flatTree[from+1], d.flatTree[from]

				from += 1
			}
		case itemPosition > 0:
			if d.IsDebug() && d.IsVerbose() {
				d.logger.Debug("Shifting backwards from %s to %s", strconv.Itoa(from), strconv.Itoa(to))
			}
			for {
				if from == to {
					break
				}
				d.flatTree[from], d.flatTree[from-1] = d.flatTree[from-1], d.flatTree[from]
				from -= 1
			}
		}
	}

	return d.flatTree, nil
}

func (d *DependencyTreeService[T]) printTree(tree []*DependencyTreeItem[T], level int, prefix string) []string {
	lines := []string{}
	spacer := ""
	if level > 0 {
		spacer = prefix
	}

	for idx, item := range tree {
		msg := spacer
		if level > 0 && !strings.HasPrefix(msg, "|") {
			msg += " "
		}

		if idx == 0 {
			if level > 0 {
				if len(tree) == 1 {
					msg += "└─ "
					prefix += "|  "
				} else {
					msg += "├─ "
					prefix += "|  "
				}
			} else {
				msg += "┌─ "
				prefix += "|  "
			}
		} else if idx == len(tree)-1 {
			msg += "└─ "
			if len(item.Children) > 0 {
				if level > 0 {
					prefix = strings.TrimSpace(prefix)
					prefix = prefix[:len(prefix)-2]
					prefix = prefix + strings.Repeat("  ", level+1)
				} else {
					prefix = strings.Repeat("  ", level+1)
				}
			} else {
				prefix = strings.Repeat("  ", level+1)
			}
		} else {
			msg += "├─ "
		}

		msg += item.Name

		lines = append(lines, msg)
		childLines := d.printTree(item.Children, level+1, prefix)
		lines = append(lines, childLines...)
	}

	return lines
}
