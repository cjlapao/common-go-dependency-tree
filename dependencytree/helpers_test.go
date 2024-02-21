package dependencytree

import (
	"strings"
	"testing"

	log "github.com/cjlapao/common-go-logger"
	"github.com/stretchr/testify/assert"
)

var expectedSimpleString = `└─ item 1`

var expectedSimpleTreeWithChildrenString = `└─ item 2
   ├─ item 2 Child 2
   ├─ item 2 Child 1
   |  ├─ item 2 Child 1 Child 2
   |  └─ item 2 Child 1 Child 1
   └─ item 2 Child 3
      ├─ item 2 Child 3 Child 2
      └─ item 2 Child 3 Child 1`

var expectedComplexString = `┌─ item 1
├─ item 2
|  ├─ item 2 Child 2
|  ├─ item 2 Child 1
|  |  ├─ item 2 Child 1 Child 2
|  |  └─ item 2 Child 1 Child 1
|  └─ item 2 Child 3
|     ├─ item 2 Child 3 Child 2
|     └─ item 2 Child 3 Child 1
├─ item 3
├─ item 4
|  └─ item 4 Child 1
├─ item 5
└─ item 6
   ├─ item 6 Child 2
   └─ item 6 Child 1`

func TestPrintSimpleTree(t *testing.T) {
	logger := log.Get()
	logger.LogLevel = log.Debug
	dpService := Get(MockObject1{})
	dpService.SetLogger(logger)

	_, _ = dpService.AddRootItem("item_1", "item 1", MockObject1{id: "item_1", someStoredValue: "item 1"})

	_, err := dpService.Build()
	assert.NoError(t, err)

	lines := dpService.printTree(dpService.tree, 0, "", true)
	assert.Equal(t, expectedSimpleString, strings.Join(lines, "\n"))
}

func TestPrintSimpleTreeWithChildren(t *testing.T) {
	logger := log.Get()
	logger.LogLevel = log.Debug
	dpService := Get(MockObject1{})
	dpService.SetLogger(logger)

	_, _ = dpService.AddRootItem("item_2", "item 2", MockObject1{id: "item_2", someStoredValue: "item 2"})
	_, _ = dpService.AddItem("item_2_child_2", "item 2 Child 2", "item_2", MockObject1{id: "item_2_child_2", someStoredValue: "item 2 Child 2"})
	_, _ = dpService.AddItem("item_2_child_1", "item 2 Child 1", "item_2", MockObject1{id: "item_2_child_1", someStoredValue: "item 2 Child 1"})
	_, _ = dpService.AddItem("item_2_child_3", "item 2 Child 3", "item_2", MockObject1{id: "item_2_child_3", someStoredValue: "item 2 Child 3"})
	_, _ = dpService.AddItem("item_2_child_1_child_2", "item 2 Child 1 Child 2", "item_2_child_1", MockObject1{id: "item_2_child_1_child_2", someStoredValue: "item 2 Child 1 Child 2"})
	_, _ = dpService.AddItem("item_2_child_1_child_1", "item 2 Child 1 Child 1", "item_2_child_1", MockObject1{id: "item_2_child_1_child_1", someStoredValue: "item 2 Child 1 Child 1"})
	_, _ = dpService.AddItem("item_2_child_3_child_2", "item 2 Child 3 Child 2", "item_2_child_3", MockObject1{id: "item_2_child_3_child_2", someStoredValue: "item 2 Child 3 Child 2"})
	_, _ = dpService.AddItem("item_2_child_3_child_1", "item 2 Child 3 Child 1", "item_2_child_3", MockObject1{id: "item_2_child_3_child_1", someStoredValue: "item 2 Child 3 Child 1"})

	_, err := dpService.Build()
	assert.NoError(t, err)

	lines := dpService.printTree(dpService.tree, 0, "", true)
	assert.Equal(t, expectedSimpleTreeWithChildrenString, strings.Join(lines, "\n"))
}

func TestPrintComplexTree(t *testing.T) {
	logger := log.Get()
	logger.LogLevel = log.Debug
	dpService := Get(MockObject1{})
	dpService.SetLogger(logger)

	_, _ = dpService.AddRootItem("item_6", "item 6", MockObject1{id: "item_6", someStoredValue: "item 6"})
	_, _ = dpService.AddRootItem("item_4", "item 4", MockObject1{id: "item_4", someStoredValue: "item 4"})
	_, _ = dpService.AddItem("item_6_child_2", "item 6 Child 2", "item_6", MockObject1{id: "item_6_child_2", someStoredValue: "item 6 Child 2"})
	_, _ = dpService.AddItem("item_6_child_1", "item 6 Child 1", "item_6", MockObject1{id: "item_6_child_1", someStoredValue: "item 6 Child 1"})
	_, _ = dpService.AddRootItem("item_1", "item 1", MockObject1{id: "item_1", someStoredValue: "item 1"})
	_, _ = dpService.AddRootItem("item_3", "item 3", MockObject1{id: "item_3", someStoredValue: "item 3"})
	_, _ = dpService.AddRootItem("item_2", "item 2", MockObject1{id: "item_2", someStoredValue: "item 2"})
	_, _ = dpService.AddRootItem("item_5", "item 5", MockObject1{id: "item_5", someStoredValue: "item 5"})
	_, _ = dpService.AddItem("item_2_child_2", "item 2 Child 2", "item_2", MockObject1{id: "item_2_child_2", someStoredValue: "item 2 Child 2"})
	_, _ = dpService.AddItem("item_2_child_1", "item 2 Child 1", "item_2", MockObject1{id: "item_2_child_1", someStoredValue: "item 2 Child 1"})
	_, _ = dpService.AddItem("item_2_child_3", "item 2 Child 3", "item_2", MockObject1{id: "item_2_child_3", someStoredValue: "item 2 Child 3"})
	_, _ = dpService.AddItem("item_2_child_1_child_2", "item 2 Child 1 Child 2", "item_2_child_1", MockObject1{id: "item_2_child_1_child_2", someStoredValue: "item 2 Child 1 Child 2"})
	_, _ = dpService.AddItem("item_2_child_1_child_1", "item 2 Child 1 Child 1", "item_2_child_1", MockObject1{id: "item_2_child_1_child_1", someStoredValue: "item 2 Child 1 Child 1"})
	_, _ = dpService.AddItem("item_2_child_3_child_2", "item 2 Child 3 Child 2", "item_2_child_3", MockObject1{id: "item_2_child_3_child_2", someStoredValue: "item 2 Child 3 Child 2"})
	_, _ = dpService.AddItem("item_2_child_3_child_1", "item 2 Child 3 Child 1", "item_2_child_3", MockObject1{id: "item_2_child_3_child_1", someStoredValue: "item 2 Child 3 Child 1"})
	_, _ = dpService.AddItem("item_4_child_1", "item 4 Child 1", "item_4", MockObject1{id: "item_4_child_1", someStoredValue: "item 4 Child 1"})

	_ = dpService.DependsOn("item_6", "item_5")
	_ = dpService.DependsOn("item_5", "item_4")
	_ = dpService.DependsOn("item_4", "item_3")
	_ = dpService.DependsOn("item_3", "item_2")
	_ = dpService.DependsOn("item_2", "item_1")

	_, err := dpService.Build()
	assert.NoError(t, err)

	lines := dpService.printTree(dpService.tree, 0, "", false)
	assert.Equal(t, expectedComplexString, strings.Join(lines, "\n"))
}

func TestShiftToWithInvalidIndex(t *testing.T) {
	logger := log.Get()
	logger.LogLevel = log.Debug
	globalDependencyTreeService = nil
	dpService := Get(MockObject1{})
	dpService.SetLogger(logger)

	_, _ = dpService.AddRootItem("item_6", "item 6", MockObject1{id: "item_6", someStoredValue: "item 6"})

	_, err := dpService.shiftTo(0, 4)

	assert.Error(t, err)
}
