## Opsfilegen
Generates an opsfile from the diff between two bosh manifests.

### Intentional Limitations
* Will not parse YAML matrices (arrays that have at least one array as an element). In practice this rarely exists in bosh manifests.

### Temporary Limitations
* *Extremely* hackday code. No tests; one big main.go.
* Only creates remove definitions.
* Will only remove array elements if they are hashes; will not create remove definitions for strings or integers, which need to be referenced positionally.
* Does not prioritize the "name" key when looking for a unique identifier for an array item, so paths sometimes have patterns like "/instances/vm_type=c3.large/stemcell"
