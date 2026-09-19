package passport

import (
	"context"
	"testing"
)

type fakeRepository struct{}

func (fakeRepository) Get(context.Context, string) (View, error) {
	return View{CandidateID: "candidate", AnonymousID: "KZ-1042", OverallScore: 82}, nil
}
func TestServiceGetsCandidatePassport(t *testing.T) {
	view, err := NewService(fakeRepository{}).Get(context.Background(), "candidate")
	if err != nil {
		t.Fatal(err)
	}
	if view.AnonymousID != "KZ-1042" {
		t.Fatalf("unexpected passport %#v", view)
	}
}
