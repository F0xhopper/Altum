/*
Copyright © 2025 Eden Phillips
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a deep work timer",
	Long: `Start a timer for a deep work session. The timer will run until you press Enter.
After stopping, you'll be prompted for a rating and notes about the session.`,
	Run: func(cmd *cobra.Command, args []string) {
		dailyNotesFolderPath := viper.GetString("daily_notes_folder_path")
		dateFormat := viper.GetString("date_format")
		
		if dailyNotesFolderPath == "" {
			fmt.Fprintf(os.Stderr, "Error: daily_notes_folder_path is required. Please set it using:\n")
			fmt.Fprintf(os.Stderr, "  altum config set daily_notes_folder_path <folder_path>\n")
			fmt.Fprintf(os.Stderr, "  or use --daily_notes_folder_path flag\n")
			os.Exit(1)
		}
		
		
		fmt.Println("Deep work timer started. Press Enter to end the current session...")
		startTime := time.Now()
		
		
		stopChan := make(chan bool)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
			stopChan <- true
		}()
		
		
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-stopChan:
				elapsed := time.Since(startTime)
				minutes := int(elapsed.Minutes())
				seconds := int(elapsed.Seconds()) % 60
				
				fmt.Printf("\nTimer stopped. Duration: %d minutes %d seconds\n", minutes, seconds)
				
				
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Rate this session (1-10): ")
				rateInput, _ := reader.ReadString('\n')
				rate := rateInput[:len(rateInput)-1] 
				
				
				fmt.Print("Notes (press Enter when done): ")
				notesInput, _ := reader.ReadString('\n')
				notes := notesInput[:len(notesInput)-1] 
				
				
				today := time.Now().Format(dateFormat)
				noteFileName := fmt.Sprintf("%s.md", today)
				noteFilePath := filepath.Join(dailyNotesFolderPath, noteFileName)
				
				
				file, err := os.OpenFile(noteFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: Failed to open note file: %v\n", err)
					os.Exit(1)
				}
				defer file.Close()
				
				titleRe := regexp.MustCompile(`^## Altum Work Sessions$`)
				sessionTitleRe := regexp.MustCompile(`^#### Session \d+`)
				

				isTitleLineFound := false
				sessionCount := 1

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())
					if !isTitleLineFound && titleRe.MatchString(line) {
						isTitleLineFound = true
						continue
					}
					if isTitleLineFound && sessionTitleRe.MatchString(line) {
						sessionCount++
						continue
					}
				}

				if err := scanner.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "Error: Failed to read note file: %v\n", err)
					os.Exit(1)
				}
                entry := ""
				if !isTitleLineFound{
					entry += fmt.Sprintf("\n## Altum Work Sessions\n")
				}
				sessionStartTime := startTime.Format("15:04:05")
				sessionEndTime := time.Now().Format("15:04:05")
				entry += fmt.Sprintf("\n#### Session %d\n", sessionCount)
				entry += fmt.Sprintf("- Time: %s - %s\n", sessionStartTime, sessionEndTime)
				entry += fmt.Sprintf("- Duration: %d minutes %d seconds\n", minutes, seconds)
				entry += fmt.Sprintf("- Rate: %s/10\n", rate)
				if notes != "" {
					entry += fmt.Sprintf("- Notes: %s\n", notes)
				}
				
				if _, err := file.WriteString(entry); err != nil {
					fmt.Fprintf(os.Stderr, "Error: Failed to write to note file: %v\n", err)
					os.Exit(1)
				}
				
				fmt.Printf("\nSession logged to: %s\n", noteFilePath)
				return
				
			case <-ticker.C:
				elapsed := time.Since(startTime)
				minutes := int(elapsed.Minutes())
				seconds := int(elapsed.Seconds()) % 60
				fmt.Printf("\rTimer: %02d:%02d", minutes, seconds)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}