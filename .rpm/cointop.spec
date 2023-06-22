%define version     1.5.4
%global commit      9d6d11e71c3c754a654b191a6813a41427c717966999ce335f0f155358a2292b
%global shortcommit %(c=%{commit}; echo ${c:0:7})

Name:           cointop
Version:        %{version}
Release:        6%{?dist}
Summary:        Interactive terminal based UI application for tracking cryptocurrencies
License:        Apache-2.0
URL:            https://cointop.sh
Source0:        https://github.com/cointop-sh/%{cointop}/archive/v%{version}.tar.gz

BuildRequires:  gcc
BuildRequires:  golang >= 1.14

%description
cointop is a fast and lightweight interactive terminal based UI application for tracking and monitoring cryptocurrency coin stats in real-time.

%prep
%setup -q -n %{name}-%{version}

%build
mkdir -p ./_build/src/github.com/cointop-sh
ln -s $(pwd) ./_build/src/github.com/cointop-sh/%{name}

export GOPATH=$(pwd)/_build:%{gopath}
GO111MODULE=off go build -ldflags="-linkmode=external -compressdwarf=false -X github.com/cointop-sh/cointop/cointop.version=%{version}" -o x .

%install
install -d %{buildroot}%{_bindir}
install -p -m 0755 ./x %{buildroot}%{_bindir}/%{name}

%files
%defattr(-,root,root,-)
%doc LICENSE README.md
%{_bindir}/%{name}

%changelog
