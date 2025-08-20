package template

import "testing"

func TestTrie(t *testing.T) {
	tests := []struct {
		name   string
		words  []string
		prefix string
		want   bool
	}{
		{"正常", []string{"apple", "app", "applet"}, "app", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewTrie()
			for _, word := range tt.words {
				trie.Insert(word)

				if r := trie.Search(word); r != tt.want {
					t.Errorf("Search() = %v, want %v", r, tt.want)
				}
			}

			if r := trie.StartsWith(tt.prefix); r != tt.want {
				t.Errorf("StartsWith() = %v, want %v", r, tt.want)
			}
		})
	}
}

type Trie struct {
	children map[rune]*Trie
	isEnd    bool
}

func NewTrie() Trie {
	return Trie{
		children: make(map[rune]*Trie),
	}
}

func (t *Trie) Insert(word string) {
	node := t
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &Trie{
				children: make(map[rune]*Trie),
			}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (t *Trie) SearchPrefix(prefix string) *Trie {
	node := t
	for _, ch := range prefix {
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

func (t *Trie) Search(word string) bool {
	node := t.SearchPrefix(word)
	return node != nil && node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	return t.SearchPrefix(prefix) != nil
}
