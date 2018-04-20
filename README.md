# gophp
Dialect of PHP written in Go with many modern features 

## Exapmples

### Generics
```php
function is<T>($o): Boolean { $o instanceof T }
```

### Optional semicolons
```php
$array = [1, 2, 3]
$antotherVar = 365
```

### Everything (almost) is epxression
```php
$booleanVar = if (is<Iterable>($array)) { someResult() } else { someOtherResult() }
$integerVar = try { someHeavyCalculation() } catch (MemoryException $e) { 0 }
$some = $booleanVar->to<Integer>() + $integerVar
```
