searchmysite-go
===============

This is a front-end to searchmysite.net, initially intended to be used on [seirdy.one](https://seirdy.one/)

More documentation and tests will come.

This is not ready for others to use yet; the API is still very much in flux.

There's a sample Systemd unit file available; tweak it to your needs. The unit file assumes that the binary is statically linked, and sandboxes it in a chroot jail without any shared libraries. `systemd-analyze security searchmysite-go.service` returns a score of <samp>0.9 SAFE 😀</samp>
