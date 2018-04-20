# gophp
Dialect of PHP written in Go with many modern features 

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

### Generics
```php
function is<T>($o): Boolean { $o instanceof T }
```



### Everything (almost) is epxression
```php
$booleanVar = if (is<Iterable>($array)) { someResult() } else { someOtherResult() }
$integerVar = try { someHeavyCalculation() } catch (MemoryException $e) { 0 }
$some = $booleanVar->to<Integer>() + $integerVar
```
