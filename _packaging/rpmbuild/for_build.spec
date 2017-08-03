%define name {{name}}
%define version {{version}}
%define unmanaged_version {{version}}
%define release 1

Name: %{name}
Version:%{version}
Release: %{release}
Summary: Monitoring Agant Tool

Group: Applications/File
License: GPL2
Source0: %{name}-%{unmanaged_version}.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-buildroot
Prefix: %{_prefix}

%description
Monitoring Agent

%prep

%build

%install
%{__rm} -rf %{buildroot}
%{__install} -Dp -m0755 %{_builddir}/%{name} %{buildroot}%{_bindir}/%{name}
%{__install} -m0777 -d %{buildroot}%{_sharedstatedir}/%{name}
%{__install} -Dp -m0644 %{_sourcedir}/%{name}.conf %{buildroot}%{_sysconfdir}/%{name}/%{name}.conf
%{__install} -Dp -m0755 %{_sourcedir}/%{name}.initd %{buildroot}%{_initddir}/%{name}

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-,root,root)
%{_bindir}/%{name}
%{_initddir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/%{name}.conf
%dir %{_sharedstatedir}/%{name}

%preun
if [ "$1" = 0 ]; then
    /sbin/service ecl2mond stop > /dev/null 2>&1
    /sbin/chkconfig --del ecl2mond
fi
exit 0
