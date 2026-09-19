## Writing an interpreter in Go

I have always been curious about compilers and how they work, so i decided to build little one myself

### Features

1.  Scanner: Takes a source string in spits tokens out, simple stuff.
2.  Parser: I added parser package and made the parser usable with different strategies, right now only recursive descent strategy is implemented and it parses just expressions
