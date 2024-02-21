package dependencytree

import (
	"errors"
	"strconv"
	"strings"
)

const firstLine = "┌─ "
const line = "|  "
const middleLine = "├─ "
const lastLine = "└─ "

func (d *DependencyTreeService[T]) printVerbosef(format string, args ...interface{}) {
	if d.IsDebug() && d.IsVerbose() {
		d.logger.Debug(format, args...)
	}
}

func (d *DependencyTreeService[T]) shiftTo(from, to int) ([]*DependencyTreeItem[T], error) {
	if from > len(d.flatTree) || from == -1 || to > len(d.flatTree) || to == -1 {
		return d.flatTree, errors.New("from or to index is out of range")
	}

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
			d.moveForwards(from, to)
		case itemPosition > 0:
			d.moveBackwards(from, to)
		}
	}

	return d.flatTree, nil
}

func (d *DependencyTreeService[T]) moveBackwards(from, to int) {
	d.printVerbosef("Shifting backwards from %s to %s", strconv.Itoa(from), strconv.Itoa(to))
	for {
		if from == to {
			break
		}
		d.flatTree[from], d.flatTree[from-1] = d.flatTree[from-1], d.flatTree[from]
		from -= 1
	}
}

func (d *DependencyTreeService[T]) moveForwards(from, to int) {
	d.printVerbosef("Shifting forwards from %s to %s", strconv.Itoa(from), strconv.Itoa(to))
	for {
		if from == to {
			break
		}
		d.flatTree[from], d.flatTree[from+1] = d.flatTree[from+1], d.flatTree[from]

		from += 1
	}
}

func (d *DependencyTreeService[T]) printTree(tree []*DependencyTreeItem[T], level int, prefix string) []string {
	lines := []string{}
	spacer := ""
	if level > 0 {
		spacer = prefix
	}

	if level == 0 && len(tree) == 1 {
		msg := ""
		msg, prefix = d.printLastItem(tree[0], level, msg, prefix)
		msg += tree[0].Name
		lines = append(lines, msg)
		childLines := d.printTree(tree[0].Children, level+1, prefix)
		lines = append(lines, childLines...)
		return lines
	}

	for idx, item := range tree {
		msg := spacer
		if level > 0 {

			if strings.HasPrefix(strings.TrimSpace(msg), "|") {
				msg += " "
			}
			if !strings.HasPrefix(msg, "|") {
				msg += " "
			}
		}

		if idx == 0 {
			msg, prefix = d.printStartItem(tree, level, msg, prefix)
		} else if idx == len(tree)-1 {
			msg, prefix = d.printLastItem(item, level, msg, prefix)
		} else {
			msg += middleLine
		}

		msg += item.Name

		lines = append(lines, msg)
		childLines := d.printTree(item.Children, level+1, prefix)
		lines = append(lines, childLines...)
	}

	return lines
}

func (d *DependencyTreeService[T]) printStartItem(tree []*DependencyTreeItem[T], level int, msg, prefix string) (string, string) {
	if level > 0 {
		if len(tree) == 1 {
			msg += lastLine
			prefix += line
		} else {
			msg += middleLine
			prefix += line
		}
	} else {
		msg += firstLine
		prefix += line
	}

	return msg, prefix
}

func (d *DependencyTreeService[T]) printLastItem(item *DependencyTreeItem[T], level int, msg, prefix string) (string, string) {
	msg += lastLine
	if len(item.Children) > 0 {
		if level > 0 {
			prefix = strings.TrimSpace(prefix)
			if prefix == "|" {
				prefix = ""
			} else {
				prefix = prefix[:len(prefix)-2]
			}
			prefix = prefix + strings.Repeat("  ", level+1)
		} else {
			prefix = strings.Repeat("  ", level+1)
		}
	} else {
		prefix = strings.Repeat("  ", level+1)
	}

	return msg, prefix
}
