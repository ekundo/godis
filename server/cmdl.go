package server

type listCmd struct {
	baseCmd
}

type listReadFunc func(listItem) (interface{}, error)

func (cmd *listCmd) listRead(shard *shard, key string, itemFunc listReadFunc) (interface{}, error) {
	item, err := shard.item(key, typeList)
	if err != nil {
		return "", err
	}
	return itemFunc(item.(listItem))
}

type listWriteFunc func(listItem) (interface{}, error)

func (cmd *listCmd) listWrite(shard *shard, key string, itemFunc listWriteFunc) (interface{}, error) {
	item, err := shard.item(key, typeList)
	if err != nil {
		switch err.(type) {
		case keyNotFoundError:
			item = newListItem()
		default:
			return nil, err
		}
	}
	listItem := item.(listItem)
	val, err := itemFunc(listItem)
	if err != nil {
		return nil, err
	}
	if listItem.size() < 1 {
		delete(shard.items, key)
	} else {
		shard.items[key] = item
	}
	return val, nil
}
