# units

This directory has sample systemd units for running the rr image periodically on a timer.

These systemd units can be added and enabled for your user with:
```
cd ~/.config/systemd/user/
cp $GOPATH/src/hkjn.me/rr/units/repo-report.{service,timer} .
systemctl --user enable repo-report.*
systemctl --user start repo-report.*
```