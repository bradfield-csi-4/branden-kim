package utils

import (
	"fmt"
	"math/rand"
)

// To save some space, we make it so that we have a initial MAX level defined and when the maximum
// level needs to increase or decrease, we have to copy the pointers of every node and increase
// each node's level
const MAX_LEVEL int = 32

type SkipNode struct {
	Key        string
	Value      string
	Level_List [MAX_LEVEL]*SkipNode
}

func InitializeSkipNode(key, value string) *SkipNode {
	skipnode := new(SkipNode)
	skipnode.Key = key
	skipnode.Value = value

	return skipnode
}

type SkipList struct {
	Header *SkipNode
	// end       *SkipNode // Add this if you want to make a doublely linked SkipList
	Level     int
	Max_Level int
}

func InitializeSkipList() *SkipList {
	skiplist := new(SkipList)
	skiplist.Level = 0
	skiplist.Max_Level = MAX_LEVEL

	// Initializing the header in the skiplist
	skiplist.Header = InitializeSkipNode("", "")

	return skiplist
}

func (s *SkipList) SearchNode(key string) (*SkipNode, error) {
	var node *SkipNode = s.Header
	var current_level int

	// Loop downwards from the highest level
	for current_level = s.Max_Level - 1; current_level >= 0; current_level-- {
		// Loop through each node in the level such that we are at the node right before
		// the node we are trying to find
		for node.Level_List[current_level] != nil && node.Level_List[current_level].Key < key {
			node = node.Level_List[current_level]
		}
	}

	// At this point we should be at the node right before where the key "should be"
	node = node.Level_List[0]

	// If the node key is the same as the one we are trying to search, then return
	if node != nil && node.Key != "" && node.Key == key {
		return node, nil
	} else {
		return nil, fmt.Errorf("Could not find the node in the skip list with key: %s", key)
	}
}

func (s *SkipList) AddNode(key string, value string) error {
	var update_vector [MAX_LEVEL]*SkipNode
	var node *SkipNode = s.Header
	var current_level int

	for current_level = s.Max_Level - 1; current_level >= 0; current_level-- {
		for node.Level_List[current_level] != nil && node.Level_List[current_level].Key < key {
			node = node.Level_List[current_level]
		}
		update_vector[current_level] = node
	}

	node = node.Level_List[0]

	if node != nil && node.Key != "" && node.Key == key {
		node.Value = value
	} else {
		// If we didn't find the node with the key, we must insert it into the SkipList
		// and update the appropriate levels
		new_level := s.generateRandomLevel()
		if new_level > s.Level {
			for i := s.Level + 1; i <= new_level; i++ {
				update_vector[i] = s.Header
			}
			s.Level = new_level
		}
		new_node := InitializeSkipNode(key, value)

		// Basically set the pointers of the new node to equal the pointers immediately previous of each level
		for i := 0; i <= new_level; i++ {
			new_node.Level_List[i] = update_vector[i].Level_List[i]
			update_vector[i].Level_List[i] = new_node
		}
	}

	return nil
}

func (s *SkipList) DeleteNode(key string) error {
	var update_vector [MAX_LEVEL]*SkipNode
	var node *SkipNode = s.Header
	var current_level int

	for current_level = s.Max_Level - 1; current_level >= 0; current_level-- {
		for node.Level_List[current_level] != nil && node.Level_List[current_level].Key < key {
			node = node.Level_List[current_level]
		}
		update_vector[current_level] = node
	}

	node = node.Level_List[0]

	if node != nil && node.Key != "" && node.Key == key {
		// Basically set the pointers of the new node to equal the pointers immediately previous of each level
		for i := 0; i <= s.Level; i++ {
			if update_vector[i].Level_List[i] != node {
				break
			} else {
				update_vector[i].Level_List[i] = node.Level_List[i]
			}
		}

		// In go, we can free the memory by removing all pointer references to the node
		// and the garbage collector should handle it for us
		for s.Level > 0 && s.Header.Level_List[s.Level] == nil {
			s.Level = s.Level - 1
		}

		return nil
	} else {
		return fmt.Errorf("Delete failed since the key: %s was not found!", key)
	}

}

func (s *SkipList) generateRandomLevel() int {
	generated_level := 1
	p := 0.5

	for rand.Float64() < p && generated_level < s.Max_Level-1 {
		generated_level++
	}

	return generated_level
}
