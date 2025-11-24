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

// Arista representa una conexión entre dos edificios/nodos
type Arista struct {
	desde int
	hacia int
	costo float64
}

// ColaPrioridad implementa heap.Interface para el algoritmo de Prim
type ColaPrioridad []Arista

func (cp ColaPrioridad) Len() int           { return len(cp) }
func (cp ColaPrioridad) Less(i, j int) bool { return cp[i].costo < cp[j].costo }
func (cp ColaPrioridad) Swap(i, j int)      { cp[i], cp[j] = cp[j], cp[i] }

func (cp *ColaPrioridad) Push(x interface{}) {
	*cp = append(*cp, x.(Arista))
}

func (cp *ColaPrioridad) Pop() interface{} {
	viejo := *cp
	n := len(viejo)
	elem := viejo[n-1]
	*cp = viejo[0 : n-1]
	return elem
}

// EncontUnion estructura para optimización
type EncontUnion struct {
	padre []int
	rango []int
}

// NuevoEncontUnion crea una nueva estructura Union-Find
func NuevoEncontUnion(n int) *EncontUnion {
	eu := &EncontUnion{
		padre: make([]int, n),
		rango: make([]int, n),
	}
	for i := 0; i < n; i++ {
		eu.padre[i] = i
		eu.rango[i] = 0
	}
	return eu
}

// Encontrar encuentra el representante del conjunto con compresión de camino
func (eu *EncontUnion) Encontrar(x int) int {
	if eu.padre[x] != x {
		eu.padre[x] = eu.Encontrar(eu.padre[x]) // Compresión de camino
	}
	return eu.padre[x]
}

// Unir une dos conjuntos usando unión por rango
func (eu *EncontUnion) Unir(x, y int) bool {
	raizX := eu.Encontrar(x)
	raizY := eu.Encontrar(y)

	if raizX == raizY {
		return false // Ya están en el mismo conjunto
	}

	// Unión por rango
	if eu.rango[raizX] < eu.rango[raizY] {
		eu.padre[raizX] = raizY
	} else if eu.rango[raizX] > eu.rango[raizY] {
		eu.padre[raizY] = raizX
	} else {
		eu.padre[raizY] = raizX
		eu.rango[raizX]++
	}
	return true
}

// Grafo representa el grafo de edificios
type Grafo struct {
	vertices int
	aristas  []Arista
	listaAdy map[int][]Arista
}

// NuevoGrafo crea un nuevo grafo
func NuevoGrafo(vertices int) *Grafo {
	return &Grafo{
		vertices: vertices,
		aristas:  make([]Arista, 0),
		listaAdy: make(map[int][]Arista),
	}
}

// AgregarArista agrega una arista al grafo
func (g *Grafo) AgregarArista(desde, hacia int, costo float64) {
	arista := Arista{desde, hacia, costo}
	g.aristas = append(g.aristas, arista)
	g.listaAdy[desde] = append(g.listaAdy[desde], arista)
	g.listaAdy[hacia] = append(g.listaAdy[hacia], Arista{hacia, desde, costo})
}

// PrimAEM implementa el algoritmo de Prim para encontrar el Árbol de Expansión Mínimo
func (g *Grafo) PrimAEM() ([]Arista, float64) {
	if g.vertices == 0 {
		return nil, 0
	}

	aem := make([]Arista, 0)
	visitado := make([]bool, g.vertices)
	costoTotal := 0.0

	// Iniciar desde el nodo 0
	visitado[0] = true
	cp := &ColaPrioridad{}
	heap.Init(cp)

	// Agregar todas las aristas del nodo inicial
	for _, arista := range g.listaAdy[0] {
		heap.Push(cp, arista)
	}

	// Procesar mientras haya aristas y no hayamos visitado todos los nodos
	for cp.Len() > 0 && len(aem) < g.vertices-1 {
		arista := heap.Pop(cp).(Arista)

		// Si el destino ya fue visitado, saltar
		if visitado[arista.hacia] {
			continue
		}

		// Agregar arista al AEM
		aem = append(aem, arista)
		costoTotal += arista.costo
		visitado[arista.hacia] = true

		// Agregar todas las aristas del nuevo nodo visitado
		for _, siguienteArista := range g.listaAdy[arista.hacia] {
			if !visitado[siguienteArista.hacia] {
				heap.Push(cp, siguienteArista)
			}
		}
	}

	return aem, costoTotal
}

// CalcularCostoTotal calcula el costo total si se conectaran todos contra todos
func (g *Grafo) CalcularCostoTotal() float64 {
	total := 0.0
	for _, arista := range g.aristas {
		total += arista.costo
	}
	return total
}

// LeerArchivoMTX lee el archivo .mtx y construye el grafo
func LeerArchivoMTX(archivo string) (*Grafo, error) {
	file, err := os.Open(archivo)
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

	grafo := NuevoGrafo(vertices)

	// Semilla para generar costos aleatorios (ya que el grafo es unweighted)
	rand.Seed(time.Now().UnixNano())

	// Leer las aristas
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			desde, _ := strconv.Atoi(parts[0])
			hacia, _ := strconv.Atoi(parts[1])

			// Convertir a índice 0-based
			desde--
			hacia--

			// Generar un costo aleatorio entre 1 y 100
			// Para un caso más realista, podríamos usar la distancia euclidiana
			// pero como no tenemos coordenadas, usamos valores aleatorios
			costo := rand.Float64()*99 + 1

			grafo.AgregarArista(desde, hacia, costo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("Grafo cargado: %d nodos, %d aristas\n", vertices, edges)
	return grafo, nil
}

// VerificarAEMConEncontUnion verifica que el AEM no tenga ciclos usando Union-Find
func VerificarAEMConEncontUnion(aem []Arista, vertices int) bool {
	eu := NuevoEncontUnion(vertices)

	for _, arista := range aem {
		if !eu.Unir(arista.desde, arista.hacia) {
			fmt.Printf("¡Ciclo detectado en arista %d-%d!\n", arista.desde, arista.hacia)
			return false
		}
	}

	// Verificar que todos los nodos estén conectados
	raiz := eu.Encontrar(0)
	for i := 1; i < vertices; i++ {
		if eu.Encontrar(i) != raiz {
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
	archivo := "power-US-Grid.mtx"
	fmt.Printf("Cargando datos desde %s...\n", archivo)

	grafo, err := LeerArchivoMTX(archivo)
	if err != nil {
		fmt.Printf("Error al leer el archivo: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("Ejecutando algoritmo de Prim...")
	inicio := time.Now()

	aem, costoMin := grafo.PrimAEM()

	transcurrido := time.Since(inicio)

	fmt.Println()
	fmt.Println("=== RESULTADOS ===")
	fmt.Printf("Tiempo de ejecución: %v\n", transcurrido)
	fmt.Printf("Costo total mínimo del AEM: %.2f\n", costoMin)
	fmt.Printf("Número de conexiones en AEM: %d\n", len(aem))
	fmt.Println()

	// Calcular costo total de todas las conexiones
	totalTodasAristas := grafo.CalcularCostoTotal()
	fmt.Printf("Costo total si se conectaran todos contra todos: %.2f\n", totalTodasAristas)
	fmt.Printf("Ahorro usando AEM: %.2f (%.2f%%)\n",
		totalTodasAristas-costoMin,
		(totalTodasAristas-costoMin)/totalTodasAristas*100)
	fmt.Println()

	// Verificar AEM con Union-Find
	fmt.Println("Verificando AEM con Union-Find...")
	if VerificarAEMConEncontUnion(aem, grafo.vertices) {
		fmt.Println("✓ AEM válido: sin ciclos y todos los nodos conectados")
	} else {
		fmt.Println("✗ AEM inválido")
	}
	fmt.Println()

	// Mostrar las primeras 20 conexiones del AEM
	fmt.Println("Primeras 20 conexiones a instalar:")
	for i, arista := range aem {
		if i >= 20 {
			break
		}
		fmt.Printf("%3d. Edificio %4d <-> Edificio %4d (Costo: %.2f)\n",
			i+1, arista.desde+1, arista.hacia+1, arista.costo)
	}

	if len(aem) > 20 {
		fmt.Printf("... y %d conexiones más\n", len(aem)-20)
	}

	fmt.Println()
	fmt.Println("=== ANÁLISIS DE COMPLEJIDAD ===")
	fmt.Println("Complejidad temporal: O(E log V)")
	fmt.Printf("  - E (aristas): %d\n", len(grafo.aristas))
	fmt.Printf("  - V (vértices): %d\n", grafo.vertices)
	fmt.Printf("  - E log V ≈ %d * log2(%d) ≈ %.0f operaciones\n",
		len(grafo.aristas), grafo.vertices,
		float64(len(grafo.aristas))*logBase2(float64(grafo.vertices)))
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
