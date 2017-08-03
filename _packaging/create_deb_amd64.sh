cd /deb/work

debuild --no-tgz-check -uc -us

mv /deb/*_amd64.deb /output
