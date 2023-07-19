# flyover-harmony

work in progress/proof of concept :)

eventually, another terminal client over the web

## roadmap/ideals

- many connections/processes simultaneously. go concurrency
- low latency in input-response cycle. WebSockets over ajax polling, other options
- buffering screen updates. reduce accumulated lag from forcing all screen updates to be shown in browser in sequence
- full browser integration/support (mouse, encoding, copy/paste, mobile)
- minimize lines of codes, features, third party libs, complexity within reason. maximize feel
