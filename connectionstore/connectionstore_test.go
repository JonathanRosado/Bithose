package connectionstore

import "testing"

func TestConnection_WithOneMatchingLabelPairAndCriterion(t *testing.T) {
	labelName := "channel"
	labelValue := "viva_la_vida"

	ch := make(chan []byte)
	connection := Connection{
		Ch: ch,
		LabelAcceptanceCriteria: []LabelAcceptanceCriterion{
			{
				LabelPair: LabelPair{
					Name:  labelName,
					Value: labelValue,
				},
				Operator: "==",
			},
		},
	}

	result, err := connection.AcceptsLabels([]LabelPair{
		{
			Name:  labelName,
			Value: labelValue,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if !result {
		t.Error("Should return true as label pair should match criteria")
	}
}

func TestConnection_WithTwoCriterionsAndOneLabelPair(t *testing.T) {
	labelName := "channel"
	labelValue := "viva_la_vida"
	labelName2 := "event"
	labelValue2 := "button_press"

	ch := make(chan []byte)
	connection := Connection{
		Ch: ch,
		LabelAcceptanceCriteria: []LabelAcceptanceCriterion{
			{
				LabelPair: LabelPair{
					Name:  labelName,
					Value: labelValue,
				},
				Operator: "==",
			},
			{
				LabelPair: LabelPair{
					Name:  labelName2,
					Value: labelValue2,
				},
				Operator: "==",
			},
		},
	}

	result, err := connection.AcceptsLabels([]LabelPair{
		{
			Name:  labelName,
			Value: labelValue,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if result {
		t.Error("Should return false as there is one criterion which does not have a corresponding label pair")
	}
}

func TestConnection_WithTwoCriterionsAndTwoLabelPairs(t *testing.T) {
	labelName := "channel"
	labelValue := "viva_la_vida"
	labelName2 := "event"
	labelValue2 := "button_press"

	ch := make(chan []byte)
	connection := Connection{
		Ch: ch,
		LabelAcceptanceCriteria: []LabelAcceptanceCriterion{
			{
				LabelPair: LabelPair{
					Name:  labelName,
					Value: labelValue,
				},
				Operator: "==",
			},
			{
				LabelPair: LabelPair{
					Name:  labelName2,
					Value: labelValue2,
				},
				Operator: "==",
			},
		},
	}

	result, err := connection.AcceptsLabels([]LabelPair{
		{
			Name:  labelName,
			Value: labelValue,
		},
		{
			Name:  labelName2,
			Value: labelValue2,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if !result {
		t.Error("Should return false as there is one criterion which does not have a corresponding label pair")
	}
}

func TestConnection_WithOneLabelDoesNotMatch(t *testing.T) {
	labelName := "channel"
	labelValue := "viva_la_vida"
	labelName2 := "event"
	labelValue2 := "button_press"
	labelValue3 := "key_press"

	ch := make(chan []byte)
	connection := Connection{
		Ch: ch,
		LabelAcceptanceCriteria: []LabelAcceptanceCriterion{
			{
				LabelPair: LabelPair{
					Name:  labelName,
					Value: labelValue,
				},
				Operator: "==",
			},
			{
				LabelPair: LabelPair{
					Name:  labelName2,
					Value: labelValue2,
				},
				Operator: "==",
			},
		},
	}

	result, err := connection.AcceptsLabels([]LabelPair{
		{
			Name:  labelName,
			Value: labelValue,
		},
		{
			Name:  labelName2,
			Value: labelValue3, // Should cause to not match
		},
	})

	if err != nil {
		t.Error(err)
	}

	if result {
		t.Error("Should return false as the second label did not match the second criterion")
	}
}

func TestConnection_WithOneCriterionAndManyLabels(t *testing.T) {
	labelName := "channel"
	labelValue := "viva_la_vida"
	labelName2 := "event"
	labelValue2 := "button_press"
	labelName3 := "event_2"
	labelValue3 := "key_press"

	ch := make(chan []byte)
	connection := Connection{
		Ch: ch,
		LabelAcceptanceCriteria: []LabelAcceptanceCriterion{
			{
				LabelPair: LabelPair{
					Name:  labelName,
					Value: labelValue,
				},
				Operator: "==",
			},
		},
	}

	result, err := connection.AcceptsLabels([]LabelPair{
		{
			Name:  labelName,
			Value: labelValue,
		},
		{
			Name:  labelName2,
			Value: labelValue2, // Should cause to not match
		},
		{
			Name:  labelName3,
			Value: labelValue3,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if !result {
		t.Error("Should return true as the only criterion specified in matched by a label")
	}
}

func TestConnection_WithGreaterThanOperator(t *testing.T) {
	labelName := "channel"
	labelValue := 5

	ch := make(chan []byte)
	connection := Connection{
		Ch: ch,
		LabelAcceptanceCriteria: []LabelAcceptanceCriterion{
			{
				LabelPair: LabelPair{
					Name:  labelName,
					Value: labelValue,
				},
				Operator: ">",
			},
		},
	}

	{
		newLabelValue := 6
		result, err := connection.AcceptsLabels([]LabelPair{
			{
				Name:  labelName,
				Value: newLabelValue,
			},
		})

		if err != nil {
			t.Error(err)
		}

		if !result {
			t.Error("Should return true as the label value is greater than the specified criterion")
		}
	}

	{
		newLabelValue := 2
		result, err := connection.AcceptsLabels([]LabelPair{
			{
				Name:  labelName,
				Value: newLabelValue,
			},
		})

		if err != nil {
			t.Error(err)
		}

		if result {
			t.Error("Should return false as the label value is less than the specified criterion")
		}
	}
}
