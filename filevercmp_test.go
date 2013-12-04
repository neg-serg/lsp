package main

import "testing"
import "regexp"

func TestSuffix(t *testing.T) {
	r := regexp.MustCompile(`(\.[A-Za-z~][A-Za-z0-9~]*)+$`)
	for _, s := range examples {
		i := suffixIndex(s)
		is := r.FindStringIndex(s)
		if i == -1 {
			if len(is) != 0 && s[:is[0]] != "" {
				t.Errorf("%s has suffix %s!", s, s[:is[0]])
			}
		} else {
			if len(is) != 0 && is[0] != i {
				t.Errorf("different suffixes!\nRE: %3d %s\nN:  %3d %s\n",
					is[0], s[:is[0]], i, s[:i])

			}
		}
	}
}

var examples = []string{
	"",
	".",
	"..",
	".0",
	".9",
	".A",
	".Z",
	".a~",
	".a",
	".b~",
	".b",
	".z",
	".zz~",
	".zz",
	".zz.~1~",
	".zz.0",
	"0",
	"9",
	"A",
	"Z",
	"a~",
	"a",
	"a.b~",
	"a.b",
	"a.bc~",
	"a.bc",
	"b~",
	"b",
	"gcc-c++-10.fc9.tar.gz",
	"gcc-c++-10.fc9.tar.gz.~1~",
	"gcc-c++-10.fc9.tar.gz.~2~",
	"gcc-c++-10.8.12-0.7rc2.fc9.tar.bz2",
	"gcc-c++-10.8.12-0.7rc2.fc9.tar.bz2.~1~",
	"glibc-2-0.1.beta1.fc10.rpm",
	"glibc-common-5-0.2.beta2.fc9.ebuild",
	"glibc-common-5-0.2b.deb",
	"glibc-common-11b.ebuild",
	"glibc-common-11-0.6rc2.ebuild",
	"libstdc++-0.5.8.11-0.7rc2.fc10.tar.gz",
	"libstdc++-4a.fc8.tar.gz",
	"libstdc++-4.10.4.20040204svn.rpm",
	"libstdc++-devel-3.fc8.ebuild",
	"libstdc++-devel-3a.fc9.tar.gz",
	"libstdc++-devel-8.fc8.deb",
	"libstdc++-devel-8.6.2-0.4b.fc8",
	"nss_ldap-1-0.2b.fc9.tar.bz2",
	"nss_ldap-1-0.6rc2.fc8.tar.gz",
	"nss_ldap-1.0-0.1a.tar.gz",
	"nss_ldap-10beta1.fc8.tar.gz",
	"nss_ldap-10.11.8.6.20040204cvs.fc10.ebuild",
	"z",
	"zz~",
	"zz",
	"zz.~1~",
	"zz.0",
	"#.b#",
}

func TestVCMP(t *testing.T) {
	/* Following tests taken from test-strverscmp.c */
	if !(filevercmp("", "") == 0) {
		t.Fail()
	}
	if !(filevercmp("a", "a") == 0) {
		t.Fail()
	}
	if !(filevercmp("a", "b") < 0) {
		t.Fail()
	}
	if !(filevercmp("b", "a") > 0) {
		t.Fail()
	}
	if !(filevercmp("a0", "a") > 0) {
		t.Fail()
	}
	if !(filevercmp("00", "01") < 0) {
		t.Fail()
	}
	if !(filevercmp("01", "010") < 0) {
		t.Fail()
	}
	if !(filevercmp("9", "10") < 0) {
		t.Fail()
	}
	if !(filevercmp("0a", "0") > 0) {
		t.Fail()
	}
	for i, is := range examples {
		for j, js := range examples {
			result := filevercmp(is, js)
			if result < 0 {
				if !(i < j) {
					t.Logf(`verrevcmp("%s", "%s") = %d`,
						is, js, verrevcmp(is, js))
					t.Log(suffixIndex(is), suffixIndex(js))
					t.Errorf(`"%s", "%s"`, is, js)
				}
			} else if 0 < result {
				if !(j < i) {
					t.Logf(`verrevcmp("%s", "%s") = %d`,
						is, js, verrevcmp(is, js))
					t.Log(suffixIndex(is), suffixIndex(js))
					t.Errorf(`"%s", "%s"`, is, js)
				}
			} else {
				if j != i {
					t.Logf(`verrevcmp("%s", "%s") = %d`,
						is, js, verrevcmp(is, js))
					t.Log(suffixIndex(is), suffixIndex(js))
					t.Errorf(`"%s", "%s"`, is, js)
				}
			}
		}
	}
}
