#!/usr/bin/env bash

set -e

export PATH=/usr/local/ruby/bin:/usr/local/go/bin:$PATH

echo 'Starting DB...'
su postgres -c '
  export PATH=/usr/lib/postgresql/9.4/bin:$PATH
  export PGDATA=/tmp/postgres
  export PGLOGS=/tmp/log/postgres
  mkdir -p $PGDATA
  mkdir -p $PGLOGS
  initdb -U postgres -D $PGDATA
  pg_ctl start -l $PGLOGS/server.log
'

source /etc/profile.d/chruby.sh
chruby 2.1.6

echo 'Installing dependencies...'
(
	cd bosh-src
	bundle install --local
	bundle exec rake spec:integration:install_dependencies
)	

echo 'Running tests...'

export GOPATH=$(realpath bosh-load-tests)

go run bosh-load-tests/src/github.com/cloudfoundry-incubator/bosh-load-tests/main.go bosh-load-tests/ci/concourse-config.json
