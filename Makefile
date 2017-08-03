NAME := ecl2mond
VERSION := 1.0.0
MAINTAINER_NAME := NTT Communications corporation
MAINTAINER_EMAIL :=  monitoring@ecl.ntt.com
BUILD_OS_TARGETS := "linux darwin windows"
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'
RELEASE_DIR := releases
RPMBUILD_DIR := _packaging/rpmbuild
DEBBUILD_DIR := _packaging/deb/work

check-variables:
	echo "NAME: ${NAME}"
	echo "VERSION: ${VERSION}"
	echo "REVISION: ${REVISION}"

setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help

test: deps
	go test $$(glide novendor)

deps: setup
	glide install

update: setup
	glide update

lint:
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $?; \
	done

fmt:
	goimports -w $$(glide nv -x)

bin/%: cmd/%/main.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ $<

build: deps
	go build -ldflags "$(LDFLAGS)" -o ${RELEASE_DIR}/${GOOS}_${GOARCH}/${NAME}${SUFFIX} cmd/ecl2mond/main.go

build-all: build-linux-amd64 build-linux-386 build-darwin-amd64 build-windows-amd64 build-windows-386

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-386:
	@$(MAKE) build GOOS=linux GOARCH=386

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-windows-386:
	@$(MAKE) build GOOS=windows GOARCH=386 SUFFIX=.exe

prepare-package-rpm-amd64: build-linux-amd64
	cp -p ${RELEASE_DIR}/linux_amd64/${NAME} ${RPMBUILD_DIR}/BUILD
	@$(MAKE) move-files-for-rpm-package

prepare-package-rpm-386: build-linux-386
	cp -p ${RELEASE_DIR}/linux_386/${NAME} ${RPMBUILD_DIR}/BUILD
	@$(MAKE) move-files-for-rpm-package

move-files-for-rpm-package:
	cp -p for_build.conf ${RPMBUILD_DIR}/SOURCES/${NAME}.conf
	cat ${RPMBUILD_DIR}/for_build.spec | \
		sed -e 's/{{name}}/${NAME}/g' -e 's/{{version}}/${VERSION}/g' > ${RPMBUILD_DIR}/SPECS/${NAME}.spec
	cp -p ${RPMBUILD_DIR}/for_build.initd ${RPMBUILD_DIR}/SOURCES/${NAME}.initd

package-rpm-amd64: prepare-package-rpm-amd64
	rpmbuild --define "_topdir `pwd`/${RPMBUILD_DIR}" \
		--target=x86_64 \
		-bb ${RPMBUILD_DIR}/SPECS/${NAME}.spec

package-rpm-386: prepare-package-rpm-386
	rpmbuild --define "_topdir `pwd`/${RPMBUILD_DIR}" \
		--target=i386 \
		-bb ${RPMBUILD_DIR}/SPECS/${NAME}.spec

package-rpm-amd64-docker: prepare-package-rpm-amd64
	cd _packaging && docker build -t ${NAME}-rpm-x86_64 -f Dockerfile.rpm_amd64 .
	docker run --rm -i \
		--volume ${GOPATH}/src/github.com/nttcom/ecl2mond/_packaging/output/:/output/ \
		--env NAME=${NAME} \
		-t ${NAME}-rpm-x86_64 \
		sh create_rpm_amd64.sh

package-rpm-386-docker: prepare-package-rpm-386
	cd _packaging && docker build -t ${NAME}-rpm-386 -f Dockerfile.rpm_386 .
	docker run --rm -i \
		--volume ${GOPATH}/src/github.com/nttcom/ecl2mond/_packaging/output/:/output/ \
		--env NAME=${NAME} \
		-t ${NAME}-rpm-386 \
		sh create_rpm_386.sh

prepare-package-deb-amd64: build-linux-amd64
	cp -p ${RELEASE_DIR}/linux_amd64/${NAME} ${DEBBUILD_DIR}/debian/${NAME}.bin
	@$(MAKE) move-files-for-deb-package

prepare-package-deb-386: build-linux-386
	cp -p ${RELEASE_DIR}/linux_386/${NAME} ${DEBBUILD_DIR}/debian/${NAME}.bin
	@$(MAKE) move-files-for-deb-package

move-files-for-deb-package:
	cp -p for_build.conf ${DEBBUILD_DIR}/debian/${NAME}.conf
	cp -p _packaging/deb/for_deb_build.initd ${DEBBUILD_DIR}/debian/${NAME}.init
	cp -p _packaging/deb/dummy.tar.gz _packaging/deb/ecl2mond_${VERSION}.orig.tar.gz
	cat _packaging/deb/rules | \
		sed -e "s/{{name}}/${NAME}/g" > ${DEBBUILD_DIR}/debian/rules
	cat _packaging/deb/control | \
		sed -e "s/{{name}}/${NAME}/g" -e "s/{{version}}/${VERSION}/g" -e "s/{{maintainer_name}}/${MAINTAINER_NAME}/g" -e "s/{{maintainer_email}}/${MAINTAINER_EMAIL}}/g" > ${DEBBUILD_DIR}/debian/control
	cat _packaging/deb/include-binaries | \
		sed -e "s/{{name}}/${NAME}/g" > ${DEBBUILD_DIR}/debian/source/include-binaries
	cat _packaging/deb/dirs | \
		sed -e "s/{{name}}/${NAME}/g" > ${DEBBUILD_DIR}/debian/dirs

package-deb-amd64: prepare-package-deb-amd64
	cd ${DEBBUILD_DIR} && debuild --no-tgz-check -uc -us

package-deb-386: prepare-package-deb-386
	cd ${DEBBUILD_DIR} && debuild --no-tgz-check -uc -us

package-deb-amd64-docker: prepare-package-deb-amd64
	cd _packaging && docker build -t ${NAME}-deb-x86_64 -f Dockerfile.deb_amd64 .
	docker run --rm -i \
		--volume ${GOPATH}/src/github.com/nttcom/ecl2mond/_packaging/output/:/output/ \
		-t ${NAME}-deb-x86_64 \
		sh create_deb_amd64.sh

package-deb-386-docker: prepare-package-deb-386
	cd _packaging && docker build -t ${NAME}-deb-386 -f Dockerfile.deb_386 .
	docker run --rm -i \
		--volume ${GOPATH}/src/github.com/nttcom/ecl2mond/_packaging/output/:/output/ \
		-t ${NAME}-deb-386 \
		sh create_deb_386.sh

package-all-docker: package-rpm-amd64-docker package-rpm-386-docker package-deb-amd64-docker package-deb-386-docker

help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help
