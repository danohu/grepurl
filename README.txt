Find urls by regular expressions

Common crawl has an index of 5 billion urls
Using this, provide a search for regex searches

e.g. all images on government sites
http://.*gov.*/.*(jpg|jpeg)


# Existing Work

Follow some tricks for regex search based on

- postgres: https://wiki.postgresql.org/images/6/6c/Index_support_for_regular_expression_search.pdf
- Cho & Rajagopalan http://oak.cs.ucla.edu/~cho/papers/cho-regex.pdf
- http://swtch.com/~rsc/regexp/regexp4.html
   [http://web.archive.org/web/20150316192430/http://swtch.com/~rsc/regexp/regexp4.html]

# Approach

- Decompose every url into trigrams (maybe 4-grams would be better???)
- Limit to the first n results

# Limitations for simplicity

- start with just domains -- case-insensitive, limited charset
- case sensitivity? domain names are not case sensitive
  [I'm thinking now that case-sensitive for the domains is probably better. the more separation you can get, the better]
- when using 3-grams, you want to order them by the least common first

- think bloom filters. not sure quite how, though??
[like -- have a bloom filter tree that would]

don't I want to store more things -- like every 4-gram
400,000 bits
if i can keep it to lowercase, I'm down to

--
keep sets A,B,C,D in sorted order on disk

for each item in A:
 B successor of Ai
  bigger tha
 Ai compare Bi?:
  equa
   hold iterators, test next item
  
 in B?
  no->
   advance

2**32 ~ 5 billion
 just a bit too many

the benefit of hashing is to limit ourselves to the most efficient items

all this is going to be using memory mapped files

so you probably want to allow 2**35 or so, and not worry about alignment
 that means each decomposition is about 20 gigs
  except that you might be able to compress


keep track of:
 the size of each n-gram set
 the size of the intersection of each pair of n-gram sets
use that info to predict a sort


start at first item

if its in the se

# Advanced exercises
- cluster sites by (guess at) which CMS they are running

average length of an url is 34 chars. so we need approximately 34 ngrams per url
that means we need to

--
server setup:
 EC2 provides high-memory R3 instances -- 15gb for $0.175/hour
 basicallly, that whole series is g

http://repository.openoil.net/wiki/MediaWiki:Common.js
