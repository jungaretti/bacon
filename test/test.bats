setup() {
    make
}

@test "ping" {
    ./bin/bacon ping
}

@test "deploy" {
    ./bin/bacon deploy config.example.yml | grep '4 unchanged'
}

@test "print" {
    ./bin/bacon print bacontest42.com
}
