# gorunremote
Run code remotely with profiling and timing.

# environment setup

- Create a user to run the submitted code as, ex: 'codeexec'
- Create a group to allow the user running the server to run code as the user created above ex: sudocodeexec
- Create a directory to store temp code files before running ex: /home/codeexec/codedir
- Lock down the codeexec user's access to the network etc.```$ sudo iptables -I OUTPUT -m owner --uid-owner codeexec -j DROP and $ sudo ip6tables -I OUTPUT -m owner --uid-owner codeexec -j DROP```
- Set group as group of and permissions to 2775 for the codedir
- Add a line to sudoers file to allow the user running the server to sudo as the code running user```%sudocodeexec   ALL=(codeexec) NOPASSWD: ALL```