Unmerge
=======

A simple utility for filling blanks in a table with the last non blank value of a cell in the same column.

## Justification for existence

The default behaviour in popular spread sheet applications is to store the value of a merged cell in the top left cell leaving the other cells in the merged cell blank, these blank cells make it harder to use these tables for searching.


### Example situation

#### You have

|thing	|Type	|1	|
|-------|-------|---|
|	|Colour	|blue	|
|	|Gender	|N/A	|

#### And you want

|thing	|Type	|1	|
|-------|-------|---|
|thing	|Colour	|blue	|
|thing	|Gender	|N/A	|

## Features

- Unmerge can work on a range of adjacent columns, in one pass, if you need non-adjacent columns do multiple passes.
- Unmerge can use regexp recipes for the delimiter of the input table.
- The output delimiter can be configured to any string.

## License

Unmerge is licensed under the GPLv3. A copy of the license can be found in [LICENSE](LICENSE).
