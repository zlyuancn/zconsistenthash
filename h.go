/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/12/13
   Description :
-------------------------------------------------
*/

package zconsistenthash

import (
    "fmt"
    "hash/crc32"
    "sort"
)

type HashFn func(data []byte) uint32

type HashCircle struct {
    replicas int
    hashFn   HashFn
    keys     []int // asc
    hashMap  map[int]string
}

// 创建一个一致性hash环
func New(replicas int, hashFn HashFn, keys ...string) *HashCircle {
    if hashFn == nil {
        hashFn = crc32.ChecksumIEEE
    }
    m := &HashCircle{
        replicas: replicas,
        hashFn:   hashFn,
        hashMap:  make(map[int]string),
    }
    m.add(keys)
    return m
}

func (m *HashCircle) add(keys []string) {
    for _, key := range keys {
        for i := 0; i <= m.replicas; i++ {
            hash := int(m.hashFn([]byte(fmt.Sprintf("%s %d", key, i))))
            m.keys = append(m.keys, hash)
            m.hashMap[hash] = key
        }
    }
    sort.Ints(m.keys)
}

// 判断是否没有任何key
func (m *HashCircle) Empty() bool {
    return len(m.keys) == 0
}

// Get获取散列中最接近提供的键的项.
func (m *HashCircle) Get(key string) string {
    if len(m.keys) == 0 {
        return ""
    }

    hash := int(m.hashFn([]byte(key)))

    // 二叉搜索副本
    idx := sort.Search(len(m.keys), func(i int) bool {
        return m.keys[i] >= hash
    })

    // 环尾
    if idx == len(m.keys) {
        idx = 0
    }

    return m.hashMap[m.keys[idx]]
}
