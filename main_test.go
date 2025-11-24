package main

import (
	"testing"
)

// TestUnionFind verifica la funcionalidad de Union-Find
func TestUnionFind(t *testing.T) {
	eu := NuevoEncontUnion(5)

	// Verificar estado inicial
	for i := 0; i < 5; i++ {
		if eu.Encontrar(i) != i {
			t.Errorf("Nodo %d debería ser su propio representante", i)
		}
	}

	// Unir nodos 0 y 1
	if !eu.Unir(0, 1) {
		t.Error("Unir(0, 1) debería retornar true")
	}

	// Verificar que están en el mismo conjunto
	if eu.Encontrar(0) != eu.Encontrar(1) {
		t.Error("Nodos 0 y 1 deberían estar en el mismo conjunto")
	}

	// Intentar unir de nuevo (debería fallar)
	if eu.Unir(0, 1) {
		t.Error("Unir(0, 1) debería retornar false la segunda vez")
	}

	// Unir más nodos
	eu.Unir(1, 2)
	eu.Unir(3, 4)

	// Verificar conjuntos
	if eu.Encontrar(0) != eu.Encontrar(2) {
		t.Error("Nodos 0 y 2 deberían estar en el mismo conjunto")
	}

	if eu.Encontrar(0) == eu.Encontrar(3) {
		t.Error("Nodos 0 y 3 no deberían estar en el mismo conjunto")
	}
}

// TestSmallGraph prueba Prim con un grafo pequeño conocido
func TestSmallGraph(t *testing.T) {
	// Crear grafo de ejemplo
	grafo := NuevoGrafo(5)
	grafo.AgregarArista(0, 1, 10.0)
	grafo.AgregarArista(0, 2, 6.0)
	grafo.AgregarArista(0, 3, 5.0)
	grafo.AgregarArista(1, 3, 15.0)
	grafo.AgregarArista(2, 3, 4.0)
	grafo.AgregarArista(2, 4, 8.0)
	grafo.AgregarArista(3, 4, 12.0)

	aem, costoTotal := grafo.PrimAEM()

	// Verificar número de aristas
	aristasEsperadas := 4 // V-1
	if len(aem) != aristasEsperadas {
		t.Errorf("AEM debería tener %d aristas, tiene %d", aristasEsperadas, len(aem))
	}

	// Verificar costo total
	// AEM óptimo: (2-3: 4) + (0-3: 5) + (0-2: 6) + (2-4: 8) = 23
	// O también: (2-3: 4) + (0-3: 5) + (0-2: 6) + (3-4: 12) = 27
	// Prim desde nodo 0 produce: (0-3: 5) + (3-2: 4) + (2-4: 8) + (0-1: 10) = 27
	costoEsperado := 27.0
	if costoTotal != costoEsperado {
		t.Errorf("Costo total esperado: %.2f, obtenido: %.2f", costoEsperado, costoTotal)
	}

	// Verificar que el AEM es válido
	if !VerificarAEMConEncontUnion(aem, grafo.vertices) {
		t.Error("AEM no es válido")
	}
}

// TestMSTProperties verifica propiedades generales del AEM
func TestMSTProperties(t *testing.T) {
	grafo := NuevoGrafo(6)
	grafo.AgregarArista(0, 1, 7.0)
	grafo.AgregarArista(0, 3, 5.0)
	grafo.AgregarArista(1, 2, 8.0)
	grafo.AgregarArista(1, 3, 9.0)
	grafo.AgregarArista(1, 4, 7.0)
	grafo.AgregarArista(2, 4, 5.0)
	grafo.AgregarArista(3, 4, 15.0)
	grafo.AgregarArista(3, 5, 6.0)
	grafo.AgregarArista(4, 5, 8.0)

	aem, costoTotal := grafo.PrimAEM()

	// Propiedad 1: AEM debe tener V-1 aristas
	if len(aem) != grafo.vertices-1 {
		t.Errorf("AEM debe tener %d aristas, tiene %d", grafo.vertices-1, len(aem))
	}

	// Propiedad 2: No debe haber ciclos
	eu := NuevoEncontUnion(grafo.vertices)
	for _, arista := range aem {
		if !eu.Unir(arista.desde, arista.hacia) {
			t.Error("AEM contiene un ciclo")
		}
	}

	// Propiedad 3: Todos los nodos deben estar conectados
	raiz := eu.Encontrar(0)
	for i := 1; i < grafo.vertices; i++ {
		if eu.Encontrar(i) != raiz {
			t.Errorf("Nodo %d no está conectado al resto", i)
		}
	}

	// Propiedad 4: Costo debe ser positivo
	if costoTotal <= 0 {
		t.Error("Costo total debe ser positivo")
	}

	// Propiedad 5: Costo AEM debe ser menor que todas las aristas
	totalTodasAristas := grafo.CalcularCostoTotal()
	if costoTotal > totalTodasAristas {
		t.Error("Costo del AEM no puede ser mayor que la suma de todas las aristas")
	}
}

// TestEmptyGraph prueba el comportamiento con grafo vacío
func TestEmptyGraph(t *testing.T) {
	grafo := NuevoGrafo(0)
	aem, costo := grafo.PrimAEM()

	if len(aem) != 0 {
		t.Error("AEM de grafo vacío debería estar vacío")
	}

	if costo != 0 {
		t.Error("Costo de grafo vacío debería ser 0")
	}
}

// TestSingleNode prueba el comportamiento con un solo nodo
func TestSingleNode(t *testing.T) {
	grafo := NuevoGrafo(1)
	aem, costo := grafo.PrimAEM()

	if len(aem) != 0 {
		t.Error("AEM de un solo nodo debería estar vacío")
	}

	if costo != 0 {
		t.Error("Costo de un solo nodo debería ser 0")
	}
}

// TestDisconnectedGraph prueba con grafo desconectado
func TestDisconnectedGraph(t *testing.T) {
	grafo := NuevoGrafo(4)
	// Solo conectar 0-1 y 2-3 (dos componentes separadas)
	grafo.AgregarArista(0, 1, 5.0)
	grafo.AgregarArista(2, 3, 5.0)

	aem, _ := grafo.PrimAEM()

	// AEM solo puede conectar nodos alcanzables desde el inicio
	// Debería tener menos de V-1 aristas
	if len(aem) >= grafo.vertices-1 {
		t.Error("Grafo desconectado no debería producir AEM completo")
	}
}

// BenchmarkPrimSmall benchmark para grafo pequeño
func BenchmarkPrimSmall(b *testing.B) {
	grafo := NuevoGrafo(5)
	grafo.AgregarArista(0, 1, 10.0)
	grafo.AgregarArista(0, 2, 6.0)
	grafo.AgregarArista(0, 3, 5.0)
	grafo.AgregarArista(1, 3, 15.0)
	grafo.AgregarArista(2, 3, 4.0)
	grafo.AgregarArista(2, 4, 8.0)
	grafo.AgregarArista(3, 4, 12.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grafo.PrimAEM()
	}
}

// BenchmarkPrimMedium benchmark para grafo mediano
func BenchmarkPrimMedium(b *testing.B) {
	grafo := NuevoGrafo(100)
	// Crear grafo denso
	for i := 0; i < 100; i++ {
		for j := i + 1; j < 100; j++ {
			if (i+j)%3 == 0 { // Conectar algunos pares
				grafo.AgregarArista(i, j, float64(i+j))
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grafo.PrimAEM()
	}
}

// BenchmarkUnionFind benchmark para Union-Find
func BenchmarkUnionFind(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eu := NuevoEncontUnion(1000)
		for j := 0; j < 999; j++ {
			eu.Unir(j, j+1)
		}
		eu.Encontrar(999)
	}
}
