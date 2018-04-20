<?php

use math\random

function recursive() {
    if random() % 3 == 0 { return }
    println(1)
    recursive()
}

recursive()