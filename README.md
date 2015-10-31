# zodiac
This is a weekend/evening project to try and crack the unzolved Zodiac Killer ciphers, mainly the one commonly referred to as the 340 Cipher. Even if that doesn't happen, I at least want to come up with the logic to programmatically solve the already solved 408 Cipher  It is built using Go and the Revel framework. The plan is to create a list of randomly generated cipher keys, and sort them by how many unique words were found in each solution. I hope to revise and apply the logic I am building to the 408 cipher until it returns what has been the accepted solution, and see what comes out when I apply that to the unsolved 340 cipher. 

**This project makes the following assumptions:**
* All symbols only translate to one letter. 
* Multiple symbols can represent the same letter. (i.e. commonly used letters)
* The ciphers contain errors, either intentional or accidental. 
* The method used to encipher both of codes messages is homophonic simple substitution. 

**Completed:**
* Browser view of ciphers, with key and symbol to letter count
* Randomly generates cipher keys, accounting for a max number of symbols for each letter, and if a letter should be reused or not.
* Indexing into ElasticSearch

**Coming Up:**
* Find and count number of words in each solution generated by random key
* Index with the word count and weights related to unique words, or run a job to back fill
* Solve the 340 Cipher
* Profit? 

**Unsolved 340 Cipher:**
![340-screenshot](340-screenshot.png)

**Solved 408 Cipher:**
![408-screenshot](408-screenshot.png)
