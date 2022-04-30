package markers

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderSVG(t *testing.T) {

	// Pre-render aircraft concurrently
	for _, v := range Aircraft {

		r := renderSVG{
			scale:        v.scale,
			d:            v.svgPath,
			pathStroked:  true,
			pathFilled:   true,
			bgFilled:     false,
			strokeWidth:  2,
			strokeColour: RGBA{1, 0, 0, 1},
			fillColour:   RGBA{0, 1, 0, 1},
			bgColour:     RGBA{0, 0, 1, 1},
			offsetX:      1,
			offsetY:      1,
		}

		_, _, err := imgFromSVG(r)
		require.NoError(t, err)
	}
}

func TestConsumeCommand(t *testing.T) {
	// define test data
	tables := []struct {
		testDesc             string
		d                    string
		expectedCommandFound bool
		expectedCommand      string
		expectedRemainingD   string
	}{
		{
			testDesc:             "^command ",
			d:                    "m 244.73958,0 c",
			expectedCommandFound: true,
			expectedCommand:      "m",
			expectedRemainingD:   " 244.73958,0 c",
		},
		{
			testDesc:             "^ command ",
			d:                    " m 244.73958,0 c",
			expectedCommandFound: true,
			expectedCommand:      "m",
			expectedRemainingD:   " 244.73958,0 c",
		},
		{
			testDesc:             "^ commandnumber",
			d:                    " m244.73958,0 c",
			expectedCommandFound: true,
			expectedCommand:      "m",
			expectedRemainingD:   "244.73958,0 c",
		},
		{
			testDesc:             "^ command-number",
			d:                    " m-244.73958,0 c",
			expectedCommandFound: true,
			expectedCommand:      "m",
			expectedRemainingD:   "-244.73958,0 c",
		},
		{
			testDesc:             "^ enumber",
			d:                    " e-5-1.55157,-6.20628",
			expectedCommandFound: false,
			expectedCommand:      "",
			expectedRemainingD:   " e-5-1.55157,-6.20628",
		},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("testing '%s'", table.testDesc), func(t *testing.T) {
			commandFound, command, remainingD, err := consumeCommand(table.d)
			require.NoError(t, err)
			t.Run("checking commandFound", func(t *testing.T) {
				assert.Equal(t, table.expectedCommandFound, commandFound)
			})
			if table.expectedCommandFound {
				t.Run("checking command", func(t *testing.T) {
					assert.Equal(t, table.expectedCommand, command)
				})
			}
			t.Run("checking remainingD", func(t *testing.T) {
				assert.Equal(t, table.expectedRemainingD, remainingD)
			})
		})
	}
}

func TestConsumeNumber(t *testing.T) {
	// define test data
	tables := []struct {
		testDesc            string
		d                   string
		expectedNumberFound bool
		expectedNumber      float64
		roundToDigits       int
		expectedRemainingD  string
	}{
		{
			testDesc:            "^command ",
			d:                   "m 244.73958,0 c",
			expectedNumberFound: false,
			expectedRemainingD:  "m 244.73958,0 c",
		},
		{
			testDesc:            "^ command ",
			d:                   " m 244.73958,0 c",
			expectedNumberFound: false,
			expectedRemainingD:  " m 244.73958,0 c",
		},
		{
			testDesc:            "^number,",
			d:                   "244.73958,0 c",
			expectedNumberFound: true,
			expectedNumber:      244.73958,
			roundToDigits:       5,
			expectedRemainingD:  ",0 c",
		},
		{
			testDesc:            "^ number,",
			d:                   " 244.73958,0 c",
			expectedNumberFound: true,
			expectedNumber:      244.73958,
			roundToDigits:       5,
			expectedRemainingD:  ",0 c",
		},
		{
			testDesc:            "^-number,",
			d:                   "-19.45177,2.9398148 -21.49332,76.729166",
			expectedNumberFound: true,
			expectedNumber:      -19.45177,
			roundToDigits:       5,
			expectedRemainingD:  ",2.9398148 -21.49332,76.729166",
		},
		{
			testDesc:            "^,-number,",
			d:                   ",-19.45177,2.9398148 -21.49332,76.729166",
			expectedNumberFound: true,
			expectedNumber:      -19.45177,
			roundToDigits:       5,
			expectedRemainingD:  ",2.9398148 -21.49332,76.729166",
		},
		{
			testDesc:            "^ -number,",
			d:                   " -19.45177,2.9398148 -21.49332,76.729166",
			expectedNumberFound: true,
			expectedNumber:      -19.45177,
			roundToDigits:       5,
			expectedRemainingD:  ",2.9398148 -21.49332,76.729166",
		},
		{
			testDesc:            "^,number ",
			d:                   ",2.9398148 -21.49332,76.729166",
			expectedNumberFound: true,
			expectedNumber:      2.9398148,
			roundToDigits:       5,
			expectedRemainingD:  " -21.49332,76.729166",
		},
		{
			testDesc:            "^,number-",
			d:                   ",2.9398148-21.49332,76.729166",
			expectedNumberFound: true,
			expectedNumber:      2.9398148,
			roundToDigits:       5,
			expectedRemainingD:  "-21.49332,76.729166",
		},
		{
			testDesc:            "^,number-",
			d:                   ",2.9398148-21.49332,76.729166",
			expectedNumberFound: true,
			expectedNumber:      2.9398148,
			roundToDigits:       5,
			expectedRemainingD:  "-21.49332,76.729166",
		},
		{
			testDesc:            "^enumber ",
			d:                   "4e-5 1.55157,-6.20628",
			expectedNumberFound: true,
			expectedNumber:      4e-5,
			roundToDigits:       5,
			expectedRemainingD:  " 1.55157,-6.20628",
		},
		{
			testDesc:            "^enumber,",
			d:                   "4e-5,1.55157,-6.20628",
			expectedNumberFound: true,
			expectedNumber:      4e-5,
			roundToDigits:       5,
			expectedRemainingD:  ",1.55157,-6.20628",
		},
		{
			testDesc:            "^enumber-",
			d:                   "4e-5-1.55157,-6.20628",
			expectedNumberFound: true,
			expectedNumber:      4e-5,
			roundToDigits:       5,
			expectedRemainingD:  "-1.55157,-6.20628",
		},
	}

	for _, table := range tables {
		t.Run(fmt.Sprintf("testing '%s'", table.testDesc), func(t *testing.T) {
			numberFound, number, remainingD, err := consumeNumber(table.d)
			require.NoError(t, err)
			t.Run("checking numberFound", func(t *testing.T) {
				assert.Equal(t, table.expectedNumberFound, numberFound)
			})
			if table.expectedNumberFound {
				t.Run("checking number", func(t *testing.T) {
					assert.Equal(t, math.Round(table.expectedNumber*(float64(table.roundToDigits)*10))/(float64(table.roundToDigits)*10), math.Round(number*(float64(table.roundToDigits)*10))/(float64(table.roundToDigits)*10))
				})
			}
			t.Run("checking remainingD", func(t *testing.T) {
				assert.Equal(t, table.expectedRemainingD, remainingD)
			})
		})
	}
}
