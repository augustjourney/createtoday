package payments

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatStatus(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Status string
		Want   string
	}{
		{
			Status: "CONFIRMED",
			Want:   StatusSucceeded,
		},
		{
			Status: "confirmed",
			Want:   StatusSucceeded,
		},
		{
			Status: "conFirmed",
			Want:   StatusSucceeded,
		},
		{
			Status: "Completed",
			Want:   StatusSucceeded,
		},
		{
			Status: "completed",
			Want:   StatusSucceeded,
		},
		{
			Status: "succeeded",
			Want:   StatusSucceeded,
		},
		{
			Status: "success",
			Want:   StatusSucceeded,
		},
		{
			Status: "canceled",
			Want:   StatusCanceled,
		},
		{
			Status: "Canceled",
			Want:   StatusCanceled,
		},
		{
			Status: "order_canceled",
			Want:   StatusCanceled,
		},
		{
			Status: "order_denied",
			Want:   StatusRejected,
		},
		{
			Status: "Declined",
			Want:   StatusRejected,
		},
		{
			Status: "eclined",
			Want:   StatusRejected,
		},
		{
			Status: "REJECTED",
			Want:   StatusRejected,
		},
		{
			Status: "rejected",
			Want:   StatusRejected,
		},
		{
			Status: "DEADLINE_EXPIRED",
			Want:   StatusExpired,
		},
		{
			Status: "",
			Want:   StatusPending,
		},
		{
			Status: ".",
			Want:   StatusPending,
		},
		{
			Status: "sdfsdfsf",
			Want:   StatusPending,
		},
	}

	for _, testCase := range cases {
		name := fmt.Sprintf(`Status %s should be %s`, testCase.Status, testCase.Want)
		t.Run(name, func(t *testing.T) {
			gotStatus := FormatStatus(testCase.Status)
			assert.Equal(t, testCase.Want, gotStatus)
			t.Log(testCase.Status, gotStatus)
		})
	}

}
