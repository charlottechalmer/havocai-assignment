package parser

type Stack struct {
	items []map[string]interface{}
}
// add item to stack
func (s *Stack) Push(item map[string]interface{}) {
	s.items = append(s.items, item)
}
// removes and returns top item
func (s *Stack) Pop() map[string]interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}
// returns top items without removing it
func (s *Stack) Peek() map[string]interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}
