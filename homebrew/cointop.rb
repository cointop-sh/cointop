require "language/go"

class Cointop < Formula
  desc "An interactive terminal based UI application for tracking cryptocurrencies"
  homepage "https://cointop.sh"
  url "https://github.com/miguelmota/cointop/archive/0.0.1.tar.gz"
  sha256 "3b2b039da68c92d597ae4a6a89aab58d9741132efd514bbf5cf1a1a151b16213"
  revision 1
  head "https://github.com/miguelmota/cointop.git"
  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    path = buildpath/"src/github.com/miguelmota/cointop"
    system "go", "get", "-u", "github.com/miguelmota/cointop"
    cd path do
      system "go", "build", "-o", "#{bin}/cointop"
    end
  end

  test do
    system "true"
  end
end
