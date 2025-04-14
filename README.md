# lettertoe

### Rules

In an NxN grid (where N can be 3, 4, 5), the goal is to make as many words as possible during the game. It's a 2 player game. Each player has the right to put a single letter onto the board. For each word a player makes, they get a point. 

We check in these directions:
- VERTICALLY 2*N columns
- HORIZONTALLY 2*N lines
- DIAGONALLY 4 lines

The numbers are doubled because both directions of a column, line or a diagonal line can be used to construct words. This means that if the N lettered character sequence is a word in both ways it'll count twice (e.g, live - evil).

The game ends when all the board slots have been filled with letters.