### Intentional Limitations
* Will not parse YAML matrices (arrays that have at least one array as an element). In practice this rarely exists in bosh manifests.

### Temporary Limitations
* Extremely hackday code.
* Will only parse arrays whose elements are hashes.
* Will only create remove definitions.
