#!/bin/sh

python -m venv env
. ./env/bin/activate

pip install -r requirements_bootstrap.txt

echo "Run: source ./env/bin/activate && python -m ligath-apt"

