<?php

namespace main

class Dog {
    public function bark() {
        println("bark!")
    }
}

$dog = new Dog

foreach 0..5 as $i {
    $dog->bark()
}