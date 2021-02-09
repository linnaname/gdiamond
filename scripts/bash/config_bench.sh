 cd gdiamond/scripts/lua

 ##server publish config bench
 wrk -t8 -c100 -d30s --latency -s config_publish.lua  http://127.0.0.1:1210

 ##server get config bench
 wrk -t8 -c100 -d30s --latency -s config_get.lua  http://127.0.0.1:1210