setup() {
    make
}

@test "ping" {
    ./bin/bacon ping
}

@test "deploy preview" {
    ./bin/bacon deploy config.example.yml
}
