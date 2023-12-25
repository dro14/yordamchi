import os


def count_lines_in_file(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        num_lines = 0
        for line in file:
            line = line.strip()
            if len(line) > 0 and not line.startswith("//"):
                num_lines += 1
        return num_lines


def count_lines_in_directory(directory_path, file_extension):
    num_lines = 0
    for root, dirs, files in os.walk(directory_path):
        for file in files:
            if file.endswith(file_extension):
                file_path = os.path.join(root, file)
                num_lines += count_lines_in_file(file_path)
    return num_lines

total_lines = count_lines_in_directory(".", ".go")
print(f"Total number of lines in *.go files: {total_lines}")
