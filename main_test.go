package main

import (
	"testing"
)

// TestUnionFind verifica la funcionalidad de Union-Find
func TestUnionFind(t *testing.T) {
	uf := NewUnionFind(5)

	// Verificar estado inicial
	for i := 0; i < 5; i++ {
		if uf.Find(i) != i {
			t.Errorf("Nodo %d debería ser su propio representante", i)
		}
	}

	// Unir nodos 0 y 1
	if !uf.Union(0, 1) {
		t.Error("Union(0, 1) debería retornar true")
	}

	// Verificar que están en el mismo conjunto
	if uf.Find(0) != uf.Find(1) {
		t.Error("Nodos 0 y 1 deberían estar en el mismo conjunto")
	}

	// Intentar unir de nuevo (debería fallar)
	if uf.Union(0, 1) {
		t.Error("Union(0, 1) debería retornar false la segunda vez")
	}

	// Unir más nodos
	uf.Union(1, 2)
	uf.Union(3, 4)

	// Verificar conjuntos
	if uf.Find(0) != uf.Find(2) {
		t.Error("Nodos 0 y 2 deberían estar en el mismo conjunto")
	}

	if uf.Find(0) == uf.Find(3) {
		t.Error("Nodos 0 y 3 no deberían estar en el mismo conjunto")
	}
}

// TestSmallGraph prueba Prim con un grafo pequeño conocido
func TestSmallGraph(t *testing.T) {
	// Crear grafo de ejemplo
	// Grafo:
	//   0 --10-- 1
	//   |  \     |
	//   6   5   15
	//   |    \   |
	//   2 --4-- 3
	//   |       |
	//   8      12
	//   |       |
	//   4 ------
	graph := NewGraph(5)
	graph.AddEdge(0, 1, 10.0)
	graph.AddEdge(0, 2, 6.0)
	graph.AddEdge(0, 3, 5.0)
	graph.AddEdge(1, 3, 15.0)
	graph.AddEdge(2, 3, 4.0)
	graph.AddEdge(2, 4, 8.0)
	graph.AddEdge(3, 4, 12.0)

	mst, totalCost := graph.PrimMST()

	// Verificar número de aristas
	expectedEdges := 4 // V-1
	if len(mst) != expectedEdges {
		t.Errorf("MST debería tener %d aristas, tiene %d", expectedEdges, len(mst))
	}

	// Verificar costo total
	// MST óptimo: (2-3: 4) + (0-3: 5) + (0-2: 6) + (2-4: 8) = 23
	// O también: (2-3: 4) + (0-3: 5) + (0-2: 6) + (3-4: 12) = 27
	// Prim desde nodo 0 produce: (0-3: 5) + (3-2: 4) + (2-4: 8) + (0-1: 10) = 27
	expectedCost := 27.0
	if totalCost != expectedCost {
		t.Errorf("Costo total esperado: %.2f, obtenido: %.2f", expectedCost, totalCost)
	}

	// Verificar que el MST es válido
	if !VerifyMSTWithUnionFind(mst, graph.vertices) {
		t.Error("MST no es válido")
	}
}

// TestMSTProperties verifica propiedades generales del MST
func TestMSTProperties(t *testing.T) {
	graph := NewGraph(6)
	graph.AddEdge(0, 1, 7.0)
	graph.AddEdge(0, 3, 5.0)
	graph.AddEdge(1, 2, 8.0)
	graph.AddEdge(1, 3, 9.0)
	graph.AddEdge(1, 4, 7.0)
	graph.AddEdge(2, 4, 5.0)
	graph.AddEdge(3, 4, 15.0)
	graph.AddEdge(3, 5, 6.0)
	graph.AddEdge(4, 5, 8.0)

	mst, totalCost := graph.PrimMST()

	// Propiedad 1: MST debe tener V-1 aristas
	if len(mst) != graph.vertices-1 {
		t.Errorf("MST debe tener %d aristas, tiene %d", graph.vertices-1, len(mst))
	}

	// Propiedad 2: No debe haber ciclos
	uf := NewUnionFind(graph.vertices)
	for _, edge := range mst {
		if !uf.Union(edge.from, edge.to) {
			t.Error("MST contiene un ciclo")
		}
	}

	// Propiedad 3: Todos los nodos deben estar conectados
	root := uf.Find(0)
	for i := 1; i < graph.vertices; i++ {
		if uf.Find(i) != root {
			t.Errorf("Nodo %d no está conectado al resto", i)
		}
	}

	// Propiedad 4: Costo debe ser positivo
	if totalCost <= 0 {
		t.Error("Costo total debe ser positivo")
	}

	// Propiedad 5: Costo MST debe ser menor que todas las aristas
	totalAllEdges := graph.CalculateTotalCost()
	if totalCost > totalAllEdges {
		t.Error("Costo del MST no puede ser mayor que la suma de todas las aristas")
	}
}

// TestEmptyGraph prueba el comportamiento con grafo vacío
func TestEmptyGraph(t *testing.T) {
	graph := NewGraph(0)
	mst, cost := graph.PrimMST()

	if len(mst) != 0 {
		t.Error("MST de grafo vacío debería estar vacío")
	}

	if cost != 0 {
		t.Error("Costo de grafo vacío debería ser 0")
	}
}

// TestSingleNode prueba el comportamiento con un solo nodo
func TestSingleNode(t *testing.T) {
	graph := NewGraph(1)
	mst, cost := graph.PrimMST()

	if len(mst) != 0 {
		t.Error("MST de un solo nodo debería estar vacío")
	}

	if cost != 0 {
		t.Error("Costo de un solo nodo debería ser 0")
	}
}

// TestDisconnectedGraph prueba con grafo desconectado
func TestDisconnectedGraph(t *testing.T) {
	graph := NewGraph(4)
	// Solo conectar 0-1 y 2-3 (dos componentes separadas)
	graph.AddEdge(0, 1, 5.0)
	graph.AddEdge(2, 3, 5.0)

	mst, _ := graph.PrimMST()

	// MST solo puede conectar nodos alcanzables desde el inicio
	// Debería tener menos de V-1 aristas
	if len(mst) >= graph.vertices-1 {
		t.Error("Grafo desconectado no debería producir MST completo")
	}
}

// BenchmarkPrimSmall benchmark para grafo pequeño
func BenchmarkPrimSmall(b *testing.B) {
	graph := NewGraph(5)
	graph.AddEdge(0, 1, 10.0)
	graph.AddEdge(0, 2, 6.0)
	graph.AddEdge(0, 3, 5.0)
	graph.AddEdge(1, 3, 15.0)
	graph.AddEdge(2, 3, 4.0)
	graph.AddEdge(2, 4, 8.0)
	graph.AddEdge(3, 4, 12.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.PrimMST()
	}
}

// BenchmarkPrimMedium benchmark para grafo mediano
func BenchmarkPrimMedium(b *testing.B) {
	graph := NewGraph(100)
	// Crear grafo denso
	for i := 0; i < 100; i++ {
		for j := i + 1; j < 100; j++ {
			if (i+j)%3 == 0 { // Conectar algunos pares
				graph.AddEdge(i, j, float64(i+j))
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.PrimMST()
	}
}

// BenchmarkUnionFind benchmark para Union-Find
func BenchmarkUnionFind(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uf := NewUnionFind(1000)
		for j := 0; j < 999; j++ {
			uf.Union(j, j+1)
		}
		uf.Find(999)
	}
}
