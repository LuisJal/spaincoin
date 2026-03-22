package block

import "github.com/spaincoin/spaincoin/core/crypto"

// MerkleTree representa un árbol de Merkle construido a partir de hashes de transacciones.
// Se usa para calcular el MerkleRoot de un bloque, lo que permite verificar
// eficientemente si una transacción pertenece al bloque.
type MerkleTree struct {
	root crypto.Hash
}

// NewMerkleTree construye un árbol de Merkle a partir de una lista de hashes de transacciones.
// Si la lista está vacía, el árbol tiene un hash raíz cero.
// Si el número de nodos es impar, el último nodo se duplica.
func NewMerkleTree(txHashes []crypto.Hash) *MerkleTree {
	if len(txHashes) == 0 {
		return &MerkleTree{root: crypto.Hash{}}
	}

	// Capa inicial: copiar los hashes de las transacciones
	layer := make([]crypto.Hash, len(txHashes))
	copy(layer, txHashes)

	// Reducir el árbol hacia arriba hasta obtener la raíz
	for len(layer) > 1 {
		layer = buildNextLayer(layer)
	}

	return &MerkleTree{root: layer[0]}
}

// buildNextLayer calcula la siguiente capa del árbol combinando pares de hashes.
// Si el número de nodos es impar, el último se duplica.
func buildNextLayer(layer []crypto.Hash) []crypto.Hash {
	// Si es impar, duplicar el último elemento
	if len(layer)%2 != 0 {
		layer = append(layer, layer[len(layer)-1])
	}

	next := make([]crypto.Hash, len(layer)/2)
	for i := 0; i < len(layer); i += 2 {
		next[i/2] = hashPair(layer[i], layer[i+1])
	}
	return next
}

// hashPair combina dos hashes de Merkle concatenándolos y aplicando SHA-256.
func hashPair(left, right crypto.Hash) crypto.Hash {
	data := make([]byte, 64)
	copy(data[:32], left[:])
	copy(data[32:], right[:])
	return crypto.HashBytes(data)
}

// RootHash devuelve el hash raíz del árbol de Merkle.
// Si el árbol está vacío, devuelve el hash cero.
func (m *MerkleTree) RootHash() crypto.Hash {
	return m.root
}
