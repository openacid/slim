>> go test github.com/openacid/slim/array -bench=.

     eltSize    eltCount  loadFactor    Overhead
           1       65536       1.000         +16
           1       32962       0.500         +59
           1       13115       0.200         +95
           1        6792       0.100        +181
           1         341       0.005       +3645
           1          79       0.001      +15734
goos: linux
goarch: amd64
pkg: github.com/openacid/slim/array
BenchmarkMemOverhead/#00-2         	       1	9644019702 ns/op
     eltSize    eltCount  loadFactor    Overhead
           2       65536       1.000         +14
           2       32579       0.500         +30
           2       13182       0.200         +69
           2        6545       0.100         +97
           2         339       0.005       +1835
           2          61       0.001      +10193
BenchmarkMemOverhead/#01-2         	       1	9593877581 ns/op
     eltSize    eltCount  loadFactor    Overhead
           4       65536       1.000         +13
           4       32550       0.500         +28
           4       13037       0.200         +33
           4        6495       0.100         +73
           4         333       0.005        +938
           4          64       0.001       +4855
BenchmarkMemOverhead/#02-2         	       1	10552236712 ns/op
     eltSize    eltCount  loadFactor    Overhead
           8       65536       1.000          +9
           8       32790       0.500         +13
           8       13108       0.200         +28
           8        6598       0.100         +31
           8         334       0.005        +479
           8          66       0.001       +2363
BenchmarkMemOverhead/#03-2         	       1	11368131201 ns/op
PASS
ok  	github.com/openacid/slim/array	41.248s
