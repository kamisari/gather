gather
======
minimal downloader

Usage:
------
1. make download list. see `example.list`.
```sh
touch gather.list
vim gather.list
```

2. check dry-run.
```sh
gather -list "/path/gather.list" -dir "/path/out/dir"
# or
gather -dir "/path/out/dir" -- "/path/gather.list"
```

3. after check than run it
```sh
gather -list "/path/gather.list" -dir "/path/out/dir" -dryrun=false
# or
gather -dry-run=false -dir "/path/out/dir" -- "/path/gather.list"
```

Options:
--------
```sh
gather -help
```

Install:
--------
```sh
go get -v -u github.com/yaeshimo/gather
```

License:
--------
MIT
