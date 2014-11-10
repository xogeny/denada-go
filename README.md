# Denada - A declarative language for creating simple DSLs

This is an implementation of Denada in Go (golang).  This is an
attempt to recreate the functionality of my previous
[Javascript implementation of Denada](https://github.com/xogeny/denada-js)
in the Go language.

## TL;DR

Denada allows you to very quickly create a DSL for expressing specific
data and/or structure.  It has the advantage (in my opinion) over XML
and/or JSON that it allows you to formulate a DSL that is human
readable and provide better diagnostic error messages.  Defining a
grammar for your DSL is super easy.

## Background

Denada is based on a project I once worked on where we needed to build
and quickly evolve a simple domain-specific language (DSL).  But an
important aspect of the project was that there were several
non-developers involved.  We developed a language very similar to this
one that had several interesting properties (which I'll come to
shortly).  Recently, I was faced with a situation where I needed to
develop a DSL for a project and decided to follow the same approach.

I can already imagine people rolling their eyes at the premise.  But
give me five more minutes and I'll explain why I did it.

There are lots of different ways to build DSLs.  Let's take a quick
walk across the spectrum of possibilities.  One approach to DSL design
is to create an "internal DSL" where you simply use clever syntactic
tricks in some host language to create what looks like a domain
specific language but is really just a set of domain specific
primitives layered on top of an existing language.  Scala is
particularly good for things like this (using `implicit` constructs)
but you can do it in a number of languages (yes, Lisp works well for
this too...but I don't care for the aesthetics of homoiconic
representations).  The problem here is that you expose the user of the
language to the (potentiall complicated) semantics of the host
language.  Depending on the use case, leveraging the host language's
semantics could be a win (you need to implement similar semantics in
your language) or a loss (you add a bunch of complexity and sharp
edges to a language for non-experts).

Another approach is to create a so-called "external DSL".  For this,
you might using a parser generator (e.g. ANTLR) to create a parser for
your language.  This allows you to completely define your semantics
(without exposing people to the host language semantics).  This allows
you to control the complexity of the language.  But you've still got
to create the parser, debug parsing issues, generate a tree, and then
write code to walk the tree.  So this can be a significant investment
of time.  Sure, parser generators can really speed things up.  But
there are cases where some of this work can be skipped.

Another route you can go is to just use various markup languages to
try and represent your data.  Representations like XML, YAML, JSON or
even INI files can be used in this way.  But for some cases this is
either overly verbose, too "technical" or unreadable.

## So Why Denada?

The general philosophy of Denada is to define a syntax *a priori*.  As
a result, you don't need to write a parser for it.  You don't get a
choice about how the language looks.  Sure, it's fun to "design"
languages.  But there are a wide range of simple DSLs that can be
implemented naturally within the limited syntax of Denada.

So, I can already hear people saying "But if you've decided on the
syntax, you've already designed **a** language, what's all this talk
about designing DSLs".  Although the syntax of Denada is defined,
**the grammar isn't**.  Denada allows us to impose a grammar on top of
the syntax in the same way that "XML Schemas" allow us to impose a
structure on top of XML.  And, like "XML Schemas" we use the same
syntax for the grammar as for the language.  But unlike anything XML
related, Denada looks (kind of) like a computer language designed for
humans to read.  It also includes a richer data model.

## An Example

To demonstrate how Denada works, let's work through a simple example.
Imagine I'm a system administator and I need a file format to list all
the assets in the company.

The generic syntax of Denada is simple.  There are two things you can
express in Denada.  One looks like a variable declaration and the
other expresses nested structure (which can contain instances of these
same two things).  For example, this is valid Denada code:

```
printer ABC {
   set location = "By my desk";
   set model = "HP 8860";
}
```

This doesn't really *mean* anything, but it conforms to the required
syntax.  Although this might be useful as is, the real use case for
Denada is defining a grammar that restricts what is permitted.  That's
because this is also completely legal:

```
aardvark ABC {
   set location = "By my desk";
   set order = "Tubulidentata";
}
```

So in this case, we want to restrict ourselves (initially) to
cataloging printers.  To do this, we specify a grammar for our assets
file.  Initially, our grammar could look like this:

```
printer _ "printer*" {
  set location = "$string" "location";
  set model = "$string" "model";
  set networkName = "$string" "name?";
}
```

Note how this looks almost exactly like our original input text?  That
is because **grammars in Denada are Denada files**.  They just have
some special annotations (not syntax!).  In this case, the "name" of
the printer is given as just `_`.  This is a wildcard in Denada means
"any identifier".  Also note the "descriptive string" following the
printer definition, `"printer*"`.  That means that this defines the
`printer` rule and the star indicates we can have zero or more of them
in our file.

Furthermore, this grammar defines the contents of a `printer`
specification (*i.e.*, what information we associated with a printer).
It shows that there can be three lines inside a printer definition.
The first is the `location` of the printer.  This is mandatory because
the rule name, `"location"` has no cardinality specified.  Similarly,
we also have a mandatory `model` property.  Finally, we have an
optional `networkName` property.  We know it is optional because the
rule name `"name?"` ends with a `?`.

By defining the grammar in this way, we specify precisely what can be
included in the Denada file.  But let's not limit ourselves to
printers.  Assume we want to list the computers in the company too.
We could simply create a new rule for computers, *e.g.,*

```
printer _ "printer*" {
  set location = "$string" "location";
  set model = "$string" "model";
  set networkName = "$string" "name?";
}

computer _ "computer*" {
  set location = "$string" "location";
  set model = "$string" "model";
  set networkName = "$string" "name?";
}
```

In this case, the contents of these definitions are the same, so we
could even do this:

```
'printer|computer' _ "asset*" {
  set location = "$string" "location";
  set model = "$string" "model";
  set networkName = "$string" "name?";
}
```

With just this simple grammar, we've created a parser for a DSL that
can parse our sample asset list above and flag errors.

### Named Rules and Recursion

To created recursively nested grammars, it is necessary to somehow
"break" the potentially infinite structure.  Since Denada grammars are
(at least up until now) isomorphic with the input structures, handling
recursion is a challenge.  In fact, any time there are repeated
patterns you'll have an issue with repeating yourself (potentially an
infinite number of times).

For this reason, the "description" field of a **definition** rule can
also include a specification of what **rules** should be matched for
the children.  The plural in "rules" is important.  The specification
for child rules appears in the description after a `>`.  What follows
the `>` can be one of the following:

    * `$root` - Use the rules that appear at the root of the document.
	
	* `$this` - Use the children of the current definition.  This
      is the default so you never have to specify it explicitly (although it
	  will work).

    * `<rulename>` - Where `<rulename>` is the name of a **fully
      qualified** definition rule.  The set of possible rule matches
	  for the children will correspond to the children of all rules
	  that match the specified rulename.  Since multiple rules can
	  have the same name, the search for a match is done across
	  **all** children from all rules.

A simple and convenient shorthand here is to prefix the rule
description with a `^`.  This indicates that the children of the
associated definition should match the siblings of that definition
(*i.e.,* this is a simple way to describe a simple recursive
relationship where the same entities can appear at every level from
this point down).

Note that the Denada cardinality syntax allows you to specify a `-`
for the cardinality.  This means exactly zero occurences of that
entity.  This is useful for creating collections of rules that are
never matched by themselves, but can be referred to in multiple
locations within the grammar for specifying the rules for children.

## Denada Syntax

Here is EBNF for the Denada language:

```
File = { Definition | Declaration } .

Definition = Preface [ string ] "{" File "}" .

Declaration = Preface [ "=" expr ] [ string ] ";" . 

QualifiersAndId = { identifier } identifier .

Modification = identifier "=" expr . 

Modifiers = "(" [ Modification { "," Modification } ] ")" .

Preface = QualifiersAndId [ Modifiers ] .
```

Lexically, we have only a few types of tokens.  Before getting into
what matters, let's point out that whitespace and C++ style comments
are ignored in the grammar (*i.e.*, they can be removed during lexical
analysis).

It is important to point out that because the grammar for a given
Denada DSL is written in Denada, we need a lot of latitude in our
identifiers (the ``identifier`` token in the EBNF grammar).  This is
because they will be used as regular expressions when they are used
within a grammar.  As such, an identifier in Denada is any sequence of
characters that doesn't contain whitespace, a comment or the reserved
characters ``{``, ``}``, ``(``, ``)``, ``/``, ``"``, ``=``, ``;`` or
``,``.

There are really only two other token types in Denada.  The first is
quoted strings (the ``string`` token in the EBNF grammar).  This is
just a sequence of characters that start with a ``"`` and end with an
unescaped ``"``.

Finally, we have the ``expr`` token.  This is a JSON value (*i.e*, not
necessarily an object, but a value).
