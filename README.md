## opsfilegen
Generates an opsfile from the diff between two bosh manifests.

### Intentional Limitations
* Will not parse YAML matrices (arrays that have at least one array as an element). In practice this rarely exists in bosh manifests.

### Temporary Limitations
* *Extremely* hackday code. No tests; one big main.go.
* Only creates remove definitions -- should also create replace definitions
* Will only remove array elements if they are hashes; will not create remove definitions for strings or integers, which need to be referenced positionally -- should either replace whole array with new array or remove the element by index (seems dangerous)
* Fails silently and provides an invalid opsfile if YAML variable definition is removed without references being removed -- should fail with descriptive error
* Definition order is random -- should alphabetize definitions based on value of path
