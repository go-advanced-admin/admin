package fields

import "testing"

func TestEmailField_emailValidation(t *testing.T) {
	f := &EmailField{}
	validEmails := []string{"test@example.com", "user.name+tag+sorting@example.com"}
	for _, email := range validEmails {
		errs, err := f.emailValidation(email)
		if err != nil || len(errs) != 0 {
			t.Errorf("Expected email '%s' to be valid, got errs: %v, backend error: %v", email, errs, err)
		}
	}

	invalidEmails := []string{"plainaddress", "@missingusername.com", "username@.com"}
	for _, email := range invalidEmails {
		errs, err := f.emailValidation(email)
		if err != nil {
			t.Errorf("Unexpected backend error: %v", err)
		}
		if len(errs) == 0 {
			t.Errorf("Expected email '%s' to be invalid", email)
		}
	}
}
