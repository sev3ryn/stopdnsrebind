# [WIP] stopdnsrebind

## Name

*stopdnsrebind* - Coredns plugin that implement `--stop-dns-rebind` from dnsmasq.

## Description

With `stopdnsrebind` enabled, users are able to block addresses from upstream nameservers which are in the private ranges plus ranges specified in public_nets parameter

## Syntax

```
stopdnsrebind {
    public_nets [IP RANGE]
}
```

- **IP RANGE** public ip range not allowed to resolve

## Examples

To demonstrate the usage of plugin stopdnsrebind, here we provide some typical examples.

~~~ corefile
. {
    stopdnsrebind {
        public_nets 8.8.8.0/24
    }
}
~~~
