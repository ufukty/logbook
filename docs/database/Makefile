all: \
	simplified.png \
	simplified.d2.png \
	versioned.png 

simplified.png: simplified.mmd
	mmdc -f -e png -s 3 -i "simplified.mmd" -o "simplified.png"

versioned.png: versioned.mmd
	mmdc -f -e png -s 3 -i "versioned.mmd" -o "versioned.png"

simplified.d2.png: simplified.d2
	d2 --layout elk "simplified.d2" "simplified.d2.png"