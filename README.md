# "cleanansiblebackups"
Small utlility written in Go to cleanup Ansible backup files created over time when using "backup: yes" when copying files.

# History
After searching for an easy way to cleanup Ansible backup files and not finding one I've decided to write my own small utility. The utility specifically (and only) targets the backups created by ansible.builtin.copy when using "backup: yes", resulting in numerous "<filename>.<ansible info>~" files.
There is no easy method in Ansible to cleanup these backup files.

# Usage
The utility takes 3 parameters:
  - The name of the file
  - The minimum number of files that must exist/remain (regardless of "age")
  - The maximum age in days to keep a file
  
 When incorporating as a task into an Ansible script:
  
```
 - name: Install the Ansible cleanup utility
   ansible.builtin.copy:
     src: cleanansiblebackups
     dest: /usr/local/bin/cleanansiblebackups
     mode: '0755'

 - name: Cleanup Ansible backup files
   shell: |
     cleanansiblebackups -name /etc/snmp/snmpd.conf -mincount 5 -age 31
```

# Build
Assuming Go is installed (v1.16+):
```
git clone https://github.com/moudsen/cleanansiblebackups
cd cleanansiblebackups
go get github.com/itroot/keysort
go build cleanansiblebackups.go
```
  
# Issues and/or Development
I accept Pull requests and I will also respond to Issues created in this Github project.
