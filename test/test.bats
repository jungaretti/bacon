setup() {
    make
}

@test "ping" {
    ./bin/bacon ping
}

@test "deploy would delete 0 records" {
    ./bin/bacon deploy config.example.yml | grep 'Would delete 0 records'
}

@test "deploy would create 0 records" {
    ./bin/bacon deploy config.example.yml | grep 'Would create 0 records'
}

@test "deploy" {
    output="$(./bin/bacon deploy config.example.yml)"
    if [[ ! $output == *"Would delete 0 records"* ]]; then
        skip 'Would delete a record'
    fi
    if [[ ! $output == *"Would create 0 records"* ]]; then
        skip 'Would create a record'
    fi

    ./bin/bacon deploy --delete --create config.example.yml
}
