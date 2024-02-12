package dependency_tree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependsOn(t *testing.T) {
	dt := &DependencyTreeItem[MockObject1]{
		isDependentOn: []string{"item1", "item2"},
	}

	t.Run("Depends on existing item", func(t *testing.T) {
		err := dt.DependsOn("item1")
		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("item item1 already exists"), err)
	})

	t.Run("Depends on new item", func(t *testing.T) {
		err := dt.DependsOn("item3")
		assert.NoError(t, err)
		assert.Equal(t, []string{"item1", "item2", "item3"}, dt.isDependentOn)
	})
}

func TestGetProperty(t *testing.T) {
	dt := &DependencyTreeItem[MockObject1]{
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	t.Run("Existing key with value", func(t *testing.T) {
		value := dt.GetProperty("key1", "default")
		assert.Equal(t, "value1", value)
	})

	t.Run("Existing key with different type value", func(t *testing.T) {
		value := dt.GetProperty("key2", "default")
		assert.Equal(t, 123, value)
	})

	t.Run("Non-existing key with default value", func(t *testing.T) {
		value := dt.GetProperty("key3", "default")
		assert.Equal(t, "default", value)
	})

	t.Run("Non-existing key without default value", func(t *testing.T) {
		value := dt.GetProperty("key3", nil)
		assert.Nil(t, value)
	})

	t.Run("Empty metadata", func(t *testing.T) {
		dt.Metadata = nil
		value := dt.GetProperty("key3", nil)
		assert.Nil(t, value)
	})
}

func TestSetProperty(t *testing.T) {
	dt := &DependencyTreeItem[MockObject1]{
		Metadata: make(map[string]interface{}),
	}

	t.Run("Set property with existing key", func(t *testing.T) {
		err := dt.SetProperty("key1", "new value")
		assert.NoError(t, err)
		assert.Equal(t, "new value", dt.Metadata["key1"])
	})

	t.Run("Set property with new key", func(t *testing.T) {
		err := dt.SetProperty("key2", 456)
		assert.NoError(t, err)
		assert.Equal(t, 456, dt.Metadata["key2"])
	})

	t.Run("Set property with new ke, empty metadata", func(t *testing.T) {
		dt.Metadata = nil
		err := dt.SetProperty("key2", 456)
		assert.NoError(t, err)
		assert.Equal(t, 456, dt.Metadata["key2"])
	})
}

func TestRequiredBy(t *testing.T) {
	dt := &DependencyTreeItem[MockObject1]{
		requiredBy: []string{"item1", "item2"},
	}

	t.Run("Get required by", func(t *testing.T) {
		result := dt.RequiredBy()
		assert.Equal(t, []string{"item1", "item2"}, result)
	})
}
