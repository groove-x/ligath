#!/bin/sh

python3 -m venv env
. ./env/bin/activate

pip3 install -r requirements_bootstrap.txt

echo "Run: source ./env/bin/activate && python3 -m collect_apt"

