#Find urls by regular expressions

Common crawl has an index of 5 billion urls
Using this, provide a search for regex searches

e.g. all images on government sites
http://.*gov.*/.*(jpg|jpeg)


## Existing Work

Follow some tricks for regex search based on

- postgres: https://wiki.postgresql.org/images/6/6c/Index_support_for_regular_expression_search.pdf
- Cho & Rajagopalan http://oak.cs.ucla.edu/~cho/papers/cho-regex.pdf
- http://swtch.com/~rsc/regexp/regexp4.html
   [http://web.archive.org/web/20150316192430/http://swtch.com/~rsc/regexp/regexp4.html]

# Approach

- Decompose every url into trigrams (maybe 4-grams would be better???)
- Limit to the first n results for now

# Limitations for simplicity

- start with just domains -- case-insensitive, limited charset
- case sensitivity? domain names are not case sensitive
  [I'm thinking now that case-sensitive for the domains is probably better. the more separation you can get, the better]
- when using 3-grams, you want to order them by the least common first

###

keep sets A,B,C,D in sorted order on disk

for each item Ai in A:
 find first item in B >= Ai
  if equal, append to output and iterate Ai
  if bigger than Ai, find first item in A >= Bi
  lather, rinse, repeat (until we run out of A or B)

keep track of:
 the size of each n-gram set
 the size of the intersection of each pair of n-gram sets
use that info to predict a sort

### is it realistic?

#### Scalability
We can divide the 5 billion urls across multiple servers -- there are no interdependencies, beyond firing out a request to all servers and collating the responses.

We probably want a few machines with lots of RAM, and ideally some local SSD storage. On EC2, that means something in the R series. But for the long run, you can get comparable specs *much* cheaper from hetzner auctions.

#### Storage space

average length of an url is 34 chars. We need as many ngrams as characters (or a couple more if we want to include start/end anchors). That meansapproximately 34 ngrams per url
that means we need to store 34 * 5 billion * (size of each node).

These 170 billion nodes get divided into [charset size] ** 3 sets. For domain names that's alphanumerics (case-insensitive), plus hyphen and dot. So 38**3, or about 55,000
For full paths the charset is a-z, A-Z, 0-9, plus -._~:?#[]@!$&'()*,;=/e. That makes us 82; 82**3 is 550,000.
This higher number is actually better, since it eases search for any one regex. It just means there's a large overhead for each machine


2 options for storage:

a) sorted list of uint32. This makes about 680 GB
b) b+-tree. This will require several times more space, but with the benefit of much faster times for intersection operations. [other variations on binary- or b-trees might also work]

#### Search times

If all URL characters were random, we'd have to go through about 300,000 keys for the first intersection. Some of that will happen in parallel, though. And with luck not all ngrams will be equally common, so we'll be able to limit our set a lot before hitting the big sets.

