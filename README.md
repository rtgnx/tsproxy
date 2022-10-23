## TSEXEC


**TSEXEC** starts socks5 proxy server to tailscale network and executes any
command with following envs:

```
  http_proxy=socks5://...
  https_proxy=socks5://...
  socks5_proxy=socks5://...
```


### Example Usage 
*CURL Example*
```
TSKEY="tskey-xxxxx" tsexec curl http://tailnet-addr/
```
*Ansible*

```
TSKEY="tskey-xxxxx" tsexec ansible-playbook -i inventory.ini playbook.yml --ssh-common-args "-o ProxyCommand=\"/usr/bin/nc -X 5 -x ${socks5_proxy} %h %p\""
```

*Custom Port*
```
TSKEY="tskey-xxxxx" SOCKS_ADDR="127.0.0.1:9999" tsexec curl http://google.com/
```
 




