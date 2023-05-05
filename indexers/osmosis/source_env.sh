#!/bin/bash

while read -r line; do
  eval "$line"
done < /run/secrets/starscope-env