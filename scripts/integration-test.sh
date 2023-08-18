set -e

echo 'Building Bacon'
make
echo

echo 'test: bacon ping works'
./bin/bacon ping || (echo 'test: FAIL: ping did not work'; exit 1)
echo 'test: SUCCESS'
echo

echo 'test: bacon deploy would not do anything'
./bin/bacon deploy config.example.yml || (echo 'test: FAIL: deploy did not work'; exit 1)
./bin/bacon deploy config.example.yml | grep 'Would delete 0 records' || (echo 'test: FAIL: Would delete 1+ records'; exit 1)
./bin/bacon deploy config.example.yml | grep 'Would create 0 records' || (echo 'test: FAIL: Would create 1+ records'; exit 1)
echo 'test: SUCCESS'
echo

echo 'test: bacon deploy works'
./bin/bacon deploy --delete --create config.example.yml || (echo 'test: FAIL: deploy --delete --create did not work'; exit 1)
echo 'test: SUCCESS'
