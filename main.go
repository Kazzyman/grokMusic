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

// @formatter:off
// GoLand is configured so as to highlight (in pink) all characters on a line to the right of three colons, e.g. ::: this is pink

// NotePosition represents a note's position on the staff; implements interface
type NotePosition struct {
	Pitch string  // such as "A5", "G5", "F5", or "E5"
	Y     float32 // Y-coordinate on the canvas for this note
}

// MarkedNote tracks a placed note for possible retraction in case of player error; implements interface
type MarkedNote struct {
	Circle *canvas.Circle
	X      float32
	Y      float32
}

func main() { 
	about_app() // show SLOC on the terminal; and, maintain a log file: musicAppLog.txt where those LOC figures are tracked. 
	
	// Initialize Fyne app  -- app.___ is a Fyne object.
	RicksFirstGUI := app.New()
	parentWindow := RicksFirstGUI.NewWindow("Rick's Find the Note game") // create the app window and title it.
	parentWindow.Resize(fyne.NewSize(1000, 1000)) 

	// Define the Grand Staff notes (A5 to F2)
	notes := []string{
		"A5", "G5", "F5", "E5", "D5", "C5", "B4", "A4", "G4", "F4", "E4", // Treble
		"D4", "C4", "B3", // Middle
		"A3", "G3", "F3", "E3", "D3", "C3", "B2", "A2", "G2", "F2", // Bass
	}
	notePositions := make([]NotePosition, len(notes))
	/*
	make is a built-in function used to create and initialize certain built-in types: slices, maps, and channels. When
		you see make([]Type, length), it creates a slice of type []Type with a specified length (and optionally a capacity, if provided
		as a third argument).
		[]NotePosition: Specifies the slice type—elements are NotePosition structs.
		len(notes): Sets the length of the slice to 24 (since len(notes) is 24).
		Result::: notePositions is a slice of 24 NotePosition elements (structs), pre-allocated and initialized with zero values for the
		::: type (Pitch: "", Y: 0.0 for each element).
	*/

	// Load the empty notePositions slice:
	// Calculate and set the Y-axis coordinates for each note to match the staff layout. We could have hardcoded each, but calculating ...
	// ... them is both fun and less error-prone!
	for i, note := range notes { // "i" will become 0 through 23
		if i < 11 { // for the first 11 lines/notes, Treble (A5 to E4), calculate and assign each note its position on the Y-axis.
			notePositions[i] = NotePosition{Pitch: note, Y: float32(40 + i*30)} // here "i" is 0 for the first iteration...
			// ... e.g., when i=0, Y=40 A5; i=1, Y=70 G5, i=2, Y=100 F5
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
	targetNoteLetter := []string{"C", "D", "E", "F", "G", "A", "B"}[rand.Intn(7)] // [rand.Intn(7)] uses a random number as index to the slice.
	var targetPositions []NotePosition  // NotePosition is a custom type, and targetPositions is then a new empty slice of those types.
	// was: targetPositions := []NotePosition{} // NotePosition is a custom type, and targetPositions is then a new empty slice of those types.
	// var is just a declaration (nil slice), while := initializes an empty slice. Functionally identical here since we append right away.
	for _, pos := range notePositions { // notePositions is a slice of 24 NotePosition elements (structs); each now loaded with a ...
		// ... Y coordinate. ::: "pos" is therefore a struct of type NotePosition; each containing one of those Y-axis coordinate values which
		// .. was calculated above. And, we toss the unneeded range position via the built-in _ bit bin variable -- “blank identifier” (Go term for _).
		if pos.Pitch[0:1] == targetNoteLetter { // targetNoteLetter could be any of C to B, as per the randomly indexed slice above. And ...
			/*
			pos could be, e.g., {Pitch: "A5", Y: 40}  pos.Pitch returns the Pitch field: "A5"; Whereas pos.Y would return the Y field: 40.
			pos.Pitch is a string like "A5". And [0:1]: is a slice expression applied to that string. In Go, strings are sliceable, meaning that
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
	staffCanvas := canvas.NewRectangle(&color.RGBA{R: 25, G: 200, B: 25, A: 155})
	staffCanvas.Resize(fyne.NewSize(1000, 1000)) // Needed to apply the colors specified on the previous line; default is very dark grey. 

	// Draw the Grand Staff
	lines := []fyne.CanvasObject{staffCanvas}
	// Treble staff (E4 bottom, F5 top)
	for i := 0; i < 5; i++ {
		y := float32(340 - i*60) // E4 (340), G4 (280), B4 (220), D5 (160), F5 (100) // F5 at 100 ?? i=9 Y=310  -- i=9 310 F5 -- 
		// when i=4, 370 - i*60 = 130
		line := canvas.NewLine(&color.Black)
		line.Position1 = fyne.NewPos(100, y)
		line.Position2 = fyne.NewPos(900, y)
		line.StrokeWidth = 2
		lines = append(lines, line)
	}

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

	// ::: Track the player's marked notes — places where they’ve placed circles/dots.
	markedNotes := []MarkedNote{} // empty slice declaration using literal {}
	// could also have done a: var markedNotes []MarkedNote // var is just a declaration (nil slice), while := initializes an empty slice.
	// it’s for storing MarkedNote structs from clicks.

	// Create staff container (a fyne object to hold staff lines and notes)
	staffContainer := container.NewWithoutLayout(lines...)
	// No Resize/Move statement for staffContainer — VBox in content overrides these!

	// Handle mouse clicks with a tappable rectangle (more fyne objects)
	staffArea := canvas.NewRectangle(&color.Transparent) // invisible overlay for detecting mouse clicks

	staffArea.Resize(fyne.NewSize(1000, 1000)) // Needed, makes the tappable area cover the entire staff drawing surface

	// Add tap handler (this is a big one, approximately 50 lines). This is our custom tap-handling callback func -- staffAreaTapped is the instance (via the receiver t in Tapped())
	staffAreaTapped := &TappableCanvas{ // &TappableCanvas is a pointer address to a custom type: TappableCanvas extends CanvasObject to handle taps (see below)
		// It snaps to closest Y from notePositions. Snaps clicks to the nearest Y from notePositions -- staffAreaTapped is the instance (via the receiver t in Tapped())
		/* 
			&TappableCanvas points to our custom TappableCanvas type, extending CanvasObject to include tap glory (see TappableCanvas below).
		 */
		/*
		How Snapping Works:
		notePositions Recap:
		Defined earlier as a slice of NotePosition structs: []NotePosition{Pitch: string, Y: float32}.

		Maps 24 notes (A5 to F2) to Y-coordinates (e.g., A5=40, G5=70, F5=100, …, F2=790).

		Example: notePositions[2] = {Pitch: "F5", Y: 100}.

		Click Input:
		clickY (from e.Position.Y) is the raw Y-coordinate where the player taps—could be anywhere (e.g., 153.7).

		Finding the Closest:
		Initial Guess: Starts with closest = notePositions[0] (A5, Y=40) and minDiff = abs(clickY - 40).

		Loop: Iterates over notePositions, calculating diff = abs(clickY - pos.Y) for each.

		Update: If diff < minDiff, updates minDiff and closest to that position.

		Result: closest ends up as the NotePosition with the Y-value nearest to clickY.
		Example: If clickY = 153.7:
		abs(153.7 - 100) = 53.7 (F5)

		abs(153.7 - 130) = 23.7 (E5)

		abs(153.7 - 160) = 6.3 (D5) → Wins! closest = {Pitch: "D5", Y: 160}.

		X Positioning:
		Hardcodes noteX:
		500 for A5 or C4 (ledger line center, 400-600 range).

		300 for others (mid-staff, between left edge 100 and ledger start 400).

		Why? Ensures notes visually align with staff or ledger lines, though it’s a simplification (no X-snapping).

		Placing the Note:
		Creates a red circle, positions it at (noteX-10, closest.Y-10) (centers the 20x20 circle), and adds it to markedNotes.

		Why Snap?
		Precision: Players don’t need pixel-perfect clicks—snapping makes it forgiving (e.g., clicking Y=153.7 snaps to D5 at 160).

		Game Logic: Ensures notes land on valid staff positions from notePositions, matching targetPositions for scoring.
		*/
		CanvasObject: staffArea, // a struct literal field initialization with an embedded field ...
		// CanvasObject: staffArea, Embeds staffArea (a transparent rectangle) as the drawable CanvasObject — makes it tappable and visible
		OnTapped: func(e *fyne.PointEvent) { // OnTapped is the callback func "from" the TappableCanvas struct which is an extended instance of CanvasObject.
			// ... It sets OnTapped, the tap-handling callback in TappableCanvas — extending CanvasObject with our click magic!
			clickX, clickY := e.Position.X, e.Position.Y // e is the argument passed to the OnTapped callback func that we are defining here. 
			// e is a *fyne.PointEvent, a struct with fields like Position (a fyne.Position with X and Y floats). It’s the event data—where the player clicked.
			// e.Position.X and e.Position.Y extract the click coordinates. e is the tap event (*fyne.PointEvent) — grabs X/Y coords
			
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

			// Snap to nearest note position — finds the closest Y from notePositions for that perfect note placement!
			closest := notePositions[0] // Start with A5 (Y=40) as our guess.
			minDiff := abs(clickY - closest.Y) // How far’s the click from A5?
			for _, pos := range notePositions { // Loop through all 24 notes (A5 to F2).
			    diff := abs(clickY - pos.Y) // Distance from click to this note’s Y.
			    if diff < minDiff { // Closer than our last best? Update!
			        minDiff = diff
			        closest = pos // New champ—e.g., click Y=153.7 snaps to D5 (Y=160).
			    }
			}

		/*
			This is a classic “nearest neighbor” algorithm—simple yet effective. It’s forgiving (no threshold—always snaps), but 
			you could add if minDiff < 20 to limit snapping range if desired.
		*/
			// Determine X position: ledger (center) or staff (right)
			var noteX float32

			if closest.Pitch == "A5" || closest.Pitch == "C4" {
				noteX = 500 // Center of ledger lines (400-600)
			} else {
				noteX = 300 // Halfway between staff left (100) and ledger left (400)
			}

			// Add a red circle/dot
			circle := canvas.NewCircle(&color.RGBA{R: 255, G: 0, B: 0, A: 255})
			circle.Resize(fyne.NewSize(20, 20)) // Needed to apply the colors specified on the previous line; default appears to be invisible ???
			circle.Move(fyne.NewPos(noteX-10, closest.Y-10)) // Center circle on position
			markedNotes = append(markedNotes, MarkedNote{Circle: circle, X: noteX, Y: closest.Y})
			staffContainer.Add(circle)  // Add replaces deprecated AddObject—keeps it modern!
			staffContainer.Refresh()   // staffContainer is an instance of container.NewWithoutLayout(lines...) , done above.
			fmt.Printf("Marked %s at X=%.0f, Y=%.0f\n", closest.Pitch, noteX, closest.Y) // debugging log to terminal.
		},
	}
	staffContainer.Add(staffAreaTapped) // staffContainer is an instance of container.NewWithoutLayout(lines...) , done above. 
	// Adds our tappable layer to staffContainer (from NewWithoutLayout above) — clicks live here!

	// Instruction and feedback
	instruction := widget.NewLabel(fmt.Sprintf("Click all %s notes on the Grand Staff", targetNoteLetter))
	instruction.TextStyle = fyne.TextStyle{Bold: true}
	// No Resize statement for instruction — VBox sets size based on text!

	feedback := canvas.NewText("", color.Black)
	feedback.TextSize = 24
	feedback.TextStyle = fyne.TextStyle{Bold: true}

	// Define content container
	content := container.NewVBox()

	// Check button — tallies player’s note placements.
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
	// No Resize statement for checkButton — HBox in content dictates button size!
	
	// Reset button (aka New Game) — wipes slate clean for a fresh challenge.
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
	// No Resize statement for resetButton — HBox in content dictates button size!
	
	// Populate content container
	content.Objects = []fyne.CanvasObject{
		instruction,
		staffContainer,
		container.NewHBox(checkButton, resetButton),
		feedback,
	}

	// Main layout
	mainContainer := container.New(layout.NewVBoxLayout(), content)

	// Set up window, and run it
	parentWindow.SetContent(mainContainer)
	parentWindow.ShowAndRun()
} // ::: end of main

// TappableCanvas extends CanvasObject to catch and handle taps (mouse clicks). // TappableCanvas extends CanvasObject to snag and handle taps (mouse clicks galore)!
type TappableCanvas struct { 
	fyne.CanvasObject // // Embeds the drawable base, size, and position (the first step towards extending it. 
	/* per grok:
	   Embedding fyne.CanvasObject (no field name, just type) gives TappableCanvas all its methods (Resize(), Move(), Refresh()) and
	   fields. It’s Go’s inheritance trick — CanvasObject is Fyne’s drawable foundation (rectangles, lines, circles). This makes
	   TappableCanvas a CanvasObject, ready for containers like staffContainer.Add(staffAreaTapped). Next, we add tap powers!
	*/
	/* said my way:
		this is a hybrid of inheritance-like behavior (via embedding) and event handling. fyne.CanvasObject is an embedded field (no name, just the type). Embedding lets
		TappableCanvas inherit all methods and fields of fyne.CanvasObject . CanvasObject is Fyne’s base interface for anything drawable, such as rectangles, lines, or
		circles. It includes methods like Resize(), Move(), and Refresh() . TappableCanvas thereby becomes a CanvasObject : and, it can therefore be added to containers
		(e.g., staffContainer.Add(staffAreaTapped) . Next, we extend it. 
	*/
	OnTapped func(*fyne.PointEvent) // A named field of the struct, OnTapped, is a function type (a callback func : for tap action; which is where the magic happens!). 
	// It takes a *fyne.PointEvent (a pointer to a struct with X/Y coordinates and other event data) and returns nothing. This is the handler you’ll define later : 
	// it defines what happens when the player clicks on the canvas.
	/* per grok:
	OnTapped func(*fyne.PointEvent) // Named field — a callback function for tap action, where the magic unfolds!
	Takes *fyne.PointEvent (X/Y coords and event data), returns nada. Set this later to define player tap behavior!
	 */
}

// Tapped is a declared method, (t *TappableCanvas) is the method receiver; or method receiver declaration. It creates a local var 't' as a pointer to an instance of TappableCanvas ...
// t functions as the instance (below in t.OnTapped(e) 
/* per grok:
	Tapped, a method with receiver (t *TappableCanvas)—declares ‘t’ as a pointer to our TappableCanvas instance!
 */
func (t *TappableCanvas) Tapped(e *fyne.PointEvent) { // Implements Fyne’s Tappable interface; tap-ready! Notice that TappableCanvas became an extended CanvasObject -- above.
	/* And; that other guy ---- ^ e will be used above as: OnTapped: func(e *fyne.PointEvent) { // e is the argument passed to OnTapped
	This is a method on TappableCanvas with a receiver t (the instance being tapped). Named Tapped : this is key! It implements Fyne’s Tappable interface, which 
	requires a Tapped(*fyne.PointEvent) method. Notice that it takes the same *fyne.PointEvent as the OnTapped field/(callback func) : the coordinates of the tap or PointEvent.
	*/
	/* per grok: 
	   Method on TappableCanvas — receiver ‘t’ is the tapped instance. ‘Tapped’ name is key: satisfies Fyne’s Tappable interface
	   with Tapped(*fyne.PointEvent). Matches OnTapped’s *fyne.PointEvent for tap coords!
	*/
	// t below is a local variable: and contains a pointer to a type (in this case the user-defined struct TappableCanvas). Which is an extended version of CanvasObject.
	/* per grok:
	t points to a TappableCanvas instance -- ‘t’ is the receiver — a pointer to this TappableCanvas instance, extended from CanvasObject.
	*/
	if t.OnTapped != nil { // This checks to assure that OnTapped was set (is not nil). In Go, function types default to nil if unassigned, preventing 
		// a panic crash that would result from calling a null/unset function. Safety is hereby assured; we'll have no nil crashes here! Proceed only if not nil.
		t.OnTapped(e) // This calls "back" to the stored OnTapped tap-handling function of the preceding struct, passing it the tap event (e). This delegates the actual ...
		// ... tap logic to whatever you plugged in (e.g., your note-placing code). It thereby fires off our custom tap handling logic.
		/* per grok:
		t.OnTapped(e) // Fires the stored OnTapped callback with tap event ‘e’ — unleashes your custom logic!
		 */
	}
}
/* grok does a recap:
// Add tap handler—a hefty ~50-line beast! Our custom tap-handling callback awaits.
staffAreaTapped := &TappableCanvas{ // &TappableCanvas points to our custom TappableCanvas type, extending CanvasObject for tap glory (see below).
    // Snaps clicks to the nearest Y from notePositions—precision meets play!
    CanvasObject: staffArea, // Embeds staffArea (a transparent rectangle) as the drawable CanvasObject—makes it tappable and visible!
    OnTapped: func(e *fyne.PointEvent) { // Sets OnTapped, the tap-handling callback in TappableCanvas—extending CanvasObject with click magic!
        clickX, clickY := e.Position.X, e.Position.Y // e is the tap event (*fyne.PointEvent)—grabs X/Y coords, not the instance (that’s t)!
        // Check if clicking an existing note to remove it
        for i, note := range markedNotes { ...
    },
}
staffContainer.Add(staffAreaTapped) // Adds our tappable layer to staffContainer (from NewWithoutLayout above)—clicks live here!

// TappableCanvas extends CanvasObject to snag and handle taps (mouse clicks galore)!
type TappableCanvas struct {
    fyne.CanvasObject // Embeds the drawable base—size, position, the works—step one to extending it!
    / *
        Embedding fyne.CanvasObject (no field name, just type) gives TappableCanvas all its methods (Resize(), Move(), Refresh()) and
        fields. It’s Go’s inheritance trick—CanvasObject is Fyne’s drawable foundation (rectangles, lines, circles). This makes
        TappableCanvas a CanvasObject, ready for containers like staffContainer.Add(staffAreaTapped). Next, we add tap powers!
	* /
OnTapped func(*fyne.PointEvent) // Named field—a callback function for tap action, where the magic unfolds!
// Takes *fyne.PointEvent (X/Y coords and event data), returns nada. Set this later to define player tap behavior!
}

// Tapped, a method with receiver (t *TappableCanvas)—declares ‘t’ as a pointer to our TappableCanvas instance!
func (t *TappableCanvas) Tapped(e *fyne.PointEvent) { // Implements Fyne’s Tappable interface—tap-ready! Extends CanvasObject (see above).
	/ *
	   Method on TappableCanvas—receiver ‘t’ is the tapped instance. ‘Tapped’ name is key: satisfies Fyne’s Tappable interface
	   with Tapped(*fyne.PointEvent). Matches OnTapped’s *fyne.PointEvent for tap coords!
	* /
	// ‘t’ is the receiver—a pointer to this TappableCanvas instance, extended from CanvasObject.
	if t.OnTapped != nil { // Checks OnTapped isn’t nil—avoids panic crashes from unset funcs. Safety first!
		t.OnTapped(e) // Fires the stored OnTapped callback with tap event ‘e’—unleashes your custom logic!
	}
}
 */

// grok explains all about interfaces:
/*
In Go, interfaces are implicit contracts — a set of method signatures a type must implement. Unlike Java or C#, you don’t declare “implements” — 
if a type has the methods, it satisfies the interface. Fyne leans heavily on this for its UI framework.

Syntax:
			type MyInterface interface {
				Method1()
				Method2(arg string) int
			}

Key: Any type with Method1() and Method2(string) int is a MyInterface—no explicit tie needed.

Fyne’s Core Interfaces:
Fyne uses interfaces to define behaviors for UI elements. Here’s how they fit your TappableCanvas:
fyne.CanvasObject:

Definition:
			type CanvasObject interface {
				Size() fyne.Size
				Resize(size fyne.Size)
				Position() fyne.Position
				Move(position fyne.Position)
				Visible() bool
				Show()
				Hide()
				Refresh()
			}

Purpose: Anything drawable on the canvas—rectangles, lines, text, etc.

Your Code:
	staffArea (a *canvas.Rectangle) implements this.

Embedding fyne.CanvasObject in TappableCanvas means it inherits these methods (e.g., Resize(1000, 900)).

Why: Ensures staffAreaTapped can be sized, positioned, and drawn in staffContainer.

fyne.Tappable:

Definition:
			type Tappable interface {
				Tapped(*fyne.PointEvent)
			}

Purpose: Makes an object respond to tap/click events (mouse or touch).

Your Code:
	TappableCanvas implements this with:

			func (t *TappableCanvas) Tapped(e *fyne.PointEvent) {
				if t.OnTapped != nil {
					t.OnTapped(e)
				}
			}

Fyne’s event system calls Tapped() when staffAreaTapped is clicked.

Why: Adds interactivity—without it, staffArea would just sit there, pretty but mute.

How They Work Together
Embedding:
	TappableCanvas embeds fyne.CanvasObject, so it’s drawable (via staffArea).

Adding Tapped() makes it Tappable too.

Result: A single type satisfying two interfaces—visual and interactive.

Fyne’s Event Loop:
When you tap the staff:
	Fyne checks if the clicked object implements Tappable.

Finds staffAreaTapped (topmost in staffContainer).

Calls Tapped(e), passing the click’s *fyne.PointEvent.

Your OnTapped runs, snapping to notePositions.

Other Fyne Interfaces ::: (Context)
		fyne.Widget (superset of CanvasObject):
		Adds CreateRenderer() for complex widgets (e.g., buttons, labels).
		
		You didn’t need this — your Rectangle is simpler. “Rectangle” being shorthand for the type canvas.Rectangle — the thing you’re working with — rather than 
		the constructor function canvas.NewRectangle().

you have:
	staffCanvas := canvas.NewRectangle(&color.RGBA{R: 255, G: 255, B: 255, A: 255}) // use of a constructor function 
	staffArea := canvas.NewRectangle(&color.Transparent) // invisible overlay for detecting mouse clicks


fyne.Draggable:
Dragged(*fyne.DragEvent)—could extend your game to drag notes!

fyne.Focusable:
	FocusGained(), FocusLost()—for keyboard input (not used here).

Why Interfaces?
	Loose Coupling: Fyne doesn’t care how TappableCanvas works—just that it has Tapped().
	Flexibility: Swap staffArea for a canvas.Circle—still works as a CanvasObject.
	Go Idiomatic: No inheritance mess—just compose behaviors via interfaces.

*/


// abs returns the absolute value of a float32
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
