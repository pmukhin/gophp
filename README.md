# gophp
Dialect of PHP written in Go with many modern features 

## Exapmples

### Optional semicolons
```php
$array = [1, 2, 3]
$antotherVar = 365
```

### Last statement is return statement
```php
function makeArray() {
  []
}
```

### Type is a constant object
```php
println(dump(Integer)) // <type 'ClassInteger'>
```

### Top level constants
```php
namespace math

const Pi = 3.14
...

println(math\Pi); // 3.14

```

### Everything (almost) is epxression
```php
$integerVar = try { someHeavyCalculation() } catch (MemoryException $e) { 0 }
```
