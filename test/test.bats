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

@test "print" {
    ./bin/bacon print bacondemo.com
}
