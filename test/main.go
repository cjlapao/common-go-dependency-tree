package main

import (
	"fmt"

	dependency_tree "github.com/cjlapao/common-go-dependency-tree/dependencytree"
)

type Item struct {
	ID   string
	Name string
}

func main() {
	dependencyService := dependency_tree.Get[Item](Item{})
	_, _ = dependencyService.AddRootItem("ITEM_6", "Item 6", Item{ID: "ITEM_6", Name: "Item 6"})
	_, _ = dependencyService.AddRootItem("ITEM_1", "Item 1", Item{ID: "ITEM_1", Name: "Item 1"})
	_, _ = dependencyService.AddRootItem("ITEM_4", "Item 4", Item{ID: "ITEM_4", Name: "Item 4"})
	_, _ = dependencyService.AddRootItem("ITEM_2", "Item 2", Item{ID: "ITEM_2", Name: "Item 2"})
	_, _ = dependencyService.AddRootItem("ITEM_5", "Item 5", Item{ID: "ITEM_5", Name: "Item 5"})
	_, _ = dependencyService.AddRootItem("ITEM_3", "Item 3", Item{ID: "ITEM_3", Name: "Item 3"})
	_, _ = dependencyService.AddItem("ITEM_2_CHILD_1", "Item 2 Child 1", "ITEM_2", Item{ID: "ITEM_2_CHILD_1", Name: "Item 2 Child 1"})
	_, _ = dependencyService.AddItem("ITEM_2_CHILD_2", "Item 2 Child 2", "ITEM_2", Item{ID: "ITEM_2_CHILD_2", Name: "Item 2 Child 2"})
	_, _ = dependencyService.AddItem("ITEM_2_CHILD_3", "Item 2 Child 3", "ITEM_2", Item{ID: "ITEM_2_CHILD_3", Name: "Item 2 Child 3"})

	_, _ = dependencyService.AddItem("ITEM_2_CHILD_1_CHILD_1", "Item 2 Child 1 Child 1", "ITEM_2_CHILD_1", Item{ID: "ITEM_2_CHILD_1_CHILD_1", Name: "Item 2 Child 1 Child 1"})
	_, _ = dependencyService.AddItem("ITEM_2_CHILD_1_CHILD_2", "Item 2 Child 1 Child 2", "ITEM_2_CHILD_1", Item{ID: "ITEM_2_CHILD_1_CHILD_2", Name: "Item 2 Child 1 Child 2"})

	_, _ = dependencyService.AddItem("ITEM_2_CHILD_3_CHILD_1", "Item 2 Child 3 Child 1", "ITEM_2_CHILD_3", Item{ID: "ITEM_2_CHILD_3_CHILD_1", Name: "Item 2 Child 3 Child 1"})
	_, _ = dependencyService.AddItem("ITEM_2_CHILD_3_CHILD_2", "Item 2 Child 3 Child 2", "ITEM_2_CHILD_3", Item{ID: "ITEM_2_CHILD_3_CHILD_2", Name: "Item 2 Child 3 Child 2"})

	_, _ = dependencyService.AddItem("ITEM_4_CHILD_1", "Item 4 Child 1", "ITEM_4", Item{ID: "ITEM_4_CHILD_1", Name: "Item 4 Child 1"})
	_, _ = dependencyService.AddItem("ITEM_6_CHILD_1", "Item 6 Child 1", "ITEM_6", Item{ID: "ITEM_6_CHILD_1", Name: "Item 6 Child 1"})
	_, _ = dependencyService.AddItem("ITEM_6_CHILD_2", "Item 6 Child 2", "ITEM_6", Item{ID: "ITEM_6_CHILD_2", Name: "Item 6 Child 2"})

	_, _ = dependencyService.AddItem("ITEM_2_CHILD_3_CHILD_1_CHILD_1", "Item 2 Child 3 Child 1 Child 1", "ITEM_2_CHILD_3_CHILD_1", Item{ID: "ITEM_2_CHILD_3_CHILD_1_CHILD_1", Name: "Item 2 Child 3 Child 1 Child 1"})
	_, _ = dependencyService.AddItem("ITEM_2_CHILD_3_CHILD_1_CHILD_2", "Item 2 Child 3 Child 1 Child 2", "ITEM_2_CHILD_3_CHILD_1", Item{ID: "ITEM_2_CHILD_3_CHILD_1_CHILD_2", Name: "Item 2 Child 3 Child 1 Child 2"})

	_ = dependencyService.DependsOn("ITEM_6", "ITEM_5")
	_ = dependencyService.DependsOn("ITEM_5", "ITEM_4")
	_ = dependencyService.DependsOn("ITEM_4", "ITEM_3")
	_ = dependencyService.DependsOn("ITEM_3", "ITEM_2")
	_ = dependencyService.DependsOn("ITEM_2", "ITEM_1")
	fmt.Println("Before:")
	dependencyService.PrintFlatTree()
	_, _ = dependencyService.Build()
	fmt.Println("After:")

	println(dependencyService.String())
}
