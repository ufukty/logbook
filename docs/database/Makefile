all: \
	groups.png \
	objectives.png \
	transactions.png

groups.png: groups.d2
	d2 --layout elk "groups.d2" "groups.png"
	magick groups.png -channel RGB -negate -brightness-contrast -10x-30  groups.dark.png

objectives.png: objectives.d2
	d2 --layout elk "objectives.d2" "objectives.png"
	magick objectives.png -channel RGB -negate -brightness-contrast -10x-30  objectives.dark.png

transactions.png: transactions.mmd
	mmdc -f -e png -s 3 -i "transactions.mmd" -o "transactions.png"