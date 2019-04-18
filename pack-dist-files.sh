#!/bin/bash

cd lambda-bin

for i in *; do zip -r "../deploy/${i%/}.zip" "$i"; done
