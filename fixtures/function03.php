<?php

use math\random;

function recursive() {
    if random() % 3 == 0 { return }
    recursive()
}

recursive()