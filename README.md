# gophp
Dialect of PHP written in Go with many modern features 

## Exapmples

### Optional semicolons
```php
$array = [1, 2, 3]
$antotherVar = 365
```

### Last statement is a return statement
```php
function makeArray(): Array { [] }
println(makeArray()) // []
```

### Everything is a value
```php
function helloWorld(String $name) { println("Hello " + $name) }
$helloFunc = helloWorld

foreach ["Pavel", "Kristina"] as $index => $value {
  function(int $i) { printf("" + $i + ": ") }($index)
  $helloFunc($value)
}

// 0: Hello Pavel
// 1: Hello Kristina
```

### Range operator
```php
foreach 0..3 as $i { println($i) }
// 0
// 1
// 2
// 3
```

### Type is a constant object
```php
println(Integer) // <type 'ClassInteger'>
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
