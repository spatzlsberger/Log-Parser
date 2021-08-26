# Log Parser
The log parsing command line tool was built to concurrently search for matching patterns or strings within a large file. It then outputs the matched lines to the console. This was created to parse large files that cannot usually be opened by other text editors. 

## Run Modes
You can either search for a simple string, or use a more complex regex pattern to look through your file. Below are some examples from the command line. These require that you are located in the root directory of this application.

    // search for a specific string within a file
    go run . -path "<path to input file>" -string "<string to search>"

    // search file for lines that match a given pattern
    go run . -path "<path to input file>" -regex "<regex string>"