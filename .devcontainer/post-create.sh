#!/bin/bash

#######
# Config ZSH
#######
# powerline fonts for zsh theme
git clone https://github.com/powerline/fonts.git
cd fonts
./install.sh
cd .. && rm -rf fonts

# Set ZSH_THEME to agnoster in ~/.zshrc
if grep -q '^ZSH_THEME=' ~/.zshrc 2>/dev/null; then
    sed -i 's/^ZSH_THEME=.*/ZSH_THEME="agnoster"/' ~/.zshrc
else
    echo 'ZSH_THEME="agnoster"' >> ~/.zshrc
fi

