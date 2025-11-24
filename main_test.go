package main

import (
	"testing"
)

func TestUnionFind(t *testing.T) {
	eu := NuevoEncontUnion(5)

	for i := 0; i < 5; i++ {
		if eu.Encontrar(i) != i {
			t.Errorf("Nodo %d debería ser su propio representante", i)
		}
	}

	if !eu.Unir(0, 1) {
		t.Error("Unir(0, 1) debería retornar true")
	}

	if eu.Encontrar(0) != eu.Encontrar(1) {
		t.Error("Nodos 0 y 1 deberían estar en el mismo conjunto")
	}

	if eu.Unir(0, 1) {
		t.Error("Unir(0, 1) debería retornar false la segunda vez")
	}

	eu.Unir(1, 2)
	eu.Unir(3, 4)

	if eu.Encontrar(0) != eu.Encontrar(2) {
		t.Error("Nodos 0 y 2 deberían estar en el mismo conjunto")
	}

	if eu.Encontrar(0) == eu.Encontrar(3) {
		t.Error("Nodos 0 y 3 no deberían estar en el mismo conjunto")
	}
}

func TestSmallGraph(t *testing.T) {
	grafo := NuevoGrafo(5)
	grafo.AgregarArista(0, 1, 10.0)
	grafo.AgregarArista(0, 2, 6.0)
	grafo.AgregarArista(0, 3, 5.0)
	grafo.AgregarArista(1, 3, 15.0)
	grafo.AgregarArista(2, 3, 4.0)
	grafo.AgregarArista(2, 4, 8.0)
	grafo.AgregarArista(3, 4, 12.0)

	aem, costoTotal := grafo.PrimAEM()

	aristasEsperadas := 4
	if len(aem) != aristasEsperadas {
		t.Errorf("AEM debería tener %d aristas, tiene %d", aristasEsperadas, len(aem))
	}

	costoEsperado := 27.0
	if costoTotal != costoEsperado {
		t.Errorf("Costo total esperado: %.2f, obtenido: %.2f", costoEsperado, costoTotal)
	}

	if !VerificarAEMConEncontUnion(aem, grafo.vertices) {
		t.Error("AEM no es válido")
	}
}

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

	if len(aem) != grafo.vertices-1 {
		t.Errorf("AEM debe tener %d aristas, tiene %d", grafo.vertices-1, len(aem))
	}

	eu := NuevoEncontUnion(grafo.vertices)
	for _, arista := range aem {
		if !eu.Unir(arista.desde, arista.hacia) {
			t.Error("AEM contiene un ciclo")
		}
	}

	raiz := eu.Encontrar(0)
	for i := 1; i < grafo.vertices; i++ {
		if eu.Encontrar(i) != raiz {
			t.Errorf("Nodo %d no está conectado al resto", i)
		}
	}

	if costoTotal <= 0 {
		t.Error("Costo total debe ser positivo")
	}

	totalTodasAristas := grafo.CalcularCostoTotal()
	if costoTotal > totalTodasAristas {
		t.Error("Costo del AEM no puede ser mayor que la suma de todas las aristas")
	}
}

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

func TestDisconnectedGraph(t *testing.T) {
	grafo := NuevoGrafo(4)
	grafo.AgregarArista(0, 1, 5.0)
	grafo.AgregarArista(2, 3, 5.0)

	aem, _ := grafo.PrimAEM()

	if len(aem) >= grafo.vertices-1 {
		t.Error("Grafo desconectado no debería producir AEM completo")
	}
}

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

func BenchmarkPrimMedium(b *testing.B) {
	grafo := NuevoGrafo(100)
	for i := 0; i < 100; i++ {
		for j := i + 1; j < 100; j++ {
			if (i+j)%3 == 0 {
				grafo.AgregarArista(i, j, float64(i+j))
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grafo.PrimAEM()
	}
}

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
