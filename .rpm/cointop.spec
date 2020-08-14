%define version     1.5.3
%global commit      da7c975a19b71cb0c62afd69565ce98eddbb54d3b875e277e0fefe32456b106e
%global shortcommit %(c=%{commit}; echo ${c:0:7})

Name:           cointop
Version:        %{version}
Release:        6%{?dist}
Summary:        Interactive terminal based UI application for tracking cryptocurrencies
License:        Apache-2.0
URL:            https://cointop.sh
Source0:        https://github.com/miguelmota/%{cointop}/archive/v%{version}.tar.gz

BuildRequires:  gcc
BuildRequires:  golang >= 1.14

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
