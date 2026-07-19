setup() {
    make
}

@test "ping" {
    ./bin/bacon ping
}

@test "deploy --output table" {
    ./bin/bacon deploy --output table config.example.yml
}

@test "deploy --output json" {
    ./bin/bacon deploy --output json config.example.yml | jq
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

@test "print" {
    ./bin/bacon print bacondemo.com
}
