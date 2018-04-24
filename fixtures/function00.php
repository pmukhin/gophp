<?php

$zero = 0

/**
 * Fibonacci recursive implementation.
 * @author Pavel Mukhin
 */
function fib(int $n): int {
    if $n < 2 { $n } else { fib($n - 1) + fib($n - 2) }
}

println(fib(0), fib(1), fib(2), fib(3), fib(4), fib(5), fib(6), fib(7))
