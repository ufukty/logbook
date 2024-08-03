all: \
	simplified.png \
	transactions.png

simplified.png: simplified.d2
	d2 --layout elk "simplified.d2" "simplified.png"

transactions.png: transactions.mmd
	mmdc -f -e png -s 3 -i "transactions.mmd" -o "transactions.png"