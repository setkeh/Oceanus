#!/bin/bash

scrot -s '%Y-%m-%d_$wx$h_scrot.png' -e 'curl --location --request POST https://screenshots.local.setkeh.com/image --form file=@$f && mv $f /home/setkeh/storage/screenshots/' | jq -r '.url' | xclip -selection c