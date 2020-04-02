## cert exp date monitor
use to monitor certficate file exp date

exporter infomation: 

1. from date

2. after date

3. domain

4. cert path

args:

```text
--path  type  dir/file
               Provide certificate file path
```

example:
```shell script
$ docker run -dit w564791/cert-exp-exporter:0.0.3 --path=certs  --path=path/ca.pem
```
