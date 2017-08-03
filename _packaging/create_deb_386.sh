cd /deb/work

debuild --no-tgz-check -uc -us

mv /deb/*_i386.deb /output/
