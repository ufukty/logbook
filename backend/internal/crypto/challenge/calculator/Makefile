all: dataset/values.json graphs selected.json

dataset/values.json: main.go
	go run .

selected.json: dataset/values.json
	cat "dataset/values.json" | jq -c '.[] | select(.d <= 120)' > selected.json

graphs: plot.py dataset/values.json
	mkdir -p graphs
	python plot.py
