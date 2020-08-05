%define version     1.5.1
%global commit      8c225f598f8f3e32afaba7ac4f799577281a160088fec5183919a699e5b4c14e
%global shortcommit %(c=%{commit}; echo ${c:0:7})

Name:           cointop
Version:        %{version}
Release:        6%{?dist}
Summary:        Interactive terminal based UI application for tracking cryptocurrencies
License:        Apache-2.0
URL:            https://cointop.sh
Source0:        https://github.com/miguelmota/%{cointop}/archive/v%{version}.tar.gz

BuildRequires:  gcc
BuildRequires:  golang >= 1.13

%description
cointop is a fast and lightweight interactive terminal based UI application for tracking and monitoring cryptocurrency coin stats in real-time.

%prep
%setup -q -n %{name}-%{version}

%build
mkdir -p ./_build/src/github.com/miguelmota
ln -s $(pwd) ./_build/src/github.com/miguelmota/%{name}

export GOPATH=$(pwd)/_build:%{gopath}
GO111MODULE=off go build -ldflags="-linkmode=external -compressdwarf=false -X github.com/miguelmota/cointop/cointop.version=v%{version}" -o x .

%install
install -d %{buildroot}%{_bindir}
install -p -m 0755 ./x %{buildroot}%{_bindir}/%{name}

%files
%defattr(-,root,root,-)
%doc LICENSE README.md
%{_bindir}/%{name}

%changelog
