# render-go-forward-slow

* Listens on `:10000`
* TCP tunnels to `$TARGET`
* Toggles delay on reads from `$TARGET` between `0s` and `15s` upon `SIGHUP`
