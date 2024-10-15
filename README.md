# Details

This is a repo for a tutorial.

A short overview of how our program will work:

- Read the text file and break it down into chunks.
- Create a `goroutine` worker pool with a fixed number of workers.
- Send each chunk to a worker for processing, and receive a map of words to occurrence count.
- Merge the resultant maps into one final result
- Write the results to a text file.
