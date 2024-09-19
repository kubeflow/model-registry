import sys


def reorder_go_functions(file_path):
    with open(file_path, 'r') as f:
        lines = f.readlines()

    header = []
    functions = []
    current_function = []
    inside_function = False

    for line in lines:
        if line.startswith("// Assert"):
            inside_function = True
        
        if inside_function:
            current_function.append(line)
            if line.rstrip() == "}":
                functions.append("".join(current_function))
                current_function = []
                inside_function = False
        elif len(functions) == 0:
            header.append(line)

    functions.sort(key=lambda func: func.split('\n')[0].strip())

    sorted_content = "".join(header) + "\n".join(functions)

    with open(file_path, 'w') as f:
        f.write(sorted_content)


if __name__ == "__main__":
    if len(sys.argv) == 2:
        reorder_go_functions(sys.argv[1])
    else:
        print("Provide internal/server/openapi/type_asserts.go path as first argument")
