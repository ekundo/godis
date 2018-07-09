package server

type dictCmd struct {
	baseCmd
}

type dictReadFunc func(dictItem) (interface{}, error)

func (cmd *dictCmd) dictRead(shard *shard, key string, itemFunc dictReadFunc) (interface{}, error) {
	item, err := shard.item(key, typeDict)
	if err != nil {
		return "", err
	}
	return itemFunc(item.(dictItem))
}

type dictWriteFunc func(dictItem) (interface{}, error)

func (cmd *dictCmd) dictWrite(shard *shard, key string, itemFunc dictWriteFunc) (interface{}, error) {
	item, err := shard.item(key, typeDict)
	if err != nil {
		switch err.(type) {
		case keyNotFoundError:
			item = newDictItem()
		default:
			return nil, err
		}
	}
	dictItem := item.(dictItem)
	val, err := itemFunc(dictItem)
	if err != nil {
		return nil, err
	}
	if dictItem.size() < 1 {
		delete(shard.items, key)
	} else {
		shard.items[key] = item
	}
	return val, nil
}
