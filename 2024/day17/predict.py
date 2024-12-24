import numpy as np
from keras import models
import keras
import tensorflow as tf

# Constants
MAX_SEQUENCE_LENGTH = 20
NUM_CLASSES = 10  # Digits 0-9


def encode_sequence(sequence, max_length=MAX_SEQUENCE_LENGTH):
    """
    Encodes a digit sequence string into a one-hot encoded NumPy array.

    Parameters:
        sequence (str): The digit sequence as a string.
        max_length (int): The fixed length to pad the sequence.

    Returns:
        np.ndarray: One-hot encoded array of shape (max_length, NUM_CLASSES)
    """
    # Ensure the sequence is padded with leading zeros
    sequence = sequence.zfill(max_length)
    encoded = keras.utils.to_categorical(
        [int(char) for char in sequence], num_classes=NUM_CLASSES
    )
    return encoded


def decode_sequence(encoded_seq):
    """
    Decodes a one-hot encoded sequence back to a digit string.

    Parameters:
        encoded_seq (np.ndarray): One-hot encoded array of shape (max_length, NUM_CLASSES)

    Returns:
        str: Decoded digit sequence as a string.
    """
    digits = np.argmax(encoded_seq, axis=-1)
    # Convert to string and remove leading zeros
    digit_str = "".join(str(digit) for digit in digits).lstrip("0")
    return digit_str if digit_str else "0"


def prepare_input(output_sequence):
    """
    Prepares the output sequence for prediction by encoding it.

    Parameters:
        output_sequence (str): The digit sequence as a string.

    Returns:
        np.ndarray: One-hot encoded array of shape (1, max_length, NUM_CLASSES)
    """
    encoded = encode_sequence(output_sequence)  # Shape: (max_length, NUM_CLASSES)
    # Expand dimensions to match model's expected input shape: (batch_size, max_length, NUM_CLASSES)
    encoded_expanded = np.expand_dims(encoded, axis=0)
    return encoded_expanded


def predict_input_from_output(model, output_sequence):
    """
    Predicts the input sequence given an output sequence using the trained model.

    Parameters:
        model (keras.Model): The trained Keras model.
        output_sequence (str): The output digit sequence as a string.

    Returns:
        str: The predicted input digit sequence as a string.
    """
    # Prepare the input
    prepared_input = prepare_input(output_sequence)

    # Make prediction
    predicted_output = model.predict(prepared_input)

    # The model's output has shape (1, max_length, NUM_CLASSES)
    # Remove the batch dimension
    predicted_output = predicted_output[0]

    # Decode the predicted input sequence
    predicted_input_sequence = decode_sequence(predicted_output)

    return predicted_input_sequence


def main():
    # Load the trained model in the native Keras format
    model = keras.models.load_model("best_model.keras")
    print("Model loaded successfully from 'best_model.keras'.")

    while True:
        # Prompt the user for an output sequence
        output_sequence = input("Enter an output digit sequence (or 'exit' to quit): ")

        if output_sequence.lower() == "exit":
            print("Exiting.")
            break

        # Validate the input
        if not output_sequence.isdigit():
            print("Invalid input. Please enter a sequence of digits.")
            continue

        # Check length and adjust if necessary
        if len(output_sequence) > MAX_SEQUENCE_LENGTH:
            print(
                f"Input sequence is longer than {MAX_SEQUENCE_LENGTH} digits. Truncating."
            )
            output_sequence = output_sequence[-MAX_SEQUENCE_LENGTH:]
        elif len(output_sequence) < MAX_SEQUENCE_LENGTH:
            print(
                f"Padding the input sequence to {MAX_SEQUENCE_LENGTH} digits with leading zeros."
            )

        # Predict the input sequence
        predicted_input = predict_input_from_output(model, output_sequence)

        print(f"Predicted Input Sequence: {predicted_input}\n")


if __name__ == "__main__":
    main()
