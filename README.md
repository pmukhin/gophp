## gophp is a dialect of PHP written in pure go with many modern features 

## Exapmples

### Optional semicolons
```php
$array = [1, 2, 3]
$antotherVar = 365
```

### Top level constants
```php
namespace math

const Pi = 3.14
...

println(math\Pi); // 3.14

```

### Last statement is a return statement
```php
function makeArray(): Array { [] }
println(makeArray()) // []
```

### Everything is a value
```php
const Greeting = "Hello "

function helloWorld(String $name) { println(Greeting + $name) }
$helloFunc = helloWorld

foreach ["Pavel", "Kristina"] as $index => $value {
  function(int $i) { print("" + $i + ": ") }($index)
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
```

### Type is a constant object
```php
println(Integer) // <type 'ClassInteger'>
```

### Everything (almost) is epxression
```php
$integerVar = try { someHeavyCalculation() } catch (MemoryException $e) { 0 }
```
