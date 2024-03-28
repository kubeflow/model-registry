#!/bin/bash

install_yq_if_not_installed() {
    # Check operating system
    os=$(uname -s)
    case $os in
        Linux*)
            if ! command -v yq &>/dev/null; then
                echo "yq is not installed. Installing..."
                # Linux installation using curl
                curl -L https://github.com/mikefarah/yq/releases/download/3.4.1/yq_linux_amd64 -o /usr/local/bin/yq && sudo chmod +x /usr/local/bin/yq
            else
                echo "yq is already installed."
            fi
            ;;
        Darwin*)
            if ! command -v yq &>/dev/null; then
                echo "yq is not installed. Installing..."
                # macOS installation using Homebrew
                brew install yq
            else
                echo "yq is already installed."
            fi
            ;;
        *)
            echo "Unsupported operating system: $os"
            ;;
    esac
}

# Call the function when the script is sourced
install_yq_if_not_installed