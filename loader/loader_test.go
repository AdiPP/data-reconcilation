package loader

import "testing"

func TestSuccessLoadFile(t *testing.T) {
	_, err := NewLoader().Load("../source_test/valid_proxy_test.csv")

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestFailedLoadFile(t *testing.T) {
	_, err := NewLoader().Load("../really_random_path.csv")
	expectation := FILE_NOT_FOUND

	if err.Error() != expectation {
		t.Errorf("Expected erorr %v but got %v", expectation, err.Error())
	}
}
