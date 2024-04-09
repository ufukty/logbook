all: \
	simplified.png \
	branching.png 

branching.png: branching.mmd
	mmdc -f -e png -s 3 -i "branching.mmd" -o "branching.png"

simplified.png: simplified.d2
	d2 --layout elk "simplified.d2" "simplified.png"