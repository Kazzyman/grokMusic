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
	"time"
)

// NotePosition represents a note's position on the staff
type NotePosition struct {
	Pitch string
	Y     float32 // Y-coordinate on the canvas
}

// MarkedNote tracks a placed note for retraction
type MarkedNote struct {
	Circle *canvas.Circle
	X      float32
	Y      float32
}

func main() {
	// Initialize Fyne app
	myApp := app.New()
	myWindow := myApp.NewWindow("Find the Note")
	myWindow.Resize(fyne.NewSize(1200, 1000))

	// Define the Grand Staff notes (A5 to F2)
	notes := []string{
		"A5", "G5", "F5", "E5", "D5", "C5", "B4", "A4", "G4", "F4", "E4", // Treble
		"D4", "C4", "B3", // Middle
		"A3", "G3", "F3", "E3", "D3", "C3", "B2", "A2", "G2", "F2", // Bass
	}
	notePositions := make([]NotePosition, len(notes))
	for i, note := range notes {
		if i < 11 { // Treble (A5 to E4)
			notePositions[i] = NotePosition{Pitch: note, Y: float32(40 + i*30)}
		} else if i < 13 { // Middle (D4 to C4)
			notePositions[i] = NotePosition{Pitch: note, Y: float32(370 + (i-11)*30)}
		} else if i == 13 { // B3
			notePositions[i] = NotePosition{Pitch: note, Y: 490}
		} else { // Bass (A3 to F2)
			notePositions[i] = NotePosition{Pitch: note, Y: float32(520 + (i-14)*30)}
		}
	}

	// Pick a random note to find
	rand.Seed(time.Now().UnixNano())
	targetNoteLetter := []string{"C", "D", "E", "F", "G", "A", "B"}[rand.Intn(7)]
	targetPositions := []NotePosition{}
	for _, pos := range notePositions {
		if pos.Pitch[0:1] == targetNoteLetter {
			targetPositions = append(targetPositions, pos)
		}
	}
	fmt.Printf("Target %s notes: %v\n", targetNoteLetter, targetPositions)

	// Create canvas for the staff
	staffCanvas := canvas.NewRectangle(&color.RGBA{R: 255, G: 255, B: 255, A: 255})
	staffCanvas.Resize(fyne.NewSize(1000, 900))

	// Draw the Grand Staff
	lines := []fyne.CanvasObject{staffCanvas}
	// Treble staff (E4 bottom, F5 top)
	for i := 0; i < 5; i++ {
		y := float32(340 - i*60) // E4 (340), G4 (280), B4 (220), D5 (160), F5 (100)
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

	// Track player's marked notes
	markedNotes := []MarkedNote{}

	// Create staff container
	staffContainer := container.NewWithoutLayout(lines...)
	staffContainer.Resize(fyne.NewSize(1000, 900))
	staffContainer.Move(fyne.NewPos(100, 100))

	// Handle mouse clicks with a tappable rectangle
	staffArea := canvas.NewRectangle(&color.Transparent)
	staffArea.Resize(fyne.NewSize(1000, 900))

	// Add tap handler
	staffAreaTapped := &TappableCanvas{
		CanvasObject: staffArea,
		OnTapped: func(e *fyne.PointEvent) {
			clickX, clickY := e.Position.X, e.Position.Y

			// Check if clicking an existing note to remove it
			for i, note := range markedNotes {
				dx := clickX - note.X
				dy := clickY - note.Y
				distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
				if distance < 10 { // Within 10px radius of note center
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

			// Add a red circle
			circle := canvas.NewCircle(&color.RGBA{R: 255, G: 0, B: 0, A: 255})
			circle.Resize(fyne.NewSize(20, 20))
			circle.Move(fyne.NewPos(noteX-10, closest.Y-10)) // Center circle on position
			markedNotes = append(markedNotes, MarkedNote{Circle: circle, X: noteX, Y: closest.Y})
			staffContainer.AddObject(circle)
			staffContainer.Refresh()
			fmt.Printf("Marked %s at X=%.0f, Y=%.0f\n", closest.Pitch, noteX, closest.Y)
		},
	}
	staffContainer.Add(staffAreaTapped)

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
			msg = fmt.Sprintf("Perfect! All %s notes found!", targetNoteLetter)
			checkButton.Disable()
		}
		fmt.Println(msg)
		feedback.Text = msg
		feedback.Refresh()
		content.Refresh()
	})
	checkButton.Resize(fyne.NewSize(200, 100))

	// Reset button
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
	resetButton.Resize(fyne.NewSize(200, 100))

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
}

// TappableCanvas extends CanvasObject to handle taps
type TappableCanvas struct {
	fyne.CanvasObject
	OnTapped func(*fyne.PointEvent)
}

func (t *TappableCanvas) Tapped(e *fyne.PointEvent) {
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
