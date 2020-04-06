package consistent
/*
一致性哈希
*/

import (
	"hash/crc32"
	"sort"
	"strconv"
)
/*
	哈希函数类型
*/
type Hash func(data []byte) uint32

/*
	主要结构，Map是一致性哈希的抽象类
*/
type Map struct {
	hash Hash
	/*
		虚拟节点的倍数
	*/
	virMpl int
	/*
		哈希环
	*/
	keys []int
	hashMap map[int]string
}

/*
	构造函数
*/
func New(virMpl int,hash Hash) *Map{
	m:=&Map{
		virMpl:virMpl,
		hash:hash,
		hashMap:make(map[int]string),
	}
	/*
		如果hash函数为空，则使用crc32
	*/
	if m.hash==nil{
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

/*
	添加真实节点
*/
func (m* Map) Add(keys ...string){
	for _,realNodeKey:=range keys{
		for i:=0;i<m.virMpl;i++{
			/*
				keys中的每个真实节点都对映着virMpl个虚拟节点
				每个虚拟节点的key（即virNodeKey）为 i+realNodekey
				（即一个“不定数”，这里用i值，加上真实节点key
			*/
			virNodeKey := []byte(strconv.Itoa(i)+realNodeKey)
			/*
				对虚拟节点做哈希
			*/
			virNodeHash:= int(m.hash(virNodeKey))
			/*
				添加进哈希环，所以虚拟节点也存在于哈希环中
			*/
			m.keys = append(m.keys,virNodeHash)
			/*
				虚拟节点的hash对映某个真实节点的key
			*/
			m.hashMap[virNodeHash] = realNodeKey
		}
	}
	sort.Ints(m.keys)
}

/*
	访问
*/
func (m *Map) Get(key string) string{
	if len(m.keys)==0 {
		return ""
	}
	nodeHash := int(m.hash([]byte(key)))
	/*
		这里用的是内置的二分查找
	*/
	indexOfNodeHash := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i]>=nodeHash
	})
	/*
		取真实节点
	*/
	return m.hashMap[indexOfNodeHash%len(m.keys)]
}













