#!/bin/sh

python3.7 -m venv env
. ./env/bin/activate

pip install -r requirements_bootstrap.txt

echo "Run: source ./env/bin/activate && python3.7 -m apt"

