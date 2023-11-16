import os


def count_lines_in_file(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        counter = 0
        for line in file:
            line = line.strip()
            if len(line) > 0 and not line.startswith("//"):
                counter += 1
        return counter


def count_lines_in_directory(directory_path, file_extension):
    total_lines = 0
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith(file_extension):
                file_path = os.path.join(root, file)
                total_lines += count_lines_in_file(file_path)
    return total_lines

# Define your project's root directory (where you want to start counting lines)
project_root = '.'

# Define the file extension you want to count lines for
file_extension = '.go'

# Count the lines and print the result
total_lines = count_lines_in_directory(project_root, file_extension)
print(f'Total number of lines in {file_extension} files: {total_lines}')
