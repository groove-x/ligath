#!/bin/sh

python3 -m venv env
. ./env/bin/activate

pip install -r requirements_bootstrap.txt

echo "Run: source ./env/bin/activate && python3 -m apt"

