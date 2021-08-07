#!/bin/sh

set -e
set -x

go run . encode example.png "Hello, Stega! สวัสดีสเตกาโนกราฟี"
go run . decode example.png
git checkout example.png
