package dependencytree

import (
	"testing"

	log "github.com/cjlapao/common-go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanTestBuild(service *DependencyTreeService[MockObject1]) {
	service.flatTree = []*DependencyTreeItem[MockObject1]{}
	service.SetDebug(false)
	service.SetVerbose(false)
}

func TestBuild(t *testing.T) {
	logger := log.Get()
	logger.LogLevel = log.Debug
	dpService := Get(MockObject1{})
	dpService.SetLogger(logger)

	t.Run("Simple tree build", func(t *testing.T) {
		_, _ = dpService.AddRootItem("item_6", "item 6", MockObject1{id: "item_6", someStoredValue: "item 6"})
		_, _ = dpService.AddRootItem("item_4", "item 4", MockObject1{id: "item_4", someStoredValue: "item 4"})
		_, _ = dpService.AddRootItem("item_1", "item 1", MockObject1{id: "item_1", someStoredValue: "item 1"})
		_, _ = dpService.AddRootItem("item_3", "item 3", MockObject1{id: "item_3", someStoredValue: "item 3"})
		_, _ = dpService.AddRootItem("item_2", "item 2", MockObject1{id: "item_2", someStoredValue: "item 2"})
		_, _ = dpService.AddRootItem("item_5", "item 5", MockObject1{id: "item_5", someStoredValue: "item 5"})

		_ = dpService.DependsOn("item_6", "item_5")
		_ = dpService.DependsOn("item_5", "item_4")
		_ = dpService.DependsOn("item_4", "item_3")
		_ = dpService.DependsOn("item_3", "item_2")
		_ = dpService.DependsOn("item_2", "item_1")

		_, err := dpService.Build()

		require.NoError(t, err)
		require.Equal(t, 6, len(dpService.flatTree))
		assert.Equal(t, "item 1", dpService.flatTree[0].Name)
		assert.Equal(t, "item 2", dpService.flatTree[1].Name)
		assert.Equal(t, "item 3", dpService.flatTree[2].Name)
		assert.Equal(t, "item 4", dpService.flatTree[3].Name)
		assert.Equal(t, "item 5", dpService.flatTree[4].Name)
		assert.Equal(t, "item 6", dpService.flatTree[5].Name)

		cleanTestBuild(dpService)
	})

	t.Run("tree build width childs", func(t *testing.T) {
		dpService.Clear()
		dpService.SetDebug(true)
		dpService.SetVerbose(true)

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

		require.NoError(t, err)
		require.Equal(t, 6, len(dpService.tree))
		assert.Equal(t, "item 1", dpService.tree[0].Name)
		assert.Equal(t, "item 2", dpService.tree[1].Name)
		assert.Equal(t, "item 3", dpService.tree[2].Name)
		assert.Equal(t, "item 4", dpService.tree[3].Name)
		assert.Equal(t, "item 5", dpService.tree[4].Name)
		assert.Equal(t, "item 6", dpService.tree[5].Name)

		cleanTestBuild(dpService)
	})
}
