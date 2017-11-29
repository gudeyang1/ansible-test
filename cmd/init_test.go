package cmd

import "errors"
import "os"
import "testing"

func TestMkdirIfNotExist(t *testing.T) {
	name := "./itsma_service-1.0.0"
	os.Remove(name)
	t.Log("Case #1 create a directory which does not exist")
	err := MkdirIfNotExist(name)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Case #2 create a directory which already exists")
	if err = MkdirIfNotExist(name); err != nil {
		t.Errorf("expected: nil, got: %v", err.Error())
	}

	if err = os.Remove(name); err != nil {
		t.Fatal("Clean data failed.", err)
	}

	t.Log("Case #3 create a directory which has the same name with an existing file")
	if _, err = os.Create(name); err != nil {
		t.Fatal("Create test data failed.", err)
	}

	defer os.Remove(name)

	if err = MkdirIfNotExist(name); err == nil {
		t.Errorf("expected: err, got: nil. Cause file with same name exists.")
	}
}

func TestCreateService(t *testing.T) {
	os.Chdir("../")
	defer os.Chdir("cmd")

	ucases := []struct {
		ItsmaService Service
		Expected     error
	}{
		{Service{"", "1.0.0", "localhost:5000"}, errors.New("Service name can NOT be empty")},
		{Service{"propel", "1.0.0", "localhost:5000"}, nil},
		{Service{"propel", "1.0.0", "localhost:5000"}, errors.New("mkdir itom-propel-1.0.0: file exists")},
	}

	for _, ucase := range ucases {
		err := CreateService(ucase.ItsmaService)
		if !CompareError(err, ucase.Expected) {
			t.Errorf("Case %v failed. expected: %v, got: %v\n", ucase, ucase.Expected, err)
		}
	}

	os.RemoveAll("itom-propel-1.0.0")
}

func CompareError(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1 != nil && err2 != nil {
		return err1.Error() == err2.Error()
	}

	return false
}
