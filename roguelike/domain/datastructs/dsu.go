package datastructs

func MakeSets(parent []int, rank []int) {
	for i := 0; i < len(parent); i++ {
		parent[i] = i
		rank[i] = 0
	}
}

func FindSet(v int, parent []int) int {
	if v == parent[v] {
		return v
	}
	parent[v] = FindSet(parent[v], parent)
	return parent[v]
}

func UnionSets(v, u int, parent []int, rank []int) {
	v = FindSet(v, parent)
	u = FindSet(u, parent)

	if v != u {
		if rank[u] >= rank[v] {
			parent[v] = u
		} else {
			parent[u] = v
		}

		if rank[u] == rank[v] {
			rank[u]++
		}
	}
}
