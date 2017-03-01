TODO
====

* write this readme, with json schemas
* howto add new scripts
* write systemd.service


EXAMPLES
========

```
mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"colorfade","value":0.5}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"colorfade","value":0.2}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":""}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"ceilingsinus"}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"ceilingsinus","b":{"amplitude":200,"offset":400}}'

```
