all: \
	diagram.png \
	anti-csrf-token-generation.png \
	server-client.png \
	placement-cache.png

diagram.png: diagram.mmd
	mmdc -f -e png -s 3 -i "diagram.mmd" -o "diagram.png"

anti-csrf-token-generation.png: anti-csrf-token-generation.mmd
	mmdc -f -e png -s 3 -i "anti-csrf-token-generation.mmd" -o "anti-csrf-token-generation.png"

placement-cache.png: placement-cache.mmd
	mmdc -f -e png -s 3 -i "placement-cache.mmd" -o "placement-cache.png"

server-client.png: server-client.mmd
	mmdc -f -e png -s 5 -i "server-client.mmd" -o "server-client.png"
