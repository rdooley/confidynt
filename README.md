Confidynt
============

Confident [12 factor config](https://12factor.net/config) in dynamo db CLI tool

### Usage


* Writing a config

```
❯ cat example.conf
environment=confidynt-example
alpha_first="sandwich
    hotdog"
comes_first=fart
other_thing=barbeque
CLOUDWATCH_THING="ref to $environment"
❯ confidynt --table=deployment write example.conf
example.conf written to deployment
```


* Reading a config
```
❯ confidynt --table=deployment read environment confidynt-example
environment=confidynt-example
comes_first=fart
CLOUDWATCH_THING="ref to $environment"
alpha_first="sandwich
    hotdog"
other_thing=barbeque
```
