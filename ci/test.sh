# automated testing with coveralls support, based on
# https://github.com/mlafeldt/chef-runner/blob/v0.7.0/script/coverage

set -e

rm -f .coverage.out
echo "mode: count" > .coverage.out

for pkg in $(go list ./...); do
  go test -covermode=count -coverprofile=.coverage.tmp "$pkg"
  grep -h -v "^mode:" .coverage.tmp >> .coverage.out || :
  rm -f .coverage.tmp
done

$HOME/gopath/bin/goveralls -coverprofile=.coverage.out -service=travis-ci

rm -f .coverage.out

