#!/bin/bash
/home/$USER/bin/dockerd-rootless-setuptool.sh uninstall -f
/home/$USER/bin/rootlesskit rm -rf /home/$USER/.local/share/docker
rm -rf ~/.config/systemd/user/docker.service
rm -rf ~/.config/docker
rm -rf ~/.docker/run

