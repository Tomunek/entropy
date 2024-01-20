# Entropy
This program calculates entropy of a file byte by byte. The lower the entropy, the more information the file contains, e.g. text file has low entropy, while compressed / encrypted file has high entropy. Entropy varies between 0 (least entropy - file contains only 1 value) and 8 (file contains all possible byte values, with equal probabilities)


## Usage
```console
entropy [flags] file
  -b int
    	Number of bins on graph (default 10)
  -g	Show output graph 
  -l int
    	Max graph bar length (default 30)
  -q	Number-only output (quiet mode)

``` 

## Examples
Get entropy of `file.txt`:
```console
entropy file.txt
```
Get entropy of file.txt in quiet mode (outputs only float value, for easier automated processing):
```console
entropy -q file.txt
```
Graph a histogram of `file.txt` with default settings:
```console
entropy -g file.txt
```
Graph a histogram of `file.txt` with 20 bins:
```console
entropy -g -b 20 file.txt
```
Graph a histogram of `file.txt` with 20 bins and max bar length of 60 characters:
```console
entropy -g -b 20 -l 60 file.txt
```