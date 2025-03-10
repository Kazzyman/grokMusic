package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math"
	"math/rand"
)

// NotePosition represents a note's position on the staff
type NotePosition struct {
	Pitch string  // such as "A5", "G5", "F5", or "E5"
	Y     float32 // Y-coordinate on the canvas for this note
}

// MarkedNote tracks a placed note for retraction
type MarkedNote struct {
	Circle *canvas.Circle
	X      float32
	Y      float32
}

func main() { // ctrl-M to navigate to matching brace. main is some 300 lines long!
	// Initialize Fyne app
	myApp := app.New()                           // app.___ is a Fyne object.
	myWindow := myApp.NewWindow("Find the Note") // create the app window and title it.
	myWindow.Resize(fyne.NewSize(1065, 980))     // I made these a bit smaller. Resize is in Fyne package.

	// Define the Grand Staff notes (A5 to F2)
	notes := []string{
		"A5", "G5", "F5", "E5", "D5", "C5", "B4", "A4", "G4", "F4", "E4", // Treble
		"D4", "C4", "B3", // Middle
		"A3", "G3", "F3", "E3", "D3", "C3", "B2", "A2", "G2", "F2", // Bass
	}
	notePositions := make([]NotePosition, len(notes))
	// notePositions := make([]NotePosition, 24) // is an equivalent form, since len(notes) = 24
	/*
					type NotePosition struct {
					    Pitch string // such as "A5", "G5", "F5", or "E5"
					    Y     float32 // Y-coordinate on the canvas for this note
					}
				make is a built-in function used to create and initialize certain built-in types: slices, maps, and channels. When
		you see make([]Type, length), it creates a slice of type []Type with a specified length (and optionally a capacity, if provided
		as a third argument).
		[]NotePosition: Specifies the slice type—elements are NotePosition structs.
		len(notes): Sets the length of the slice to 24 (since len(notes) is 24).
		Result::: notePositions is a slice of 24 NotePosition elements (structs), pre-allocated and initialized with zero values for the
		::: type (Pitch: "", Y: 0.0 for each element).

	*/

	// Load the empty notePositions slice:
	// Calculate and set the Y-axis coordinates for each note to match the staff layout. We could have hardcoded each, but calculating them is both fun and less error-prone!
	for i, note := range notes { // "i" will become 0 through 23
		if i < 11 { // for the first 11 lines/notes, Treble (A5 to E4), calculate and assign each note its position on the Y-axis.
			notePositions[i] = NotePosition{Pitch: note, Y: float32(40 + i*30)} // here "i" is 0 for the first iteration...
			// ... e.g., when i=0, Y=40 A5; i=1, Y=70 G5, i=2, Y=100 F5 // ::: so F5 is at 100 here too!
			// when i=3, 40+(3*30)=130 E5 -- i=4 160 D5 -- i=5 190 C5 -- i=6 220 B4 -- i=7 250 A4 -- i=8 280 G4 -- i=9 310 F4 -- i=10 Y= 340 E4
		} else if i < 13 { // for the next two lines/notes, (D4 and C4), calculate and assign their positions on the Y-axis.
			notePositions[i] = NotePosition{Pitch: note, Y: float32(370 + (i-11)*30)}
		} else if i == 13 { // the Y-axis of the B3 note/line is unique...
			notePositions[i] = NotePosition{Pitch: note, Y: 490} // ... so this one gets hardcoded as Y = 490
		} else { // the remaining notes on the Bass clef, (A3 to F2), are calculated with respect to the iterated value of "i", thusly:
			notePositions[i] = NotePosition{Pitch: note, Y: float32(520 + (i-14)*30)}
		}
		/*
					Results:
					i=0:  {Pitch: "A5", Y: 40}
			... see above for i=1 to i=9
					i=10: {Pitch: "E4", Y: 340}
					i=11: {Pitch: "D4", Y: 370}
					i=12: {Pitch: "C4", Y: 400}

					i=13: {Pitch: "B3", Y: 490}
					i=14: {Pitch: "A3", Y: 520}
					i=23: {Pitch: "F2", Y: 790}

				or:

					notePositions[0] = {Pitch: "A5", Y: 40}
					notePositions[1] = {Pitch: "G5", Y: 70}
					...
					notePositions[23] = {Pitch: "F2", Y: 790}
		*/
	}

	// Pick a note, randomly, for the player to place at each of its proper locations on the Grand Staff
	// rand.Seed(time.Now().UnixNano()) // deprecated, not needed
	targetNoteLetter := []string{"C", "D", "E", "F", "G", "A", "B"}[rand.Intn(7)] // [rand.Intn(7)] uses a random number as index to the slice.
	var targetPositions []NotePosition                                            // NotePosition is a custom type, and targetPositions is then a new empty slice of those types.
	// was: targetPositions := []NotePosition{} // NotePosition is a custom type, and targetPositions is then a new empty slice of those types.
	// var is just a declaration (nil slice), while := initializes an empty slice. Functionally identical here since we append right away.
	for _, pos := range notePositions { // notePositions is a slice of 24 NotePosition elements (structs); each now loaded with a ...
		// ... Y coordinate. ::: "pos" is therefore a struct of type NotePosition; each containing one of those Y-axis coordinate values which
		// .. was calculated above. And, we toss the unneeded range position via the built-in _ bit bin variable -- “blank identifier” (Go term for _).
		if pos.Pitch[0:1] == targetNoteLetter { // targetNoteLetter could be any of C to B, as per the randomly indexed slice above. And ...
			/*
				pos could be, e.g., {Pitch: "A5", Y: 40}  pos.Pitch returns the Pitch field: "A5"; Whereas pos.Y would return the Y field: 40.
				pos.Pitch is a string like "A5". And [0:1]: is a slice expression applied to that string. In Go, strings are sliceable, meaning
					you can always extract a substring using the syntax string[start:end]
					e.g., s[start:end] extracts a portion of s from index start (inclusive) to index end (exclusive).
					::: So, pos.Pitch[0:1] trims off the octave identifier found in each NotePosition struct; turning things like "C4" to just plain "C".
			*/
			targetPositions = append(targetPositions, pos) // load/add/append pos to the targetPositions slice. pos having been found by range.
			// Appending: Collects all NotePosition structs where the note letter matches (e.g., all "C" positions).
		}
	}
	/*
			Example Run:
			targetNoteLetter = "C" (randomly picked).

			Loop checks notePositions:
			pos.Pitch = "A5": "A5"[0:1] = "A" ≠ "C", skip.
		...
			pos.Pitch = "C5": "C5"[0:1] = "C" == "C", append {Pitch: "C5", Y: 160}.
			pos.Pitch = "C4": "C4"[0:1] = "C" == "C", append {Pitch: "C4", Y: 400}.
			pos.Pitch = "C3": "C3"[0:1] = "C" == "C", append {Pitch: "C3", Y: 670}.
		...
			Output of the following Printf statement: Target C notes: [{C5 160} {C4 400} {C3 670}].
	*/
	fmt.Printf("Target %s notes: %v\n", targetNoteLetter, targetPositions) // log activity to the console/terminal.

	// Create canvas for the staff: canvas.___ is a fyne object. Compare Fyne calls near top of main.
	staffCanvas := canvas.NewRectangle(&color.RGBA{R: 255, G: 255, B: 255, A: 255})
	staffCanvas.Resize(fyne.NewSize(1000, 900)) // size of actual app window. canvas is (1065, 980)
	// obviously, the canvas fits inside the window

	// Draw the Grand Staff
	lines := []fyne.CanvasObject{staffCanvas}
	// Treble staff (E4 bottom, F5 top)
	for i := 0; i < 5; i++ {
		y := float32(340 - i*60) // E4 (340), G4 (280), B4 (220), D5 (160), F5 (100) // F5 at 100 ?? i=9 Y=310  -- i=9 310 F5 --  ???
		// when i=4, 370 - i*60 = 130
		line := canvas.NewLine(&color.Black)
		line.Position1 = fyne.NewPos(100, y)
		line.Position2 = fyne.NewPos(900, y)
		line.StrokeWidth = 2
		lines = append(lines, line)
	}
	/*
		i=0 Y=340 the E4 line bottom line of treble clef.  i=1 280, 2 220, 3 160, 4 100 (for i := 0; i < 5; i++ {)
		Matches the code (340 - i*60), but note these don’t align perfectly with notePositions (e.g., F5 is 100 here, but 130 in notePositions). Intentional offset?
	*/
	// Bass staff (G2 bottom, A3 top)
	for i := 0; i < 5; i++ {
		y := float32(760 - i*60) // G2 (760), B2 (700), D3 (640), F3 (580), A3 (520)
		line := canvas.NewLine(&color.Black)
		line.Position1 = fyne.NewPos(100, y)
		line.Position2 = fyne.NewPos(900, y)
		line.StrokeWidth = 2
		lines = append(lines, line)
	}
	// Middle C ledger line (C4)
	for x := 400; x < 600; x += 20 {
		ledger := canvas.NewLine(&color.Black)
		ledger.Position1 = fyne.NewPos(float32(x), 400)
		ledger.Position2 = fyne.NewPos(float32(x+10), 400)
		ledger.StrokeWidth = 2
		lines = append(lines, ledger)
	}
	// A5 ledger line (above G5)
	for x := 400; x < 600; x += 20 {
		ledger := canvas.NewLine(&color.Black)
		ledger.Position1 = fyne.NewPos(float32(x), 40)
		ledger.Position2 = fyne.NewPos(float32(x+10), 40)
		ledger.StrokeWidth = 2
		lines = append(lines, ledger)
	}

	// ::: Track player's marked notes—where they’ve placed circles.
	markedNotes := []MarkedNote{} // empty slice declaration using literal {}
	//  it’s for storing MarkedNote structs from clicks.

	// Create staff container (a fyne object to hold staff lines and notes)
	staffContainer := container.NewWithoutLayout(lines...)
	staffContainer.Resize(fyne.NewSize(1000, 900)) // same dimensions as staffCanvas.Resize(fyne.NewSize(
	staffContainer.Move(fyne.NewPos(100, 100))     // meant to center the canvas (staff container) in the window

	// Handle mouse clicks with a tappable rectangle (more fyne objects)
	staffArea := canvas.NewRectangle(&color.Transparent)
	staffArea.Resize(fyne.NewSize(1000, 900)) // same dimensions as staffCanvas.Resize(fyne.NewSize(

	// Add tap handler (this is a big one, approximately 50 lines)
	staffAreaTapped := &TappableCanvas{ // &TappableCanvas is a pointer address to a custom type, i.e., ...
		// TappableCanvas extends CanvasObject to handle taps, far below.
		// Snaps to closest Y from notePositions
		CanvasObject: staffArea,
		OnTapped: func(e *fyne.PointEvent) {
			clickX, clickY := e.Position.X, e.Position.Y

			// Check if clicking an existing note to remove it
			for i, note := range markedNotes {
				dx := clickX - note.X
				dy := clickY - note.Y
				distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
				if distance < 10 { // ::: Within 10px radius of note center
					staffContainer.Remove(note.Circle)
					markedNotes = append(markedNotes[:i], markedNotes[i+1:]...)
					staffContainer.Refresh()
					fmt.Printf("Removed note at X=%.0f, Y=%.0f\n", note.X, note.Y)
					return
				}
			}

			// Snap to nearest note position
			closest := notePositions[0]
			minDiff := abs(clickY - closest.Y)
			for _, pos := range notePositions {
				diff := abs(clickY - pos.Y)
				if diff < minDiff {
					minDiff = diff
					closest = pos
				}
			}

			// Determine X position: ledger (center) or staff (right)
			var noteX float32
			if closest.Pitch == "A5" || closest.Pitch == "C4" {
				noteX = 500 // Center of ledger lines (400-600)
			} else {
				noteX = 300 // Halfway between staff left (100) and ledger left (400)
			}

			// Add a red circle/dot
			circle := canvas.NewCircle(&color.RGBA{R: 255, G: 0, B: 0, A: 255})
			circle.Resize(fyne.NewSize(20, 20))
			circle.Move(fyne.NewPos(noteX-10, closest.Y-10)) // Center circle on position
			markedNotes = append(markedNotes, MarkedNote{Circle: circle, X: noteX, Y: closest.Y})
			// staffContainer.AddObject(circle)  // AddObject is deprecated. ?? what could go in its place ???
			staffContainer.Add(circle)                                                   // Add replaces deprecated AddObject—keeps it modern!
			staffContainer.Refresh()                                                     // staffContainer is an instance of container.NewWithoutLayout(lines...) , done above.
			fmt.Printf("Marked %s at X=%.0f, Y=%.0f\n", closest.Pitch, noteX, closest.Y) // debugging log to terminal.
		},
	}
	staffContainer.Add(staffAreaTapped) // staffContainer is an instance of container.NewWithoutLayout(lines...) , done above.

	// Instruction and feedback
	instruction := widget.NewLabel(fmt.Sprintf("Click all %s notes on the Grand Staff", targetNoteLetter))
	instruction.TextStyle = fyne.TextStyle{Bold: true}
	instruction.Resize(fyne.NewSize(1000, 100))
	feedback := canvas.NewText("", color.Black)
	feedback.TextSize = 24
	feedback.TextStyle = fyne.TextStyle{Bold: true}

	// Define content container
	content := container.NewVBox()

	// Check button
	var checkButton *widget.Button
	checkButton = widget.NewButton("Check", func() {
		correctCount := 0
		wrongCount := 0
		coveredTargets := make(map[float32]bool) // Track unique target Y positions hit
		uniqueMarks := make(map[string]bool)     // Track unique X,Y positions marked

		// Deduplicate marks based on X,Y position
		dedupedMarks := []MarkedNote{}
		for _, mark := range markedNotes {
			key := fmt.Sprintf("%.0f,%.0f", mark.X, mark.Y)
			if !uniqueMarks[key] {
				uniqueMarks[key] = true
				dedupedMarks = append(dedupedMarks, mark)
			}
		}

		fmt.Println("Check clicked")
		for _, mark := range dedupedMarks {
			isCorrect := false
			for _, target := range targetPositions {
				if abs(mark.Y-target.Y) < 15 && !coveredTargets[target.Y] {
					isCorrect = true
					coveredTargets[target.Y] = true // Mark this target as covered
					break
				}
			}
			if isCorrect {
				correctCount++
			} else {
				wrongCount++
			}
		}

		missing := len(targetPositions) - correctCount
		msg := fmt.Sprintf("Found %d/%d %s notes, %d wrong positions", correctCount, len(targetPositions), targetNoteLetter, wrongCount)
		if missing == 0 && wrongCount == 0 {
			msg = fmt.Sprintf("Perfect! All %s notes found!", targetNoteLetter) // Success message to player.
			checkButton.Disable()
		}
		fmt.Println(msg)
		feedback.Text = msg
		feedback.Refresh()
		content.Refresh()
	})

	// Check button—tallies player’s note placements.
	checkButton.Resize(fyne.NewSize(200, 100)) // Check button’s footprint.
	// ::: Success message: "Perfect! All %s notes found!" if all targets hit. (line 292)

	// Reset button (aka New Game)—wipes slate clean for a fresh challenge.
	resetButton := widget.NewButton("New Game", func() {
		targetNoteLetter = []string{"C", "D", "E", "F", "G", "A", "B"}[rand.Intn(7)]
		targetPositions = []NotePosition{}
		for _, pos := range notePositions {
			if pos.Pitch[0:1] == targetNoteLetter {
				targetPositions = append(targetPositions, pos)
			}
		}
		fmt.Printf("Target %s notes: %v\n", targetNoteLetter, targetPositions)
		for _, mark := range markedNotes {
			staffContainer.Remove(mark.Circle)
		}
		markedNotes = []MarkedNote{}
		instruction.SetText(fmt.Sprintf("Click all %s notes on the Grand Staff", targetNoteLetter))
		feedback.Text = ""
		checkButton.Enable()
		staffContainer.Refresh()
		content.Refresh()
	})
	resetButton.Resize(fyne.NewSize(200, 100)) // dimensions of "New Game" button.

	// Populate content container
	content.Objects = []fyne.CanvasObject{
		instruction,
		staffContainer,
		container.NewHBox(checkButton, resetButton),
		feedback,
	}

	// Main layout
	mainContainer := container.New(layout.NewVBoxLayout(), content)

	// Set up window
	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
} // ::: end of main

// TappableCanvas extends CanvasObject to handle taps
type TappableCanvas struct { // GoLand adds "go to interfaces" widget in margin.
	fyne.CanvasObject
	OnTapped func(*fyne.PointEvent)
}

func (t *TappableCanvas) Tapped(e *fyne.PointEvent) { // GoLand adds "go to interfaces" widget in margin.
	if t.OnTapped != nil {
		t.OnTapped(e)
	}
}

// abs returns the absolute value of a float32
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
