# scheme

```
globalContext {
  functionTable
  classTable
  interfaceTable
  contantTable
}

fileContext = {
  name: "<string>.php"
  path: "dir"
  namespace: "namespace"
  variables: map[string]string
  uses: map[string]string
}

functionContext {
  fileContext
}
```
