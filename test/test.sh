rm action
curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
sleep 1 ; echo ""
curl -XPOST http://localhost:8080/init -d @hello.json
sleep 1 
curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
sleep 1 ; echo ""
curl -XPOST http://localhost:8080/init -d @ciao.json
sleep 1
curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
