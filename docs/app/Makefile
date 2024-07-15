all: \
	diagram.png \
	anti-csrf-token-generation.png

diagram.png: diagram.mmd
	mmdc -f -e png -s 3 -i "diagram.mmd" -o "diagram.png"

anti-csrf-token-generation.png: anti-csrf-token-generation.mmd
	mmdc -f -e png -s 3 -i "anti-csrf-token-generation.mmd" -o "anti-csrf-token-generation.png"
