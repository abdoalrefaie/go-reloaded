Go Reloaded: Text Formatter
Go Reloaded is a command-line utility written in Go that reads a text file, applies a series of string manipulations and grammatical formatting rules, and writes the corrected text to an output file.

🚀 Features
This tool automatically processes and cleans up text based on specific inline modifiers and general grammatical rules:

Number Conversions:

(hex): Converts the preceding hexadecimal word into its decimal equivalent.

(bin): Converts the preceding binary word into its decimal equivalent.

Text Casing Modifiers:

(up): Converts the preceding word to UPPERCASE.

(low): Converts the preceding word to lowercase.

(cap): Capitalizes the first letter of the preceding word.

Multi-word variants: You can specify how many preceding words to modify by adding a number (e.g., (up, 2), (low, 3), (cap, 4)).

Punctuation Formatting: * Automatically removes spaces before commas, periods, exclamation points, question marks, colons, and semicolons (.,!?:;), and ensures there is exactly one space after them.

Quote Formatting: * Fixes spacing around single (') and double (") quotation marks so they snugly enclose the text (e.g., ' hello ' becomes 'hello').

Article Correction: * Automatically changes a to an (or A to An) if the following word begins with a vowel or an 'h' (a, e, i, o, u, h).

🛠️ Usage
To run the program, you need to provide an input file (containing the text to be formatted) and an output file (where the formatted text will be saved).

Bash
go run . <input_file> <output_file>
Note: If the output file does not exist, the program will create it. If it does exist, it will be overwritten.

Example
input.txt

Plaintext
It was the best of times, it was the worst of times (up) .
1E (hex) files were added to the directory.
There is a hour left before the show starts.
Please type ' hello world ' on your screen.
I am feeling very happy (cap, 2) today!
Command:

Bash
go run . input.txt output.txt
output.txt

Plaintext
It was the best of times, it was the worst of TIMES.
30 files were added to the directory.
There is an hour left before the show starts.
Please type 'hello world' on your screen.
I am feeling Very Happy today!
🏗️ How it Works
The program executes its transformations sequentially to ensure all rules are properly applied:

Parses Modifiers: It scans for tags like (hex) or (up, 2) and alters the preceding words accordingly.

Fixes Punctuation: Realigns standard punctuation marks to their proper spacing.

Fixes Quotes: Pairs up quotation marks and strips out extraneous padding spaces.

Corrects Articles: Scans for a and an to ensure they match the phonetic start of the next word.