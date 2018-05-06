# SqlTune
this project is inspire from eversql , suggestion sql rewrite for mysql 

sqltune use some source of [vitess](https://github.com/vitessio/vitess)

## RUN

Set The GOPATH

* windows

`set GOPATH=YOUR_GO_DIR;YOUR_PROJECT_DIR`

* linux

  `export GOPATH=YOUR_GO_DIR:YOUR_PROJECT_DIR`

build

`go build -i main\main.go`

run
`main -insql "select * from t where a*2 = 2"`

restult:

```
Output sql is:

select * from t where a = (2) / 2

```
## TODO 


## License

Unless otherwise noted, the SqlTune source files are distributed
under the Apache Version 2.0 license found in the LICENSE file.
