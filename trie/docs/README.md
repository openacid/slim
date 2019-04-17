# Illustration of Reducing a Trie to SlimTrie

## Steps

### Bitrie

![](imgs/bitrie.jpg)

### Initialize Trie

![](imgs/slim-init.jpg)


### Reduce Leaf Nodes

![](imgs/slim-cut-leaf.jpg)


### Reduce Inner Single Branch

![](imgs/slim-cut-inner.jpg)


### Remove Skip from Leaf Nodes

![](imgs/slim-final.jpg)


### Remove Leaf Nodes

![](imgs/slim-no-leaf.jpg)

## Update

```
cd imgs
```

Edit `.dot` files

```
make clean
make
```


## Dependency

-   graphviz
