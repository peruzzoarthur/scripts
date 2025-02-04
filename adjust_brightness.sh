#!/bin/bash

# Get current brightness (VCP 10) with sudo
current_brightness=$(sudo ddcutil getvcp 10 | grep -oP '(?<=Current value: )\d+')

# Check if the current brightness value was found
if [ -z "$current_brightness" ]; then
  echo "Error: Could not retrieve the current brightness value."
  exit 1
fi

# Set increment value
increment=10

# Debug: Output current brightness
echo "Current brightness: $current_brightness"

# Optionally, check if brightness is within acceptable bounds (0-100)
if [[ "$1" == "up" ]]; then
  new_brightness=$((current_brightness + increment))
  echo "Increasing brightness by $increment: $new_brightness"
elif [[ "$1" == "down" ]]; then
  new_brightness=$((current_brightness - increment))
  echo "Decreasing brightness by $increment: $new_brightness"
else
  echo "Usage: $0 up|down"
  exit 1
fi

# Ensure the new brightness is within the 0-100 range
if (( new_brightness > 100 )); then
  new_brightness=100
  echo "Brightness adjusted to max (100)"
elif (( new_brightness < 0 )); then
  new_brightness=0
  echo "Brightness adjusted to min (0)"
fi

# Set the new brightness with sudo
echo "Setting brightness to $new_brightness"
sudo ddcutil setvcp 10 "$new_brightness"

