## Dialect of PHP written in Go with many modern features
```php
namespace main

use os\{args, File}

function run(Array $files) {
  foreach $files as $name {
     $file = try { 
       File::open($args) 
     } catch (os\Exception $exception) {
       println("failed to open file: {0}"->format($exception->getMessage()))
       continue
    }
    $content = if ($file->isDir()) {
      "dir: {0}"->format($file->path())
    } else {
      $file->readAll()
    }
    println($content)
  }
}

run(args()[1:])

```
## Why
The project is a sort of a research and there's no aim to prepare a drop-in replacement for current implementations of PHP like native or Hack.

## Motivation
PHP is commonly known as a `bad design fractal` and lacks a lot of progress made in Programming Languages theory for last 20 years. On the other side PHP is still the language of the Web. This project is just one vision of how the issues of the language might be solved.

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
