#!/bin/bash

ShowError() {
    msg="$*"
    printf "\033[1;31m[ ERROR ]:\033[91m $msg\033[0m\n"
    exit 1
}

if ! command -v figlet &>/dev/null; then
  ShowError "Install figlet to run this script"
fi

PIXEL=$(figlet "PIXEL")
ECONOMY=$(figlet -- "ECONOMY")

if [[ " $* " == *" --raw "* ]]; then
  while IFS= read -r PIXEL_LINE && IFS= read -r ECONOMY_LINE <&3; do
    echo "\033[36;1m$PIXEL_LINE\033[0;1m$ECONOMY_LINE"
  done < <(echo "$PIXEL") 3< <(echo "$ECONOMY")
  exit
fi

while IFS= read -r PIXEL_LINE && IFS= read -r ECONOMY_LINE <&3; do
  printf "\033[36;1m%s\033[0;1m%s\n" "$PIXEL_LINE" "$ECONOMY_LINE"
done < <(echo "$PIXEL") 3< <(echo "$ECONOMY")