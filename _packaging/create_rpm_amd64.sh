echo $NAME

rpmbuild --define "_topdir /rpmbuild" \
    --target=x86_64 \
    -bb /rpmbuild/SPECS/$NAME.spec

mv /rpmbuild/RPMS/x86_64/$NAME-*.x86_64.rpm /output
