all: \
	anti-csrf-token-generation.png \
	service-discovery.png \
	server-client.png \
	placement-cache.png

service-discovery.png: service-discovery.mmd
	mmdc -f -e png -s 3 -i "service-discovery.mmd" -o "service-discovery.png"

anti-csrf-token-generation.png: anti-csrf-token-generation.mmd
	mmdc -f -e png -s 3 -i "anti-csrf-token-generation.mmd" -o "anti-csrf-token-generation.png"

placement-cache.png: placement-cache.mmd
	mmdc -f -e png -s 3 -i "placement-cache.mmd" -o "placement-cache.png"

server-client.png: server-client.mmd
	mmdc -f -e png -s 5 -i "server-client.mmd" -o "server-client.png"
