package server

type dictItem interface {
	item
	put(string, string) bool
	get(string) (string, error)
	keys() []string
	entries() map[string]string
	remove(string) error
	size() int
}

type baseDictItem struct {
	baseItem
	dict map[string]string
}

func newDictItem() dictItem {
	return &baseDictItem{baseItem: baseItem{}, dict: make(map[string]string)}
}

func (dictItem *baseDictItem) size() int {
	return len(dictItem.dict)
}

func (dictItem *baseDictItem) put(key string, value string) bool {
	_, ok := dictItem.dict[key]
	dictItem.dict[key] = value
	return ok
}

func (dictItem *baseDictItem) get(key string) (string, error) {
	v, found := dictItem.dict[key]
	if !found {
		return "", fieldNotFoundError{key}
	}
	return v, nil
}

func (dictItem *baseDictItem) remove(key string) error {
	_, found := dictItem.dict[key]
	if !found {
		return fieldNotFoundError{key}
	}
	delete(dictItem.dict, key)
	return nil
}

func (dictItem *baseDictItem) keys() []string {
	res := make([]string, 0, len(dictItem.dict))
	for k := range dictItem.dict {
		res = append(res, k)
	}
	return res
}

func (dictItem *baseDictItem) entries() map[string]string {
	res := make(map[string]string, len(dictItem.dict))
	for k, v := range dictItem.dict {
		res[k] = v
	}
	return res
}

func (dictItem *baseDictItem) itemType() itemType {
	return typeDict
}

var _ item = (*baseDictItem)(nil)
var _ dictItem = (*baseDictItem)(nil)
