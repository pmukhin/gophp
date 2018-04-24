use os\args

function run(Array $args) {
    foreach ($args as $value) {
        println($value)
    }
}

function makeArray() {
    $array = []
    foreach 0..25 as $i { $array->append($i) }
    $array
}

foreach (makeArray() as $key => $value) { run(args()) }
