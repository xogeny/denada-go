# Denada in Go

This is an **alpha** implementation of Denada in Go (golang).  This is
an attempt to recreate the functionality of my previous
[Javascript implementation of Denada](https://github.com/xogeny/denada-js)
but in Go (and, more importantly, with static type checking).

##  Tools and Processes

I bootstrapped this effort using the approach outlined
[here](http://noypi-linux.blogspot.com/2014/07/golang-parser-generator-ebnfyacclex.html)
which meant doing:

```
$ go get github.com/cznic/ebnf2y
$ go get github.com/cznic/golex
$ go get github.com/blynn/nex
```

Currently, I still utilize `nex` in the build process.  But I no
longer regenerate the `yacc` input from ebnf (although I keep the ebnf
around for reference, at least of now).  I don't use `golex` any
longer either.
