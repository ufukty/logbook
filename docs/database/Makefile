all: \
	simplified.png \
	versioned.png 

versioned.png: versioned.mmd
	mmdc -f -e png -s 3 -i "versioned.mmd" -o "versioned.png"

simplified.png: simplified.d2
	d2 --layout elk "simplified.d2" "simplified.png"