#!/bin/bash
LD_LIBRARY_PATH="$HOME/bin/mipsel-uClibc/lib:$LD_LIBRARY_PATH" PATH="$HOME/bin/mipsel-uClibc/bin:$PATH" GOARCH="mipsle" GOOS="linux" CC="mipsel-linux-gcc" CXX="mipsel-linux-g++" CGO_ENABLED="1" buffalo build --environment development_sml482 --ldflags "-w -s" -o bin/wishmaster
$HOME/bin/mipsel-uClibc/bin/mipsel-linux-strip bin/wishmaster
