#!/bin/bash

PYTHON_SCRIPT_PATH="./backend/app.py"
VENV_PATH="./backend/venv"

source "$VENV_PATH/bin/activate"

while true
do
    python3 "$PYTHON_SCRIPT_PATH"
    exit_code=$?
    if [ $exit_code -ne 0 ]; then
        echo "Script crashed with exit code $exit_code. Restarting..." >&2
        sleep 1
    fi
done