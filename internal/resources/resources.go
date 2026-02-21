package resources

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed CoolSigge.JPG
var coolSiggeBytes []byte

// CoolSigge is the embedded CoolSigge.JPG image, available as a Fyne static resource.
var CoolSigge = fyne.NewStaticResource("CoolSigge.JPG", coolSiggeBytes)
