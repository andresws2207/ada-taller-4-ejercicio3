package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Edge representa una conexión entre dos edificios/nodos
type Edge struct {
	from int
	to   int
	cost float64
}

// PriorityQueue implementa heap.Interface para Prim's algorithm
type PriorityQueue []Edge

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Edge))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// UnionFind estructura para optimización
type UnionFind struct {
	parent []int
	rank   []int
}

// NewUnionFind crea una nueva estructura Union-Find
func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 0
	}
	return uf
}

// Find encuentra el representante del conjunto con compresión de camino
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Compresión de camino
	}
	return uf.parent[x]
}

// Union une dos conjuntos usando unión por rango
func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // Ya están en el mismo conjunto
	}

	// Unión por rango
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	return true
}

// Graph representa el grafo de edificios
type Graph struct {
	vertices int
	edges    []Edge
	adjList  map[int][]Edge
}

// NewGraph crea un nuevo grafo
func NewGraph(vertices int) *Graph {
	return &Graph{
		vertices: vertices,
		edges:    make([]Edge, 0),
		adjList:  make(map[int][]Edge),
	}
}

// AddEdge agrega una arista al grafo
func (g *Graph) AddEdge(from, to int, cost float64) {
	edge := Edge{from, to, cost}
	g.edges = append(g.edges, edge)
	g.adjList[from] = append(g.adjList[from], edge)
	g.adjList[to] = append(g.adjList[to], Edge{to, from, cost})
}

// PrimMST implementa el algoritmo de Prim para encontrar el MST
func (g *Graph) PrimMST() ([]Edge, float64) {
	if g.vertices == 0 {
		return nil, 0
	}

	mst := make([]Edge, 0)
	visited := make([]bool, g.vertices)
	totalCost := 0.0

	// Iniciar desde el nodo 0
	visited[0] = true
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Agregar todas las aristas del nodo inicial
	for _, edge := range g.adjList[0] {
		heap.Push(pq, edge)
	}

	// Procesar mientras haya aristas y no hayamos visitado todos los nodos
	for pq.Len() > 0 && len(mst) < g.vertices-1 {
		edge := heap.Pop(pq).(Edge)

		// Si el destino ya fue visitado, saltar
		if visited[edge.to] {
			continue
		}

		// Agregar arista al MST
		mst = append(mst, edge)
		totalCost += edge.cost
		visited[edge.to] = true

		// Agregar todas las aristas del nuevo nodo visitado
		for _, nextEdge := range g.adjList[edge.to] {
			if !visited[nextEdge.to] {
				heap.Push(pq, nextEdge)
			}
		}
	}

	return mst, totalCost
}

// CalculateTotalCost calcula el costo total si se conectaran todos contra todos
func (g *Graph) CalculateTotalCost() float64 {
	total := 0.0
	for _, edge := range g.edges {
		total += edge.cost
	}
	return total
}

// ParseMTXFile lee el archivo .mtx y construye el grafo
func ParseMTXFile(filename string) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var vertices, edges int

	// Leer encabezados y comentarios
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "%") {
			continue // Saltar comentarios
		}
		// Primera línea sin comentario contiene: vertices vertices edges
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			vertices, _ = strconv.Atoi(parts[0])
			edges, _ = strconv.Atoi(parts[2])
			break
		}
	}

	graph := NewGraph(vertices)

	// Semilla para generar costos aleatorios (ya que el grafo es unweighted)
	rand.Seed(time.Now().UnixNano())

	// Leer las aristas
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			from, _ := strconv.Atoi(parts[0])
			to, _ := strconv.Atoi(parts[1])

			// Convertir a índice 0-based
			from--
			to--

			// Generar un costo aleatorio entre 1 y 100
			// Para un caso más realista, podríamos usar la distancia euclidiana
			// pero como no tenemos coordenadas, usamos valores aleatorios
			cost := rand.Float64()*99 + 1

			graph.AddEdge(from, to, cost)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("Grafo cargado: %d nodos, %d aristas\n", vertices, edges)
	return graph, nil
}

// VerifyMSTWithUnionFind verifica que el MST no tenga ciclos usando Union-Find
func VerifyMSTWithUnionFind(mst []Edge, vertices int) bool {
	uf := NewUnionFind(vertices)

	for _, edge := range mst {
		if !uf.Union(edge.from, edge.to) {
			fmt.Printf("¡Ciclo detectado en arista %d-%d!\n", edge.from, edge.to)
			return false
		}
	}

	// Verificar que todos los nodos estén conectados
	root := uf.Find(0)
	for i := 1; i < vertices; i++ {
		if uf.Find(i) != root {
			fmt.Printf("¡Grafo no conectado! Nodo %d aislado\n", i)
			return false
		}
	}

	return true
}

func main() {
	fmt.Println("=== Ejercicio 3: Red Eléctrica Óptima (Prim) ===")
	fmt.Println()

	// Parsear el archivo
	filename := "power-US-Grid.mtx"
	fmt.Printf("Cargando datos desde %s...\n", filename)

	graph, err := ParseMTXFile(filename)
	if err != nil {
		fmt.Printf("Error al leer el archivo: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Ejecutando algoritmo de Prim...")
	start := time.Now()

	mst, minCost := graph.PrimMST()

	elapsed := time.Since(start)

	fmt.Println()
	fmt.Println("=== RESULTADOS ===")
	fmt.Printf("Tiempo de ejecución: %v\n", elapsed)
	fmt.Printf("Costo total mínimo del MST: %.2f\n", minCost)
	fmt.Printf("Número de conexiones en MST: %d\n", len(mst))
	fmt.Println()

	// Calcular costo total de todas las conexiones
	totalAllEdges := graph.CalculateTotalCost()
	fmt.Printf("Costo total si se conectaran todos contra todos: %.2f\n", totalAllEdges)
	fmt.Printf("Ahorro usando MST: %.2f (%.2f%%)\n",
		totalAllEdges-minCost,
		(totalAllEdges-minCost)/totalAllEdges*100)
	fmt.Println()

	// Verificar MST con Union-Find
	fmt.Println("Verificando MST con Union-Find...")
	if VerifyMSTWithUnionFind(mst, graph.vertices) {
		fmt.Println("✓ MST válido: sin ciclos y todos los nodos conectados")
	} else {
		fmt.Println("✗ MST inválido")
	}
	fmt.Println()

	// Mostrar las primeras 20 conexiones del MST
	fmt.Println("Primeras 20 conexiones a instalar:")
	for i, edge := range mst {
		if i >= 20 {
			break
		}
		fmt.Printf("%3d. Edificio %4d <-> Edificio %4d (Costo: %.2f)\n",
			i+1, edge.from+1, edge.to+1, edge.cost)
	}

	if len(mst) > 20 {
		fmt.Printf("... y %d conexiones más\n", len(mst)-20)
	}

	fmt.Println()
	fmt.Println("=== ANÁLISIS DE COMPLEJIDAD ===")
	fmt.Println("Complejidad temporal: O(E log V)")
	fmt.Printf("  - E (aristas): %d\n", len(graph.edges))
	fmt.Printf("  - V (vértices): %d\n", graph.vertices)
	fmt.Printf("  - E log V ≈ %d * log2(%d) ≈ %.0f operaciones\n",
		len(graph.edges), graph.vertices,
		float64(len(graph.edges))*logBase2(float64(graph.vertices)))
	fmt.Println()
	fmt.Println("Union-Find optimizado con:")
	fmt.Println("  - Compresión de camino: O(α(n)) ≈ O(1) amortizado")
	fmt.Println("  - Unión por rango: mejora la eficiencia de las operaciones")
}

// logBase2 calcula el logaritmo en base 2
func logBase2(x float64) float64 {
	if x <= 1 {
		return 0
	}
	count := 0.0
	for x > 1 {
		x /= 2
		count++
	}
	return count
}
