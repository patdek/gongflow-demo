ngflow-demo
========================
This demo users angular, ng-flow-standalone, boostrap
and the demo application from the ng-flow github repo: 
https://github.com/flowjs/ng-flow/tree/master/samples/basic

Regenerate static assets
========================
This is only needed if you mess with the static assets 
embedded in the demo, which probably isn't worth it anyway.

prompt: go get https://github.com/jteeuwen/go-bindata
prompt: $GOPATH/bin/go-bindata html/...

The output will be a fresh new bindata.go with the static
assets inside of it -- then you can recompile and enjoy.
