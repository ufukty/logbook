package columns

import (
	"fmt"
	"testing"
)

func TestGroupNames(t *testing.T) {
	tcs := []GroupName{
		"Engineering",
		"Finance",
		"Marketing",
		"HR",
		"Operations",
		"Sales",
		"Legal",
		"Customer Support",
		"Product Development",
		"Research and Development",
		"Strategy",
		"Design",
		"Logistics",
		"Procurement",
		"Compliance",
		"Data Analytics",
		"Innovation",
		"IT",
		"Quality Assurance",
		"Security",
		"Growth",
		"Custodial Services",
		"Accounts Payable",
		"Creative Services",
		"Internal Audit",
		"Project Planning",
		"Architectural Design",
		"Construction",
		"Data Science",
		"Automation Team",
		"Coding Ninjas",
		"Rockstar Developers",
		"Wizard Consultants",
		"Growth Gurus",
		"Bean Counters",
		"Pencil Pushers",
		"Code Monkeys",
		"Number Crunchers",
		"The Fixers",
		"Task Force",
		"Bug Squashers",
		"Idea Factory",
		"The Think Tank",
		"Paper Shufflers",
		"Button Pushers",
		"Pixel Pushers",
		"Firefighters",
		"Deadline Dodgers",
		"Chaos Coordinators",
		"The Dream Team",
	}
	for _, tc := range tcs {
		if err := tc.Validate(); err != nil {
			t.Error(fmt.Errorf("got error for valid value %q: %v", tc, err))
		}
	}
}
