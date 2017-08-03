rpmbuild --define "_topdir /rpmbuild" \
    --target=i386 \
    -bb /rpmbuild/SPECS/$NAME.spec

mv /rpmbuild/RPMS/i386/$NAME-*.i386.rpm /output

