package merkledag

import (
	"encoding/json"
	"strings"
)

// Hash2File 函数根据给定的哈希和路径从 KVStore 中获取文件内容
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 从 KVStore 中获取对象
	value, err := store.Get(hash)
	if err != nil {
		return nil
	}

	// 解析对象
	var obj Object
	err = json.Unmarshal(value, &obj)
	if err != nil {
		return nil
	}

	// 将路径分割成部分
	parts := strings.Split(path, "/")

	// 如果对象是一个目录（或树）
	if len(obj.Links) > 0 {
		// 遍历目录中的链接
		for i, link := range obj.Links {
			// 如果链接的名称与路径的第一部分匹配
			if link.Name == parts[0] {
				// 如果这是路径的最后一部分
				if len(parts) == 1 {
					// 如果链接的类型是 list，获取文件的完整内容
					if string(obj.Data[i][0]) == "list" {
						return retrieveList(store, link.Hash)
					}
					// 否则，直接返回链接的哈希值
					return link.Hash
				}
				// 如果不是路径的最后一部分，递归调用 Hash2File 函数
				return Hash2File(store, link.Hash, strings.Join(parts[1:], "/"), hp)
			}
		}
		// 如果没有找到匹配的链接，返回 nil
		return nil
	}

	// 如果对象不是目录，那么它应该是一个文件，所以返回文件的数据
	return obj.Data[0]
}

// retrieveList 函数遍历链表并获取每个块的内容
func retrieveList(store KVStore, hash []byte) []byte {
	var data []byte
	var nextHash []byte = hash

	// 遍历链表
	for nextHash != nil {
		// 从 KVStore 中获取节点
		value, err := store.Get(nextHash)
		if err != nil {
			return nil
		}

		// 解析节点
		var node ListNode
		err = json.Unmarshal(value, &node)
		if err != nil {
			return nil
		}

		// 将节点的哈希值添加到数据中
		data = append(data, node.Hash...)
		// 更新下一个哈希值
		nextHash = node.Next
	}

	// 返回数据
	return data
}
