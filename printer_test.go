package crawler

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vyeve/crawler/mocks"
)

func TestPrinter_addLinkToSite(t *testing.T) {
	site := make(siteLinks)
	site.addLinkToSite("foo", "bar")
	site.addLinkToSite("foo", "bar")
	if le := len(site); le != 1 {
		t.Fatalf("expected %d, but actual: %d", 1, le)
	}
	val, ok := site["foo"]
	if !ok {
		t.Fatal("unexpected result")
	}
	if le := len(val); le != 1 {
		t.Fatalf("expected %d, but actual: %d", 1, le)
	}
}

func TestPrinterPrint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	site := make(siteLinks)
	site.addLinkToSite("foo", "bar-1")
	site.addLinkToSite("foo", "bar-2")
	site.addLinkToSite("docker", "k8s")

	wr := mocks.NewMockWriter(ctrl)
	first := wr.EXPECT().Write([]byte("docker\n")).Return(len("docker\n"), nil)
	second := wr.EXPECT().Write([]byte("\tk8s\n")).Return(len("\tk8s\n"), nil).After(first)
	third := wr.EXPECT().Write([]byte("foo\n")).Return(len("foo\n"), nil).After(second)
	fourth := wr.EXPECT().Write([]byte("\tbar-1\n")).Return(len("\tbar-1\n"), nil).After(third)
	wr.EXPECT().Write([]byte("\tbar-2\n")).Return(len("\tbar-2\n"), nil).After(fourth)
	err := site.Print(wr)
	if err != nil {
		t.Error(err)
	}
}

func TestPrinterPrintWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	site := make(siteLinks)
	site.addLinkToSite("docker", "k8s")
	testErr := errors.New("test")

	wr := mocks.NewMockWriter(ctrl)
	first := wr.EXPECT().Write([]byte("docker\n")).Return(len("docker\n"), nil)
	wr.EXPECT().Write([]byte("\tk8s\n")).Return(len("\tk8s\n"), testErr).After(first)
	err := site.Print(wr)
	if err != testErr {
		t.Errorf("expected error %v. got: %v", testErr, err)
	}
}

func TestPrinterPrintWithError2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	site := make(siteLinks)
	site.addLinkToSite("docker", "k8s")
	testErr := errors.New("test")

	wr := mocks.NewMockWriter(ctrl)
	wr.EXPECT().Write([]byte("docker\n")).Return(len("docker\n"), testErr)
	err := site.Print(wr)
	if err != testErr {
		t.Errorf("expected error %v. got: %v", testErr, err)
	}
}

func TestSortLinks(t *testing.T) {
	links := map[string]struct{}{
		"a": {},
		"z": {},
		"v": {},
		"b": {},
	}
	expResult := []string{"a", "b", "v", "z"}
	res := sortLinks(links)
	if !reflect.DeepEqual(expResult, res) {
		t.Errorf("expected: %v, actual: %s", expResult, res)
	}
}
