package dependency_tree

import (
	"testing"

	log "github.com/cjlapao/common-go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockObject1 struct {
	id              string
	someStoredValue string
}

func (d *MockObject1) ID() string {
	return d.id
}

func (d *MockObject1) StoredValue() string {
	return d.someStoredValue
}

type MockObject2 struct {
	id              string
	someStoredValue string
}

func (d *MockObject2) ID() string {
	return d.id
}

func (d *MockObject2) StoredValue() string {
	return d.someStoredValue
}

func TestGet(t *testing.T) {
	// Test case 1: globalDependencyTreeServiceNew is nil
	t.Run("globalDependencyTreeService is nil", func(t *testing.T) {
		lock.Lock()
		globalDependencyTreeService = nil
		lock.Unlock()

		service := Get(MockObject1{})

		if service == nil {
			t.Errorf("Expected non-nil service, got nil")
		}
	})

	// Test case 2: globalDependencyTreeServiceNew is not nil
	t.Run("globalDependencyTreeService is not nil", func(t *testing.T) {
		lock.Lock()
		// mockClass := DependencyTreeObjectMock{
		// 	id:              "test",
		// 	someStoredValue: "test",
		// }
		mockService := DependencyTreeService[MockObject1]{}

		globalDependencyTreeService = []interface{}{&mockService}
		lock.Unlock()

		service := Get(MockObject1{})

		if service == nil {
			t.Errorf("Expected non-nil service, got nil")
		}
	})

	t.Run("globalDependencyTreeService is not nil and gets right interface", func(t *testing.T) {
		lock.Lock()
		mockClass1 := MockObject1{
			id:              "test",
			someStoredValue: "test",
		}
		mockTreeObject1 := DependencyTreeItem[MockObject1]{
			ID:            "test",
			Name:          "test",
			isDependentOn: []string{},
			parentName:    "root",
			obj:           mockClass1,
			requiredBy:    []string{},
		}

		mockClass2 := MockObject2{
			id:              "test1",
			someStoredValue: "test1",
		}
		mockTreeObject2 := DependencyTreeItem[MockObject2]{
			ID:            "test1",
			Name:          "test1",
			isDependentOn: []string{},
			parentName:    "root",
			obj:           mockClass2,
			requiredBy:    []string{},
		}

		mockService1 := DependencyTreeService[MockObject1]{}
		_ = mockService1.AddDependencyTreeItem(&mockTreeObject1)

		mockService2 := DependencyTreeService[MockObject2]{}
		_ = mockService2.AddDependencyTreeItem(&mockTreeObject2)

		globalDependencyTreeService = []interface{}{&mockService1, &mockService2}
		lock.Unlock()

		service1 := Get(MockObject1{})
		assert.Equal(t, service1.flatTree[0].obj.ID(), mockClass1.ID())
		assert.Equal(t, service1.flatTree[0].obj.StoredValue(), mockClass1.StoredValue())

		service2 := Get(MockObject2{})
		assert.Equal(t, service2.flatTree[0].obj.ID(), mockClass2.ID())
		assert.Equal(t, service2.flatTree[0].obj.StoredValue(), mockClass2.StoredValue())
	})
}

func TestSetDebug(t *testing.T) {
	t.Run("Set debug modefor DependencyTreeService[DependencyTreeObjectMock]", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		service.SetDebug(true)

		assert.True(t, service.IsDebug())
	})

	t.Run("Set debug modefor DependencyTreeService[DependencyTreeObjectMock]", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		service.SetDebug(true)
		assert.True(t, service.IsDebug())

		service.SetDebug(false)
		assert.False(t, service.IsDebug())
	})
}

func TestSetVerbose(t *testing.T) {
	t.Run("Set verbose modefor DependencyTreeService[DependencyTreeObjectMock]", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		service.SetVerbose(true)

		assert.True(t, service.IsVerbose())
	})

	t.Run("Set verbose to false modefor DependencyTreeService[DependencyTreeObjectMock]", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		service.SetVerbose(true)
		assert.True(t, service.IsVerbose())

		service.SetVerbose(false)
		assert.False(t, service.IsVerbose())
	})
}

func TestSetLogger(t *testing.T) {
	mockLogger := &log.LoggerService{} // Create a mock logger

	// Test case 1: Set logger for DependencyTreeService[DependencyTreeObjectMock]
	t.Run("Set logger for DependencyTreeService[DependencyTreeObjectMock]", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{} // Create a DependencyTreeService instance
		service.SetLogger(mockLogger)                    // Set the mock logger

		// Assert that the logger is set correctly
		if service.logger != mockLogger {
			t.Errorf("Expected logger to be set to mockLogger, got %v", service.logger)
		}
	})

	// Test case 2: Set logger for DependencyTreeService[ObjectMock1]
	t.Run("Set logger for DependencyTreeService[ObjectMock1]", func(t *testing.T) {
		service := &DependencyTreeService[MockObject2]{} // Create a DependencyTreeService instance
		service.SetLogger(mockLogger)                    // Set the mock logger

		// Assert that the logger is set correctly
		if service.logger != mockLogger {
			t.Errorf("Expected logger to be set to mockLogger, got %v", service.logger)
		}
	})
}

func TestAddRootItem(t *testing.T) {
	mockClass := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass,
		requiredBy:    []string{},
	}

	t.Run("Add root item successfully", func(t *testing.T) {
		service := &DependencyTreeService[*MockObject1]{}
		_, err := service.AddRootItem("test", "test", &mockClass)

		require.Nil(t, err)
		assert.Len(t, service.flatTree, 1)
		assert.Equal(t, service.flatTree[0].obj.ID(), mockClass.ID())
		assert.Equal(t, "root", service.flatTree[0].Parent())
	})

	t.Run("Fail to create root item with missing id", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}

		_, err := service.AddRootItem("", "test", mockClass)

		assert.NotNil(t, err)
		assert.Len(t, service.flatTree, 0)
	})

	t.Run("Fail to create root item with missing name", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}

		_, err := service.AddRootItem("test", "", mockClass)

		assert.NotNil(t, err)
		assert.Len(t, service.flatTree, 0)
	})

	t.Run("Fail to add root item", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject)

		_, err := service.AddRootItem("test", "test", mockClass)

		assert.NotNil(t, err)
		assert.Len(t, service.flatTree, 1)
	})
}

func TestAddItem(t *testing.T) {
	mockParentClass := MockObject1{
		id:              "testParent",
		someStoredValue: "testParent",
	}

	mockChildClass := MockObject1{
		id:              "testChild",
		someStoredValue: "testChild",
	}
	mockChildTreeObject := DependencyTreeItem[MockObject1]{
		ID:            "testChild",
		Name:          "testChild",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockChildClass,
		requiredBy:    []string{},
	}

	t.Run("Add dependency successfully", func(t *testing.T) {
		service := &DependencyTreeService[*MockObject1]{}
		_, errParent := service.AddRootItem("TestParent", "TestParent", &mockParentClass)
		_, errChild := service.AddItem("testChild", "testChild", "TestParent", &mockChildClass)

		require.Nil(t, errParent)
		require.Nil(t, errChild)
		assert.Len(t, service.flatTree, 2)
		assert.Equal(t, service.flatTree[0].obj.ID(), mockParentClass.ID())
		assert.Equal(t, "root", service.flatTree[0].Parent())
		assert.Equal(t, service.flatTree[1].obj.ID(), mockChildClass.ID())
		assert.Equal(t, "TestParent", service.flatTree[1].Parent())
	})

	t.Run("Fail to create DependencyTreeItem with missing parent", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_, err := service.AddItem("", "test", "test", mockChildClass)

		require.NotNil(t, err)
		assert.Len(t, service.flatTree, 0)
	})

	t.Run("Fail to create DependencyTreeItem with missing id", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_, errParent := service.AddRootItem("TestParent", "TestParent", mockParentClass)
		require.Nil(t, errParent)

		_, err := service.AddItem("", "test", "TestParent", mockParentClass)

		assert.NotNil(t, err)
		assert.Len(t, service.flatTree, 1)
	})

	t.Run("Fail to create DependencyTreeItem with missing name", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_, errParent := service.AddRootItem("TestParent", "TestParent", mockParentClass)
		require.Nil(t, errParent)

		_, err := service.AddItem("test", "", "TestParent", mockParentClass)

		assert.NotNil(t, err)
		assert.Len(t, service.flatTree, 1)
	})

	t.Run("Fail to add DependencyTreeItem", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_, errParent := service.AddRootItem("TestParent", "TestParent", mockParentClass)
		require.Nil(t, errParent)

		_ = service.AddDependencyTreeItem(&mockChildTreeObject)

		_, err := service.AddItem("testChild", "testChild", "TestParent", mockParentClass)

		assert.NotNil(t, err)
		assert.Len(t, service.flatTree, 2)
	})
}

func TestAddDependencyTreeItem(t *testing.T) {
	mockClass1 := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject1 := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass1,
		requiredBy:    []string{},
	}

	t.Run("Add DependencyTreeItem[DependencyTreeObjectMock] to DependencyTreeService[DependencyTreeObjectMock]", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		err := service.AddDependencyTreeItem(&mockTreeObject1)

		assert.Nil(t, err)
		assert.Equal(t, service.flatTree[0].obj.ID(), mockClass1.ID())
	})

	t.Run("Add DependencyTreeItem[ObjectMock] to DependencyTreeService[ObjectMock] twice", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		err := service.AddDependencyTreeItem(&mockTreeObject1)

		assert.Nil(t, err)
		assert.Equal(t, service.flatTree[0].obj.ID(), mockClass1.ID())

		err = service.AddDependencyTreeItem(&mockTreeObject1)
		assert.NotNil(t, err)
	})
}

func TestRemove(t *testing.T) {
	mockClass1 := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject1 := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass1,
		requiredBy:    []string{},
	}

	mockClass2 := MockObject1{
		id:              "test2",
		someStoredValue: "test2",
	}
	mockTreeObject2 := DependencyTreeItem[MockObject1]{
		ID:            "test2",
		Name:          "test2",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass2,
		requiredBy:    []string{},
	}

	t.Run("Remove existing item", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		err := service.RemoveDependencyTreeItem(&mockTreeObject1)

		assert.Nil(t, err)
		assert.Len(t, service.flatTree, 1)
		assert.Equal(t, service.flatTree[0].obj.ID(), mockClass2.ID())
	})

	t.Run("Remove non-existing item", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)

		nonExistingItem := &DependencyTreeItem[MockObject1]{
			ID: "non-existing",
		}

		err := service.RemoveDependencyTreeItem(nonExistingItem)

		assert.NotNil(t, err)
		assert.Len(t, service.flatTree, 1)
		assert.Equal(t, service.flatTree[0].obj.ID(), mockClass1.ID())
	})
}

func TestGetItem(t *testing.T) {
	mockClass1 := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject1 := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass1,
		requiredBy:    []string{},
	}

	mockClass2 := MockObject1{
		id:              "test2",
		someStoredValue: "test2",
	}
	mockTreeObject2 := DependencyTreeItem[MockObject1]{
		ID:            "test2",
		Name:          "test2",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass2,
		requiredBy:    []string{},
	}

	t.Run("Get existing item by ID", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		item := service.GetItem("test")

		assert.NotNil(t, item)
		assert.Equal(t, item.ID, mockClass1.ID())
	})

	t.Run("Get non-existing item by ID", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		item := service.GetItem("non-existing")

		assert.Nil(t, item)
	})

	t.Run("Get existing item by name (case-insensitive)", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		item := service.GetItem("Test2")

		assert.NotNil(t, item)
		assert.Equal(t, item.ID, mockClass2.ID())
	})

	t.Run("Get non-existing item by name (case-insensitive)", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		item := service.GetItem("Non-Existing")

		assert.Nil(t, item)
	})
}

func TestGetParentItems(t *testing.T) {
	mockClass1 := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject1 := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass1,
		requiredBy:    []string{},
	}

	mockClass2 := MockObject1{
		id:              "test2",
		someStoredValue: "test2",
	}
	mockTreeObject2 := DependencyTreeItem[MockObject1]{
		ID:            "test2",
		Name:          "test2",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass2,
		requiredBy:    []string{},
	}

	t.Run("Get parent items when parent exists", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		parent := "root"
		expectedItems := []*DependencyTreeItem[MockObject1]{&mockTreeObject1, &mockTreeObject2}

		items := service.GetItemByParent(parent)

		assert.Equal(t, expectedItems, items)
	})

	t.Run("Get parent items when parent does not exist", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		parent := "non-existing"

		items := service.GetItemByParent(parent)

		assert.Empty(t, items)
	})
}

func TestGetItemChildren(t *testing.T) {
	mockClass1 := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject1 := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass1,
		requiredBy:    []string{},
	}

	mockClass2 := MockObject1{
		id:              "test2",
		someStoredValue: "test2",
	}
	mockTreeObject2 := DependencyTreeItem[MockObject1]{
		ID:            "test2",
		Name:          "test2",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass2,
		requiredBy:    []string{},
		Children: []*DependencyTreeItem[MockObject1]{
			&mockTreeObject1,
		},
	}

	t.Run("Get children by ID", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		children := service.GetItemChildren("test2")

		assert.NotNil(t, children)
		assert.Len(t, children, 1) // Replace 0 with the expected number of children
	})

	t.Run("Get children by name (case-insensitive)", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		children := service.GetItemChildren("Test2")

		assert.NotNil(t, children)
		assert.Len(t, children, 1) // Replace 0 with the expected number of children
	})

	t.Run("Get children with non-existing name or ID", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		children := service.GetItemChildren("non-existing")

		assert.NotNil(t, children)
		assert.Len(t, children, 0)
	})
}

func TestString(t *testing.T) {
	logger := log.Get()
	logger.LogLevel = log.Debug
	dpService := Get(MockObject1{})
	dpService.SetLogger(logger)
	dpService.Clear()

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

	result := dpService.String()
	assert.Equal(t, expectedComplexString, result)
}

func TestGetItemDependencies(t *testing.T) {
	mockClass1 := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject1 := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass1,
		requiredBy:    []string{},
	}

	mockClass2 := MockObject1{
		id:              "test2",
		someStoredValue: "test2",
	}
	mockTreeObject2 := DependencyTreeItem[MockObject1]{
		ID:            "test2",
		Name:          "test2",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass2,
		requiredBy:    []string{},
	}

	t.Run("Get item dependencies by ID", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		dependencies := service.GetItemDependencies("test")
		assert.Len(t, dependencies, 0)

		dependencies = service.GetItemDependencies("test2")
		assert.Len(t, dependencies, 0)
	})

	t.Run("Get item dependencies by name", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		dependencies := service.GetItemDependencies("test")
		assert.Len(t, dependencies, 0)

		dependencies = service.GetItemDependencies("test2")
		assert.Len(t, dependencies, 0)
	})

	t.Run("Get item dependencies with case-insensitive search", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		dependencies := service.GetItemDependencies("TEST")
		assert.Len(t, dependencies, 0)

		dependencies = service.GetItemDependencies("TEST2")
		assert.Len(t, dependencies, 0)
	})

	t.Run("Get item dependencies for non-existing item", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject1)
		_ = service.AddDependencyTreeItem(&mockTreeObject2)

		dependencies := service.GetItemDependencies("non-existing")
		assert.Len(t, dependencies, 0)
	})
}

func TestTree(t *testing.T) {
	mockClass := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass,
		requiredBy:    []string{},
	}

	t.Run("Return the correct tree", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject)

		tree := service.Tree()

		assert.Len(t, tree, 1)
		assert.Equal(t, tree[0].ID, mockTreeObject.ID)
		assert.Equal(t, tree[0].Name, mockTreeObject.Name)
		assert.Equal(t, tree[0].isDependentOn, mockTreeObject.isDependentOn)
		assert.Equal(t, tree[0].parentName, mockTreeObject.parentName)
		assert.Equal(t, tree[0].obj, mockTreeObject.obj)
		assert.Equal(t, tree[0].requiredBy, mockTreeObject.requiredBy)
	})
}

func TestFlatTree(t *testing.T) {
	mockClass := MockObject1{
		id:              "test",
		someStoredValue: "test",
	}
	mockTreeObject := DependencyTreeItem[MockObject1]{
		ID:            "test",
		Name:          "test",
		isDependentOn: []string{},
		parentName:    "root",
		obj:           mockClass,
		requiredBy:    []string{},
	}

	t.Run("Return flat tree successfully", func(t *testing.T) {
		service := &DependencyTreeService[MockObject1]{}
		_ = service.AddDependencyTreeItem(&mockTreeObject)

		flatTree := service.FlatTree()

		assert.Len(t, flatTree, 1)
		assert.Equal(t, flatTree[0].obj.ID(), mockClass.ID())
	})
}
