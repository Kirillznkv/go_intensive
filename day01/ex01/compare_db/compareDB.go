package compare_db

import (
	"fmt"

	read "../read_db"
)

func inSlice_name(name string, sl *[]read.Cake) bool {
	for _, v := range *sl {
		if v.Name == name {
			return true
		}
	}
	return false
}

func inSlice_ingredient(name string, sl *[]read.Item) bool {
	for _, v := range *sl {
		if v.Name == name {
			return true
		}
	}
	return false
}

func compareCakeName(oldData, newData *read.Recipes) {
	for _, v := range newData.CakeList {
		if ok := inSlice_name(v.Name, &oldData.CakeList); ok == false {
			fmt.Printf("ADDED cake \"%s\"\n", v.Name)
		}
	}
	for _, v := range oldData.CakeList {
		if ok := inSlice_name(v.Name, &newData.CakeList); ok == false {
			fmt.Printf("REMOVED cake \"%s\"\n", v.Name)
		}
	}
}

func compareCookingTime(oldData, newData *read.Recipes) {
	for _, v := range oldData.CakeList {
		for _, find := range newData.CakeList {
			if v.Name == find.Name && v.Time != find.Time {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", v.Name, find.Time, v.Time)
			}
		}
	}
}

func findItem(name string, items *[]read.Item) *read.Item {
	for _, v := range *items {
		if name == v.Name {
			return &v
		}
	}
	return nil
}

func outputChangedItems(oldIt, newIt *read.Item, cakeName string) {
	if oldIt.Count != newIt.Count {
		fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldIt.Name, cakeName, newIt.Count, oldIt.Count)
	}
	if oldIt.Unit != newIt.Unit {
		if oldIt.Unit == "" {
			fmt.Printf("ADDED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", newIt.Unit, newIt.Name, cakeName)
		} else if newIt.Unit == "" {
			fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", oldIt.Unit, newIt.Name, cakeName)
		} else {
			fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldIt.Name, cakeName, newIt.Unit, oldIt.Unit)
		}
	}
}

func outputChangedIngredients(oldData, newData *[]read.Item, cakeName string) {
	for _, n := range *newData {
		if ok := inSlice_ingredient(n.Name, oldData); ok == false {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", n.Name, cakeName)
		}
	}
	for _, o := range *oldData {
		if ok := inSlice_ingredient(o.Name, newData); ok == false {
			fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", o.Name, cakeName)
		} else {
			outputChangedItems(&o, findItem(o.Name, newData), cakeName)
		}
	}
}

func compareIngredient(oldData, newData *read.Recipes) {
	for _, o := range oldData.CakeList {
		for _, n := range newData.CakeList {
			if o.Name == n.Name {
				outputChangedIngredients(&o.Ingredient, &n.Ingredient, o.Name)
				break
			}
		}
	}
}

func Compare(oldF, newF string) {
	oldData := read.ReadDB(&oldF)
	newData := read.ReadDB(&newF)
	compareCakeName(&oldData, &newData)
	compareCookingTime(&oldData, &newData)
	compareIngredient(&oldData, &newData)
}
