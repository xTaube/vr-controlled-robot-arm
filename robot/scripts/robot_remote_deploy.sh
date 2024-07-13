#!/bin/bash
function ArgHelp {
    echo "You need to specify argument [-b] [-p]"
    echo "where:"
    echo "  -b: Target branch"
    echo "  -p: Port"
    echo ""

    exit 1
}

# Exit if arguments are not passed
if [[ $# != 4 ]]; then
    ArgHelp
fi

while getopts b:p: flag
do
    case "${flag}" in
        b) target_branch=${OPTARG};;
        p) port=${OPTARG};;
        *) ArgHelp;;
    esac
done

ARDUINO_CLI_BIN=$PWD/local/bin

cd $PWD/v-arm/vr-controlled-robot-arm

git checkout -f $target_branch > /dev/null
git pull > /dev/null

$ARDUINO_CLI_BIN/arduino-cli compile --fqbn arduino:avr:uno robot
$ARDUINO_CLI_BIN/arduino-cli upload -p $port --fqbn arduino:avr:uno robot
