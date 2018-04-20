# gophp
Dialect of PHP written in Go with many modern features 

## Exapmples:
```php
function is<T>($o): Boolean { $o instanceof T }
$array = [1, 2, 3]
$booleanVar = if (is<Iterable>($array)) { someResult() } else { someOtherResult() }
```
