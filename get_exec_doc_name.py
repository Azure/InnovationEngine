import os
import glob

# Get a list of all markdown files in the AIDocs subfolder and its subdirectories
files = glob.glob('scenarios/ocd/AIDocs/**/*.md', recursive=True)

# Get the most recently created file
latest_file = max(files, key=os.path.getctime)

# Print the name of the most recently created file
print(os.path.basename(latest_file))