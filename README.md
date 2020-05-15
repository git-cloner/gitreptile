# gitreptile
get some names of repository which stars is high

First ,you need clone this repository from https://github.com/git-cloner/gitreptile.git

```
git clone https://github.com/git-cloner/gitreptile.git
```

Second,  cd reptile, create file `file.txt`,then start webserver:

```
cd reptile

go run main.go
```

Last, cd request and start reptile
```
cd request
go run main.go -start="10000" -end="400000"
```
you know that: if you want to know the top stars of github, you need access website like `https://github.com/search?q=stars%3A10000..400000`,10000-40000 is the range of the stars of repository belong to,so you can set up params of start and end to limit the range

if you checkout some stars of repository which in range 1000 - 2000 ，you could run commond like `go run main.go -start="1000" -end="2000"`,then you could find file.txt which incloud the names of repository after 100000 second。