# gh_useradd
Takes a list of GitHub usernames and creates unix users of the same name.  Adds authorized_keys file from keys associated with GitHub account.

## Uses
Primarily, this script is used for temporary deployments like containers and cloud resources.  Run it at boot with a list of users and they'll have access immdiately.
