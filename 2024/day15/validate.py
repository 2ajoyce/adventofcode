import re

def process_and_validate(file_input, file_output):
    """
    Process the input file to extract coordinates, validate moves, and write valid coordinates to an output file.
    """
    try:
        # Step 1: Read and filter lines containing "Moved fish to [{x y}]"
        valid_lines = []
        with open(file_input, 'r') as infile:
            for line in infile:
                match = re.search(r'Moved fish to \[\{(\d+)\s+(\d+)\}\]', line)
                if match:
                    x, y = int(match.group(1)), int(match.group(2))
                    valid_lines.append((x, y))

        # Step 2: Validate moves and write to the output file
        previous_x, previous_y = None, None
        is_valid = True

        with open(file_output, 'w') as outfile:
            for i, (x, y) in enumerate(valid_lines):
                if previous_x is not None and previous_y is not None:
                    delta_x = abs(x - previous_x)
                    delta_y = abs(y - previous_y)

                    # Validation: no diagonal moves and only one coordinate changes
                    if delta_x > 0 and delta_y > 0:
                        print(f"Error: Diagonal move detected at line {i+1}: [{x}, {y}]")
                        is_valid = False
                    elif delta_x > 1 or delta_y > 1:
                        print(f"Error: Move too large at line {i+1}: [{x}, {y}]")
                        is_valid = False

                # Write valid coordinates to output file
                outfile.write(f"{x} {y}\n")
                previous_x, previous_y = x, y

        # Final validation result
        if is_valid:
            print("Processing and validation successful: No issues detected.")
        else:
            print("Processing completed with errors: Check logs for issues.")

    except FileNotFoundError:
        print(f"Error: File not found - {file_input}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")


if __name__ == "__main__":
    # Input and output file names
    input_file = "run.txt"
    output_file = "run_processed.txt"

    # Process the file and validate moves
    process_and_validate(input_file, output_file)
