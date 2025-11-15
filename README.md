# Taller 4 - Análisis y Diseño de Algoritmos

**Estudiante:** Andres Alfredo Wu Solano  
**Fecha:** Noviembre 2025

---

# Ejercicio 3: Red Eléctrica Óptima (Prim)

## Descripción del Problema
Una ciudad necesita conectar N edificios con cables eléctricos. Se tiene el costo de instalar un cable entre cada par de edificios. El objetivo es encontrar el costo mínimo para conectar todos los edificios usando el algoritmo de Prim.

**Input:**
- N edificios
- Lista de posibles conexiones con su costo: (edificio1, edificio2, costo)

**Output:**
- Costo total mínimo
- Lista de conexiones a instalar
- Costo total si se conectaran todos contra todos (para comparación)

## Enfoque de Solución

### Algoritmo de Prim
El algoritmo de Prim es un algoritmo **greedy** (voraz) que construye el árbol de expansión mínimo (MST) agregando iterativamente la arista de menor costo que conecta un nodo visitado con uno no visitado.

**Pasos del algoritmo:**
1. Iniciar con un nodo arbitrario (nodo 0)
2. Marcar el nodo como visitado
3. Agregar todas sus aristas a una cola de prioridad
4. Mientras haya nodos sin visitar:
   - Extraer la arista de menor costo de la cola
   - Si el destino no ha sido visitado:
     - Agregar la arista al MST
     - Marcar el destino como visitado
     - Agregar todas las aristas del nuevo nodo a la cola
5. Retornar el MST y su costo total

### Optimización con Union-Find
La estructura **Union-Find** (Disjoint Set Union) se utiliza para:
- **Verificar** que el MST no contenga ciclos
- **Confirmar** que todos los nodos están conectados
- Operaciones optimizadas con:
  - **Compresión de camino**: Aplana la estructura durante búsquedas
  - **Unión por rango**: Mantiene árboles balanceados

## Complejidad

### Temporal: **O(E log V)**
- **E**: Número de aristas (6,594)
- **V**: Número de vértices (4,941)

**Desglose:**
- Inicialización: O(V)
- Cada arista se procesa máximo una vez: O(E)
- Cada operación de heap (push/pop): O(log E) = O(log V²) = O(log V)
- **Total**: O(E log V) ≈ 6,594 × log₂(4,941) ≈ 85,722 operaciones

**Union-Find:**
- Find con compresión de camino: O(α(n)) ≈ O(1) amortizado
- Union con rango: O(α(n)) ≈ O(1) amortizado
- α(n) es la función inversa de Ackermann (crece extremadamente lento)

### Espacial: **O(V + E)**
- **Grafo**: O(V + E) para lista de adyacencia
- **Heap**: O(E) en el peor caso
- **Arrays auxiliares** (visited, parent, rank): O(V)
- **MST resultante**: O(V-1) = O(V)
- **Total**: O(V + E)

### ✅ Algoritmo de Prim
- Complejidad temporal: **O(E log V)**
- Usa una cola de prioridad (heap) para seleccionar eficientemente la arista de menor costo
- Implementación iterativa que evita recursión innecesaria

### ✅ Union-Find Optimizado
- **Compresión de camino**: Optimiza búsquedas futuras aplanando la estructura
- **Unión por rango**: Mantiene el árbol balanceado para operaciones más rápidas
- Usado para verificar que el MST no contenga ciclos
- Complejidad amortizada: O(α(n)) ≈ O(1) donde α es la función inversa de Ackermann

### ✅ Análisis de Datos
- Parseo del formato Matrix Market (.mtx)
- Generación de costos aleatorios (1-100) para aristas sin peso
- Cálculo del costo total del MST
- Cálculo del costo si se conectaran todas las aristas
- Comparación de ahorro

## Estructura del Código

### Estructuras de Datos

```go
type Edge struct {
    from int
    to   int
    cost float64
}

type Graph struct {
    vertices int
    edges    []Edge
    adjList  map[int][]Edge
}

type UnionFind struct {
    parent []int
    rank   []int
}
```

### Funciones Principales

1. **PrimMST()**: Implementación del algoritmo de Prim
2. **Union-Find Operations**:
   - `Find()`: Encuentra el representante con compresión de camino
   - `Union()`: Une dos conjuntos usando unión por rango
3. **ParseMTXFile()**: Lee y parsea el archivo del dataset
4. **VerifyMSTWithUnionFind()**: Verifica la validez del MST

## Resultados con US Power Grid

### Dataset
- **Nodos**: 4,941 edificios
- **Aristas**: 6,594 conexiones posibles
- **Fuente**: US Power Grid Network (Pajek/UF Sparse Matrix Collection)

### Resultados de Ejecución

```
=== RESULTADOS ===
Tiempo de ejecución: ~5ms
Costo total mínimo del MST: 206,912.11
Número de conexiones en MST: 4,940

Costo total si se conectaran todos contra todos: 334,616.78
Ahorro usando MST: 127,704.66 (38.16%)

Verificación: ✓ MST válido (sin ciclos, todos conectados)
```

### Análisis de Complejidad

- **Operaciones teóricas**: E log V ≈ 6,594 × log₂(4,941) ≈ 85,722 operaciones
- **Tiempo real**: ~5 milisegundos
- **Espacio**: O(V + E) para almacenar el grafo

## Cómo Ejecutar

### Prerequisitos
- Go 1.21 o superior instalado
- Dataset `power-US-Grid.mtx` en el mismo directorio

### Compilar y Ejecutar

```bash
# Navegar al directorio
cd power-US-Grid

# Ejecutar directamente
go run main.go

# O compilar primero y luego ejecutar
go build -o power-grid.exe main.go
./power-grid.exe
```

### Salida Esperada
El programa generará:
1. Información de carga del grafo
2. Tiempo de ejecución del algoritmo
3. Costo total mínimo del MST
4. Número de conexiones necesarias
5. Comparación con costo total de todas las aristas
6. Verificación del MST usando Union-Find
7. Lista de las primeras 20 conexiones a instalar
8. Análisis de complejidad

## Casos de Prueba

### Caso 1: US Power Grid Dataset (Principal)
**Descripción:** Red eléctrica de Estados Unidos con 4,941 nodos y 6,594 aristas

**Fuente:** UF Sparse Matrix Collection (Pajek/USpowerGrid)

**Características:**
- Grafo no dirigido
- Grafo no ponderado (se generan pesos aleatorios 1-100)
- Representa conexiones reales de la red eléctrica de EE.UU.

**Resultados Esperados:**
```
Grafo cargado: 4941 nodos, 6594 aristas
Tiempo de ejecución: ~5ms
Costo total mínimo del MST: ~200,000-210,000 (varía por pesos aleatorios)
Número de conexiones en MST: 4940 (V-1)
Ahorro vs todas las aristas: ~35-40%
Verificación: ✓ MST válido
```

### Caso 2: Grafo Pequeño (Ejemplo Manual)
Para probar con un grafo más pequeño, puedes modificar el código para usar este ejemplo:

```go
// Ejemplo: 5 edificios con 7 conexiones posibles
graph := NewGraph(5)
graph.AddEdge(0, 1, 10.0)
graph.AddEdge(0, 2, 6.0)
graph.AddEdge(0, 3, 5.0)
graph.AddEdge(1, 3, 15.0)
graph.AddEdge(2, 3, 4.0)
graph.AddEdge(2, 4, 8.0)
graph.AddEdge(3, 4, 12.0)
```

**MST Esperado:**
- Aristas: (2,3,4), (0,3,5), (0,2,6), (2,4,8)
- Costo total: 23
- Conexiones: 4 (V-1)

### Caso 3: Verificación de Propiedades del MST

El programa verifica automáticamente:

1. **Sin ciclos:** Union-Find detecta si agregar una arista crearía un ciclo
2. **Conectividad completa:** Todos los nodos tienen el mismo representante
3. **V-1 aristas:** El MST tiene exactamente V-1 aristas
4. **Minimalidad:** Prim garantiza el costo mínimo

### Pruebas de Rendimiento

| Dataset | Nodos | Aristas | Tiempo | MST Aristas | Complejidad |
|---------|-------|---------|--------|-------------|-------------|
| US Grid | 4,941 | 6,594   | ~5ms   | 4,940       | O(E log V)  |
| Pequeño | 5     | 7       | ~2µs   | 4           | O(E log V)  |
| Mediano | 100   | ~1,100  | ~234µs | 99          | O(E log V)  |

**Benchmarks (CPU: Intel i7-1355U):**
```
BenchmarkPrimSmall-12      512104    1961 ns/op    (~2µs)
BenchmarkPrimMedium-12       4760  234464 ns/op    (~234µs)
BenchmarkUnionFind-12       76516   15622 ns/op    (~16µs para 1000 nodos)
```

### Ejecutar Tests

```bash
# Ejecutar todos los tests
go test -v

# Ejecutar tests con coverage
go test -cover

# Ejecutar benchmarks
go test -bench . -run XXX

# Test específico
go test -run TestSmallGraph -v
```

**Tests Incluidos:**
- ✅ `TestUnionFind`: Verifica operaciones Find/Union
- ✅ `TestSmallGraph`: Prueba con grafo pequeño conocido
- ✅ `TestMSTProperties`: Verifica propiedades del MST (V-1 aristas, sin ciclos, conectividad)
- ✅ `TestEmptyGraph`: Caso borde con grafo vacío
- ✅ `TestSingleNode`: Caso borde con un solo nodo
- ✅ `TestDisconnectedGraph`: Verifica comportamiento con componentes desconectadas

**Resultados de Tests:**
```
=== RUN   TestUnionFind
--- PASS: TestUnionFind (0.00s)
=== RUN   TestSmallGraph
--- PASS: TestSmallGraph (0.00s)
=== RUN   TestMSTProperties
--- PASS: TestMSTProperties (0.00s)
=== RUN   TestEmptyGraph
--- PASS: TestEmptyGraph (0.00s)
=== RUN   TestSingleNode
--- PASS: TestSingleNode (0.00s)
=== RUN   TestDisconnectedGraph
--- PASS: TestDisconnectedGraph (0.00s)
PASS
```

## Output

El programa genera:
1. **Costo total mínimo**: Suma de las aristas en el MST
2. **Lista de conexiones**: Las 4,940 conexiones a instalar
3. **Costo de comparación**: Costo total de todas las aristas posibles
4. **Ahorro**: Diferencia y porcentaje ahorrado
5. **Verificación**: Validación del MST usando Union-Find

## Optimizaciones Implementadas

1. **Cola de Prioridad (Heap)**:
   - Selección eficiente de la arista mínima: O(log E)
   - Implementada usando `container/heap` de Go

2. **Union-Find**:
   - Compresión de camino: reduce altura del árbol
   - Unión por rango: mantiene árboles balanceados
   - Detección rápida de ciclos: O(α(n))

3. **Lista de Adyacencia**:
   - Acceso rápido a vecinos de cada nodo
   - Reduce iteraciones innecesarias

## Complejidad Final

- **Temporal**: O(E log V)
  - Cada arista se procesa una vez: O(E)
  - Cada operación de heap: O(log E) = O(log V²) = O(2 log V) = O(log V)
  - Total: O(E log V)

- **Espacial**: O(V + E)
  - Grafo: O(V + E)
  - Heap: O(E)
  - Union-Find: O(V)
  - Arrays auxiliares: O(V)

## Notas

- El dataset original es un grafo no ponderado, por lo que se generan costos aleatorios (1-100) para simular costos de instalación
- Los resultados varían entre ejecuciones debido a la generación aleatoria de costos
- La estructura Union-Find garantiza que el MST encontrado es válido (sin ciclos y completamente conectado)

---

## Referencias

- Dataset: [UF Sparse Matrix Collection - US Power Grid](http://www.cise.ufl.edu/research/sparse/matrices/Pajek/USpowerGrid)
- Algoritmo de Prim: Introduction to Algorithms (CLRS)
- Union-Find: Algorithms 4th Edition (Sedgewick & Wayne)
