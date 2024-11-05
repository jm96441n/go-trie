# Trie Implementation

This is a small `trie` implementation in go, the dictionary is ~370,000 words and can be found in the `unique-words.txt` list (I used the `words_alpha.txt` list from [this repo](https://github.com/dwyl/english-words)). This spins up a webserver running on port 8080 where you can query for words in the list. An example is:
```
curl "localhost:8080/search?zenith"
```
