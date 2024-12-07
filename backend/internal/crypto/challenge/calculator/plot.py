import json
import matplotlib.pyplot as plt

# Load data from a JSON file
file_path = "dataset/values.json"  # Replace with your JSON file name or path
with open(file_path, "r") as file:
    data = json.load(file)

# bigger terms are not going to be useful
cutoff = 20

# Extract values
a_values = [entry["a"] for entry in data][:cutoff]
x_values = [entry["x"] for entry in data][:cutoff]
d_values = [entry["d"] for entry in data][:cutoff]
indices = list(range(cutoff))

# combined
plt.figure(figsize=(10, 5))
plt.plot(indices, d_values, marker="o", linestyle="", color="green", label="d")
plt.plot(indices, x_values, marker="o", linestyle="", color="orange", label="x")
plt.plot(indices, a_values, marker="o", linestyle="", label="a")
plt.title("Growth of a, x, d values combined")
plt.xlabel("indices")
plt.ylabel("values")
plt.grid()
plt.legend()
plt.savefig(f"graphs/c{cutoff}.png")
plt.close()
