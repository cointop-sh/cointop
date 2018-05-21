%global commit      cf0ad3b8a1cd8bc11d2eb59c3790785d6509fdf8
%global shortcommit %(c=%{commit}; echo ${c:0:7})

Name:           cointop
Version:        1.0.5
Release:        6%{?dist}
Summary:        Terminal based application for tracking cryptocurrencies
License:        Apache-2.0
URL:            https://cointop.sh
Source0:        https://github.com/miguelmota/%{cointop}/archive/%{version}.tar.gz

BuildRequires:  gcc
BuildRequires:  golang >= 1.9-10

%description
cointop is a fast and lightweight interactive terminal based UI application for tracking and monitoring cryptocurrency coin stats in real-time.

%prep
%setup -q -n %{name}-%{version}

%build
mkdir -p ./_build/src/github.com/miguelmota
ln -s $(pwd) ./_build/src/github.com/miguelmota/%{name}

export GOPATH=$(pwd)/_build:%{gopath}
go build -ldflags=-linkmode=external -o x .

%install
install -d %{buildroot}%{_bindir}
install -p -m 0755 ./x %{buildroot}%{_bindir}/%{name}

%files
%defattr(-,root,root,-)
%doc LICENSE.md README.md
%{_bindir}/%{name}

%changelog
