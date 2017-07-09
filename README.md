##Opsfilegen
Generates an opsfile from the diff between two bosh manifests.

### Intentional Limitations
* Will not parse YAML matrices (arrays that have at least one array as an element). In practice this rarely exists in bosh manifests.

### Temporary Limitations
* *Extremely* hackday code. No tests; one big main.go.
* Only creates remove definitions.
* Will only parse arrays whose elements are hashes, strings, or integers.
