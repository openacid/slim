# convert dot to jpg

DOTs = $(wildcard *.dot)
JPGs = $(DOTs:.dot=.jpg)

all: $(JPGs)

%.jpg:%.dot
	dot -Tjpg  -Gsize=3,10\! -Gdpi=100 $< -o $@

clean:
	rm $(JPGs)
