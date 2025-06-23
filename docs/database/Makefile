all: \
	privileges.png \
	objectives.png \
	transactions.png

privileges.png: privileges.d2
	d2 --layout elk "privileges.d2" "privileges.png"

objectives.png: objectives.d2
	d2 --layout elk "objectives.d2" "objectives.png"

transactions.png: transactions.mmd
	mmdc -f -e png -s 3 -i "transactions.mmd" -o "transactions.png"