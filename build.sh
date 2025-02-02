#!/bin/bash
cd "$(dirname "$0")" || exit

source scripts/logoGen.sh

ShowError() {
  code="$1"
  msg="${*:2}"
  printf "\033[1;31m[ ERROR ]:\033[91m $msg\033[0m\n"
  exit "$code"
}

Sep() {
  local w=$(tput cols)
  local sep=$(printf '====----%.0s' $(seq 1 $((w / 8 + 1))))
  sep=${sep:0:w}
  echo -e "$sep\n"
}

if [[ " $* " == *" --noPy "* ]]; then
  cd src > /dev/null 2>&1 || ShowError 1 "Can't find the src directory, you are probably running this script from the wrong place, go to the root directory where you have PixelEconomy and run ./build.sh"
  printf "compilation...\n"
  go build -o "../bin/PixelEconomy.elf" > /dev/null 2>&1 || ShowError 1 "It looks like the compilation failed to complete, try again later"
  cd ..
  printf "\033[32;1mThe compilation was successful!\033[0m\n"
  if [[ " $* " == *" --run "* ]]; then
    printf "Running...\n"
    sleep 1
    chmod +x bin/PixelEconomy.elf
    bin/PixelEconomy.elf
  fi
  exit 0
fi

if ! command -v python3 &>/dev/null; then
  ShowError 1 "it looks like Python is not installed, to continue install it, or run the script with the --noPy flag"
fi

script="scripts/build.py"

if [ ! -f "$script" ]; then
  ShowError 1 "Can't find the src directory, you are probably running this script from the wrong place, go to the root directory where you have PixelEconomy and run ./build.sh"
fi

python3 "$script" "$@"
exit=$?

if [ $exit -ne 0 ]; then
  Sep
  ShowError "$exit" "It looks like the compilation failed to complete, try again later (Exit code: $exit)"
fi

exit $exit
