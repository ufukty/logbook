all: \
	objectives.png \
	transactions.png

objectives.png: objectives.d2
	d2 --layout elk "objectives.d2" "objectives.png"
	magick objectives.png -channel RGB -negate -brightness-contrast -10x-30  objectives.dark.png

transactions.png: transactions.mmd
	mmdc -f -e png -s 3 -i "transactions.mmd" -o "transactions.png"