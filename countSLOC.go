package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var countOffilesExplored int

func about_app() { // ::: - -
	countSLOC()
	fmt.Printf("\n\tThis app consists of the following lines of code across %d files:\n", countOffilesExplored)
	fmt.Println("\t\tmain.go, countSLOC.go, and go.mod\n")
	fmt.Printf("\tThe above figues were calculated by reading those %d files, \n\tin real-time, by countSLOC(), a custom internal function.\n\n", countOffilesExplored)
}

func countSLOC() { // ::: - -

	// todo, do a regular expression to extract the file names (last / to .go inclusive)
	// ... then, print out the names of those files over in about_app()

	numberOfFilesExplored := 0

	filenameOfThisFile1 := "/Users/quasar/grokMusic6/main.go"
	blankLines1, singleComments1, commentBlock11, commentBlock21, commentBlock31, runes11, runes21, runes31, totalLines1, nonEmptyLines1 := reportSLOCstats(filenameOfThisFile1)
	numberOfFilesExplored++
	filenameOfThisFile2 := "/Users/quasar/grokMusic6/go.mod"
	blankLines2, singleComments2, commentBlock12, commentBlock22, commentBlock32, runes12, runes22, runes32, totalLines2, nonEmptyLines2 := reportSLOCstats(filenameOfThisFile2)
	numberOfFilesExplored++
	/*
		filenameOfThisFile5 := "/Users/quasar/Jap2-main/elementsOfsloc.go"
		blankLines5, singleComments5, commentBlock15, commentBlock25, commentBlock35, runes15, runes25, runes35, totalLines5, nonEmptyLines5 := reportSLOCstats(filenameOfThisFile5)
		numberOfFilesExplored++
		filenameOfThisFile6 := "/Users/quasar/Jap2-main/functions.go"
		blankLines6, singleComments6, commentBlock16, commentBlock26, commentBlock36, runes16, runes26, runes36, totalLines6, nonEmptyLines6 := reportSLOCstats(filenameOfThisFile6)
		numberOfFilesExplored++
		filenameOfThisFile7 := "/Users/quasar/Jap2-main/globalVariables.go"
		blankLines7, singleComments7, commentBlock17, commentBlock27, commentBlock37, runes17, runes27, runes37, totalLines7, nonEmptyLines7 := reportSLOCstats(filenameOfThisFile7)
		numberOfFilesExplored++
		filenameOfThisFile8 := "/Users/quasar/Jap2-main/locateCard.go"
		blankLines8, singleComments8, commentBlock18, commentBlock28, commentBlock38, runes18, runes28, runes38, totalLines8, nonEmptyLines8 := reportSLOCstats(filenameOfThisFile8)
		numberOfFilesExplored++
		filenameOfThisFile9 := "/Users/quasar/Jap2-main/memoryFunctions.go"
		blankLines9, singleComments9, commentBlock19, commentBlock29, commentBlock39, runes19, runes29, runes39, totalLines9, nonEmptyLines9 := reportSLOCstats(filenameOfThisFile9)
		numberOfFilesExplored++
		filenameOfThisFile10 := "/Users/quasar/Jap2-main/objectsAndMethods.go"
		blankLines10, singleComments10, commentBlock110, commentBlock210, commentBlock310, runes110, runes210, runes310, totalLines10, nonEmptyLines10 := reportSLOCstats(filenameOfThisFile10)
		numberOfFilesExplored++
		filenameOfThisFile11 := "/Users/quasar/Jap2-main/pick_a_card_functions.go"
		blankLines11, singleComments11, commentBlock111, commentBlock211, commentBlock311, runes111, runes211, runes311, totalLines11, nonEmptyLines11 := reportSLOCstats(filenameOfThisFile11)
		numberOfFilesExplored++
		filenameOfThisFile14 := "/Users/quasar/Jap2-main/pick_a_card_random_all.go"
		blankLines14, singleComments14, commentBlock114, commentBlock214, commentBlock314, runes114, runes214, runes314, totalLines14, nonEmptyLines14 := reportSLOCstats(filenameOfThisFile14)
		numberOfFilesExplored++
		filenameOfThisFile15 := "/Users/quasar/Jap2-main/prompts&directions.go"
		blankLines15, singleComments15, commentBlock115, commentBlock215, commentBlock315, runes115, runes215, runes315, totalLines15, nonEmptyLines15 := reportSLOCstats(filenameOfThisFile15)
		numberOfFilesExplored++
		filenameOfThisFile16 := "/Users/quasar/Jap2-main/statsFunctions.go"
		blankLines, singleComments, commentBlock01, commentBlock02, commentBlock03, runes01, runes02, runes03, totalLines16, nonEmptyLines16 := reportSLOCstats(filenameOfThisFile16)
		numberOfFilesExplored++
	*/

	filenameOfThisFile17 := "/Users/quasar/grokMusic6/countSLOC.go"

	blankLines91, singleComments91, commentBlock104, commentBlock205, commentBlock206, runes104, runes205, runes306, totalLines17, nonEmptyLines17 := reportSLOCstats(filenameOfThisFile17)
	numberOfFilesExplored++

	countOffilesExplored = numberOfFilesExplored // Used only in countSLOC() and the associated about_app()

	totalLines := totalLines1 + totalLines2 + totalLines17

	nonEmptyLines := nonEmptyLines1 + nonEmptyLines2 + nonEmptyLines17

	blankLinesTotal := blankLines2 + blankLines1 + blankLines91

	singleCommentsTotal := singleComments2 + singleComments1 + singleComments91

	commentBlock1Total := commentBlock12 + commentBlock11 + commentBlock104
	commentBlock2Total := commentBlock22 + commentBlock21 + commentBlock205
	commentBlock3Total := commentBlock32 + commentBlock31 + commentBlock206

	runes1Total := runes12 + runes11 + runes104
	runes2Total := runes22 + runes21 + runes205
	runes3Total := runes32 + runes31 + runes306
	// was: runes3Total := runes315 + runes03 + runes314 + runes311 + runes310 + runes39 + runes38 + runes37 + runes36 + runes35 + runes32 + runes31 + runes306
	// 
	grandTotal := blankLinesTotal + singleCommentsTotal + commentBlock2Total + runes2Total

	sumOfCodePlusNon := grandTotal + nonEmptyLines

	fileHandle, err := os.OpenFile("musicAppLog.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check_error(err)

	currentTime := time.Now()

	_, err29 := fmt.Fprintf(fileHandle, "\nWhen the App began at: %s", currentTime.Format("15:04:05 on Monday 01-02-2006"))
	check_error(err29)
	_, err1 := fmt.Fprintf(fileHandle, "\n\nThe Total lines of Code (exclusive of data) = %d t-SLOC\n", totalLines)
	check_error(err1)
	_, err19 := fmt.Fprintf(fileHandle, "and, the Total lines of executable Code = %d e-SLOC\n\n", nonEmptyLines)
	check_error(err19)

	fmt.Printf("\n\tTotal lines of Code (exclusive of data) = %d t-SLOC\n\n", totalLines)

	fmt.Printf("\tTotal lines of executable Code = %d e-SLOC\n\n", nonEmptyLines)
	fmt.Printf("\tBlnkLns:%d + SnglCmLns:%d + ComBlks:%d + runes:%d = total of cmnts+spc = %d lines\n\n", blankLinesTotal, singleCommentsTotal, commentBlock2Total, runes2Total, grandTotal)

	fmt.Printf("\tTotal of comments etc. + e-SLOC = %d = t-SLOC\n\n", sumOfCodePlusNon)

	if runes3Total > 0 || runes1Total > 0 || commentBlock3Total > 0 || commentBlock1Total > 0 { // if any of these was > 0
		fmt.Println("\n\n === hey we actually got something from where there should not have been anything === \n\n")
	}

}

/*
.
*/
func reportSLOCstats(filepath string) (blankLines, singleComments, commentBlock1, commentBlock2, commentBlock3, runes1, runes2, runes3, totalLines, sloc int) {
	// Patterns to identify comments, blank lines, and strings
	singleLineCommentPattern := `^\s*//`
	multiLineCommentPattern := `(?s)/\*.*?\*/`
	blankLinePattern := `^\s*$`
	stringLiteralPattern := `(?s)"(?:\\.|[^\\"])*?"|` + "`(?:\\.|[^`])*?`"

	// Compile regular expressions
	singleLineCommentRe := regexp.MustCompile(singleLineCommentPattern)
	multiLineCommentRe := regexp.MustCompile(multiLineCommentPattern)
	blankLineRe := regexp.MustCompile(blankLinePattern)
	stringLiteralRe := regexp.MustCompile(stringLiteralPattern)

	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		// return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalLines = 0
	sloc = 0
	inMultiLineComment := false
	inMultiLineString := false

	for scanner.Scan() {
		line := scanner.Text()
		totalLines++

		// ::: Check for blank lines
		if blankLineRe.MatchString(line) {
			blankLines++
			continue
		}

		// ::: Check for single-line comments
		if singleLineCommentRe.MatchString(line) {
			singleComments++
			continue
		}

		// ::: Check for multi-line comment blocks
		if inMultiLineComment {
			if strings.Contains(line, "*/") {
				inMultiLineComment = false
				line = multiLineCommentRe.ReplaceAllString(line, "")
				if blankLineRe.MatchString(line) || singleLineCommentRe.MatchString(line) {
					commentBlock1++ // Does not normally accumulate anything.
					continue
				}
			} else {
				commentBlock2++ // This is where we find lines that match.
				continue
			}
		}
		if strings.Contains(line, "/*") {
			inMultiLineComment = true
			line = multiLineCommentRe.ReplaceAllString(line, "")
			if blankLineRe.MatchString(line) || singleLineCommentRe.MatchString(line) { // blankLines, singleComments, commentBlock1, commentBlock2, commentBlock3, runes1, runes2, runes3
				commentBlock3++ // Does not normally accumulate anything.
				continue
			}
		}

		// ::: Check for multi-line strings // string literals // Runes
		if inMultiLineString {
			if strings.Count(line, "`")%2 != 0 || strings.Count(line, "\"")%2 != 0 {
				inMultiLineString = false
				line = stringLiteralRe.ReplaceAllString(line, "")
				if blankLineRe.MatchString(line) || singleLineCommentRe.MatchString(line) {
					runes1++ // Does not normally accumulate anything.
					continue
				}
			} else {
				runes2++ // This is where we find lines that match.
				continue
			}
		}
		if strings.Count(line, "`")%2 != 0 || strings.Count(line, "\"")%2 != 0 {
			inMultiLineString = true
			line = stringLiteralRe.ReplaceAllString(line, "")
			if blankLineRe.MatchString(line) || singleLineCommentRe.MatchString(line) {
				runes3++ // Does not normally accumulate anything.
				continue
			}
		}

		sloc++
	}

	if err := scanner.Err(); err != nil {
		// return 0, 0, err
	}

	return blankLines, singleComments, commentBlock1, commentBlock2, commentBlock3, runes1, runes2, runes3, totalLines, sloc
}

// Creates a func named check_error which takes one parameter "e" of type error
func check_error(e error) { // ::: - -
	if e != nil {
		panic(e) // use panic() to display error code
	}
}
