# go4

[![travis badge](https://travis-ci.org/camlistore/go4.svg?branch=master)](https://travis-ci.org/camlistore/go4 "Travis CI")

[go4.org](http://go4.org) is a collection of packages for
Go programmers.

They started out living in [Perkeep](https://perkeep.org)'s repo
and elsewhere but they have nothing to do with Perkeep, so we're
moving them here.

## Details

* **single repo**. go4 is a single repo. That means things can be
    changed and rearranged globally atomically with ease and
    confidence.

* **no backwards compatibility**. go4 makes no backwards compatibility
    promises. If you want to use go4, vendor it. And next time you
    update your vendor tree, update to the latest API if things in go4
    changed. The plan is to eventually provide tools to make this
    easier.

* **forward progress** because we have no backwards compatibility,
    it's always okay to change things to make things better. That also
    means the bar for contributions is lower. We don't have to get the
    API 100% correct in the first commit.

* **no Go version policy** go4 packages are usually built and tested
    with the latest Go stable version. However, go4 has no overarching
    version policy; each package can declare its own set of supported
    Go versions.

* **code review** contributions must be code-reviewed. We're trying
    out Gerrithub, to see if we can find a mix of Github Pull Requests
    and Gerrit that works well for many people. We'll see.

* **CLA compliant** contributors must agree to the Google CLA (the
    same as Go itself). This ensures we can move things into Go as
    necessary in the future. It also makes lawyers at various
    companies happy.  The CLA is **not** a copyright *assignment*; you
    retain the copyright on your work. The CLA just says that your
    work is open source and you have permission to open source it. See
    https://golang.org/doc/contribute.html#cla

* **docs, tests, portability** all code should be documented in the
    normal Go style, have tests, and be portable to different
    operating systems and architectures. We'll try to get builders in
    place to help run the tests on different OS/arches. For now we
    have Travis at least.

## Contact

For any question, or communication when a Github issue is not appropriate,
please contact the [Perkeep mailing
list](https://groups.google.com/forum/#!forum/perkeep).

