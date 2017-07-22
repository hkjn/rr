# rr

Repo rr holds repo report tools.

## Running client

There's a `report_report.service`, which can be run under systemd:

```
$ sudo cp repo_report.{service,timer} /usr/lib/systemd/user/
$ systemctl --user start repo_report.timer
```