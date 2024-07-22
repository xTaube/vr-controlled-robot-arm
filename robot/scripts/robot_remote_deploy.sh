#!/bin/bash
function ArgHelp {
    echo "You need to specify argument [-p]"
    echo "where:"
    echo "  -p: Port"
    echo ""

    exit 1
}

# Exit if arguments are not passed
if [[ $# != 2 ]]; then
    ArgHelp
fi

while getopts p: flag
do
    case "${flag}" in
        p) port=${OPTARG};;
        *) ArgHelp;;
    esac
done

ARDUINO_CLI_BIN=$PWD/local/bin
ROBOT_BUILD_DIR=$PWD/v-arm/robot-build

$ARDUINO_CLI_BIN/arduino-cli upload -p $port --fqbn arduino:avr:uno --input-dir $ROBOT_BUILD_DIR
