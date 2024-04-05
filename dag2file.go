package merkledag

import (
	"encoding/json"
	"strings"
)

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	// 根据hash和path， 返回对应的文件, hash对应的类型是tree
	// 从KVStore中获取对象
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

	// 如果对象是一个树（目录）
	if len(obj.Links) > 0 {
		// 遍历树中的链接
		for _, link := range obj.Links {
			// 如果链接名称与路径的第一部分匹配
			if link.Name == parts[0] {
				// 如果这是路径的最后一部分，返回链接数据
				if len(parts) == 1 {
					return link.Hash
				}
				// 否则，递归调用Hash2File，使用找到的链接的hash和路径的其余部分
				return Hash2File(store, link.Hash, strings.Join(parts[1:], "/"), hp)
			}
		}
		// 如果没有找到链接，返回错误
		return nil
	}

	// 如果对象不是树，那么它应该是一个文件，所以返回其数据
	return obj.Data[0]
}
