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

arduino_cli_dir="$(pwd)/arduino-cli/bin/"

cd vr-controlled-robot-arm

echo "Target branch: $target_branch"
echo "Port: $port"

git checkout -f $target_branch
git pull

$arduino_cli_dir/arduino-cli compile --fqbn arduino:avr:uno robot
$arduino_cli_dir/arduino-cli upload -p $port --fqbn arduino:avr:uno robot
