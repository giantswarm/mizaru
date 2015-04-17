# mizaru
Chaos Monkey for docker containers.

```
 __  __ _
|  \/  (_)______ _ _ __ _   _
| |\/| | |_  / _` | '__| | | |
| |  | | |/ / (_| | |  | |_| |
|_|  |_|_/___\__,_|_|   \__,_|)'
```

Mizaru is a little helper to generate IPTable DROP rules to test resilience against netsplits and other network problems.

Idea inspired (or shamelessly stolen from) [Aphry`s Jepsen](https://github.com/aphyr/jepsen).

## Usage

```
sudo ./mizaru [options] mode [--timeout=seconds] [--block]
```

 * `mode` - See the modelist below. Defaults to `netsplit`
 * `timeout` - A timeout in seconds before reverting the applied rules
 * `block` - Blocks until a signal is received. Reverts rules and quits.

NOTE: you must invoke this with sudo as we are calling the iptables tool.

## Modes

 * `fast` - everything is allowd - all custom iptable rules will be deleted
 * `netsplit` - The docker containers are split into a majority and minority group which cannot communicate with each other
 * `bridge` - Same as `netsplit`, but one container can talk to both sides
 * `single-bridge` - Same as `bridge`, but only one container on each side can talk to each other
 * `ring` - Each container can talk to two other containers so that they form a ring

## Future Ideas

* Add a http-server mode so it can be more easily triggered by load testing tools like gatling.io
* Add modes for trafficshaping, e.g. package loss or latency
* Add mode for pausing containers to simulate GC lags
* Support blacklisting of CIDRs

Add your own ideas as issues [here](http://github.com/giantswarm/mizaru/issues).
