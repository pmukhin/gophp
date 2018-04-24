<?php

const Greeting = "Hello"

function helloWorld(String $name) {
    println(Greeting + " " + $name)
}

$helloFunc = helloWorld

foreach ["Pavel", "Kristina"] as $index => $value {
    function(int $i) {
        print($i->__toString() + ": ")
    }($index)
    $helloFunc($value)
}