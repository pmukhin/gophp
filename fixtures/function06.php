<?php

function helloWorld(String $name) {
    println("Hello " + $name)
}

$helloFunc = helloWorld

foreach ["Pavel", "Kristina"] as $index => $value {
    function(int $i) { print("" + $i + ": ") }($index)
    $helloFunc($value)
}
