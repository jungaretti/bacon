setup() {
    make
}

@test "ping" {
    ./bin/bacon ping
}

@test "deploy" {
    ./bin/bacon deploy config.example.yml
}

@test "print" {
    ./bin/bacon print bacontest42.com
}
